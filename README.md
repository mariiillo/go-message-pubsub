[nanomsg][1]
============
nanomsg is a socket library that provides several common communication patterns. It aims to make the networking layer
fast, scalable, and easy to use. Implemented in C, it works on a wide range of operating systems with no further depen
dencies.

PubSub
------
PubSub is a communication topology where a single entity called Publisher produces messages that it sends out to other entities called Subscribers. Subscribers may receive everything the publisher sends, or they may subscribe to message subsets called Topics.
Broadcast messages to multiple destinations. Receivers can subscribe to specific topics.

[1]: http://nanomsg.org/%
