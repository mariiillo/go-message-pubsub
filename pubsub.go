package main

import (
	"os/exec"
	"time"

	"fmt"

	"log"

	"os"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
)

func newPublisherSocket(url string) (mangos.Socket, error) {
	socket, err := pub.NewSocket()
	if err != nil {
		return nil, err
	}

	socket.AddTransport(ipc.NewTransport())
	socket.AddTransport(tcp.NewTransport())

	err = socket.Listen(url)
	if err != nil {
		return nil, err
	}

	return socket, nil
}

func newSubscriberSocket(url string) (mangos.Socket, error) {
	socket, err := sub.NewSocket()
	if err != nil {
		return nil, err
	}

	socket.AddTransport(tcp.NewTransport())
	socket.AddTransport(ipc.NewTransport())

	err = socket.Dial(url)
	if err != nil {
		return nil, err
	}

	return socket, nil
}

func subscribe(socket mangos.Socket, topic string) error {
	err := socket.SetOption(mangos.OptionSubscribe, []byte(topic))
	if err == nil {
		err = socket.SetOption(mangos.OptionRecvDeadline, 10*time.Second)
	}
	return err
}

func publish(socket mangos.Socket, topic, message string) error {
	err := socket.Send([]byte(fmt.Sprintf("%s|%s", topic, message)))
	return err
}

func receive(socket mangos.Socket) (string, error) {
	message, err := socket.Recv()
	return string(message), err
}

func runServer(url string, topics []string) {
	socket, err := newPublisherSocket(url)
	if err != nil {
		log.Fatalf("Cannot listen on %s: %s\n", url, err.Error())
	}

	for i := 0; i < 5; i++ {
		for _, topic := range topics {
			time.Sleep(1 * time.Second)
			fmt.Printf("Publishing a message for topic %s\n", topic)
			err := publish(socket, topic, fmt.Sprintf("Message for %s", topic))
			if err != nil {
				log.Fatalf("Cannot publish message for topic %s: %s\n", topic, err.Error())
			}
		}
	}
}

func runClient(name, url string, topics []string) {
	socket, err := newSubscriberSocket(url)
	if err != nil {
		log.Fatalf("Cannot dial into %s: %s\n", url, err.Error())
	}

	for _, topic := range topics {
		err := subscribe(socket, topic)
		if err != nil {
			log.Fatalf("Cannot subscribe to topic %s: %s\n", topic, err.Error())
		}
	}

	for i := 0; i < 5*len(topics); i++ {
		message, err := receive(socket)
		if err != nil {
			log.Fatalf("Error receiving message: %s\n", message)
		}
		fmt.Printf("Client %s received: %s\n", name, message)
	}
}

func main() {
	url := "tcp://localhost:56565"
	if len(os.Args) == 1 {
		client1 := exec.Command("./pubsub", "C1", "Technology")
		client1.Stdout = os.Stdout
		client1.Stderr = os.Stderr
		client2 := exec.Command("./pubsub", "C2", "Technology", "Weather")
		client2.Stdout = os.Stdout
		client2.Stderr = os.Stderr
		client3 := exec.Command("./pubsub", "C3", "Finance")
		client3.Stdout = os.Stdout
		client3.Stderr = os.Stderr

		fmt.Println("Starting client 1")
		err := client1.Start()
		if err != nil {
			log.Fatalf("Failed starting client1: %s", err.Error())
		}

		fmt.Println("Starting client 2")
		err = client2.Start()
		if err != nil {
			log.Fatalf("Failed starting client2: %s", err.Error())
		}

		fmt.Println("Starting client 3")
		err = client3.Start()
		if err != nil {
			log.Fatalf("Failed starting client3: %s", err.Error())
		}

		fmt.Println("Starting the server")
		runServer(url, []string{"Technology", "Weather", "Finance"})

		time.Sleep(1 * time.Second)
		fmt.Println("Waiting for the clients to exit")
		client1.Wait()
		client2.Wait()
		client3.Wait()
		fmt.Println("Server ends.")
	} else {
		fmt.Println(os.Args[1], "is starting")
		runClient(os.Args[1], url, os.Args[2:])
		fmt.Println("Client", os.Args[1], "ends.")
	}
}
