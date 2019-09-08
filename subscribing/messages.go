package subscribing

//subscriberMessage is the base message type
//that's sent to the feed for updates.
type subscriberMessage struct {
	key        *subscriptionKey
	subscriber Subscriber
}

//subscriberStopped is sent when the subscriber
//stops handling messages.
type subscriberStopped struct {
	subscriberMessage
}
