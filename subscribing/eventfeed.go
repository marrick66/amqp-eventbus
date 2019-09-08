package subscribing

//EventFeed only allows a Subscriber implementer to receive
//events from an exchange/topic combination.
type EventFeed interface {
	Subscribe(exchange string, topic string, subscriber Subscriber) error
}
