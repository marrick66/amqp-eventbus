package subscribing

import (
	"fmt"
	"log"
	"sync"

	"github.com/marrick66/amqp-eventbus/converters"
	"github.com/streadway/amqp"
)

//subscriptionKey is used to keep track of the registered
//subscribers for the event feed.
type subscriptionKey struct {
	exchange string
	topic    string
}

//AMQPEventFeed handles maintenance of subscribers
//to the channel.
type AMQPEventFeed struct {
	channel         *amqp.Channel
	converter       converters.ByteConverter
	subscriptions   map[subscriptionKey][]Subscriber
	subscriberMsgs  chan interface{}
	errors          <-chan *amqp.Error
	subscriberGroup *sync.WaitGroup
}

//New creates an instance of the AMQPEventFeed for subscribing. This is not self healing, since
//the AMQP connection or channel can close at any time. This would stop the message loop on each
//subscriber preventing any further handling.
func New(channel *amqp.Channel, converter converters.ByteConverter) (EventFeed, error) {
	if channel == nil {
		return nil, fmt.Errorf("Channel cannot be nil")
	}

	if converter == nil {
		return nil, fmt.Errorf("Converter cannot be nil")
	}

	errors := make(chan *amqp.Error)

	return &AMQPEventFeed{
		channel:         channel,
		converter:       converter,
		subscriptions:   make(map[subscriptionKey][]Subscriber),
		subscriberMsgs:  make(chan interface{}),
		errors:          channel.NotifyClose(errors),
		subscriberGroup: &sync.WaitGroup{}}, nil
}

//Subscribe assigns a temporary queue for the subscriber, and binds it to the exchange with
//the topic specified. A goroutine event loop for handling these items is started once subscribed.
func (feed *AMQPEventFeed) Subscribe(exchange string, topic string, subscriber Subscriber) error {
	if subscriber == nil {
		return fmt.Errorf("subscriber cannot be nil")
	}

	key := subscriptionKey{
		exchange: exchange,
		topic:    topic}

	//Create a list of subscribers for this exchange/topic combination if it doesn't exist:
	if _, ok := feed.subscriptions[key]; !ok {
		feed.subscriptions[key] = append(feed.subscriptions[key], subscriber)
	} else {
		//Otherwise, just append the subscriber if it doesn't already exist in the list. Iterating
		//over the list isn't ideal, but easy for now.
		for index := range feed.subscriptions[key] {
			if feed.subscriptions[key][index] == subscriber {
				return fmt.Errorf("this subscriber is already registered for this exchange/topic")
			}
		}

		feed.subscriptions[key] = append(feed.subscriptions[key], subscriber)
	}

	//Declare the exchange, since this may be called prior to a publish.
	//The exchange declaration settings have to be the same as the publisher's,
	//so next on the TODO list should be to extract these into a common configuration object.
	if err := feed.channel.ExchangeDeclare(
		exchange,
		"topic",
		true,
		true,
		false,
		false,
		nil); err != nil {
		return err
	}

	var queue amqp.Queue
	var err error

	//Declare a temporary exclusive queue for this subscriber.
	if queue, err = feed.channel.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil); err != nil {
		return nil
	}

	//Bind the exchange/topic to the temp queue so messages to the exchange/topic
	//are forwarded.
	if err = feed.channel.QueueBind(queue.Name, topic, exchange, false, nil); err != nil {
		return err
	}

	var msgs <-chan amqp.Delivery

	if msgs, err = feed.channel.Consume(
		queue.Name,
		topic,
		true,
		true,
		false,
		false,
		nil); err != nil {
		return err
	}

	//Start the handling loop for this subscriber:
	feed.subscriberGroup.Add(1)
	go feed.startHandling(msgs, &key, subscriber)

	return nil
}

//Handles all messages until the AMQP channel or connection is closed (or errors out). Once stopped, a message
//is sent to the subscriber event channel and the waitgroup is signaled.
func (feed *AMQPEventFeed) startHandling(msgs <-chan amqp.Delivery, key *subscriptionKey, subscriber Subscriber) {
	var event interface{}
	var err error

	for msg := range msgs {
		if event, err = feed.converter.FromBytes(msg.Body, subscriber.EventType()); err != nil {
			log.Printf("Error converting message: %s", err)
			continue
		}

		if err = subscriber.Handle(event); err != nil {
			log.Printf("Error handling message %v: %s\n", msg, err)
			continue
		}
	}

	//Signal that this subscriber stopped listening, for whatever reason.
	feed.subscriberMsgs <- &subscriberStopped{
		subscriberMessage: subscriberMessage{
			key:        key,
			subscriber: subscriber}}

	feed.subscriberGroup.Done()
}
