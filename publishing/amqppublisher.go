package publishing

import (
	"fmt"

	"github.com/marrick66/amqp-eventbus/converters"
	"github.com/streadway/amqp"
)

//AMQPPublisher struct stores the channel to publish on,
//as well as a map of the exchange/topic combinations that have
//already been declared.
type AMQPPublisher struct {
	channel   *amqp.Channel
	exchanges map[string]bool
	converter converters.ByteConverter
}

//New creates an AMQPPublisher instance from an existing
//Channel object.
func New(channel *amqp.Channel, converter converters.ByteConverter) (Publisher, error) {
	if channel == nil {
		return nil, fmt.Errorf("Channel cannot be nil")
	}

	if converter == nil {
		return nil, fmt.Errorf(("Converter cannot be nil"))
	}

	return &AMQPPublisher{
		channel:   channel,
		converter: converter,
		exchanges: make(map[string]bool)}, nil
}

//Publish will declare the exchange/topic combination on the channel, if it has not already done
//so. Otherwise, it just converts the event and sends it to the channel.
func (publisher *AMQPPublisher) Publish(exchange string, topic string, event interface{}) error {
	var body []byte
	var err error

	//If we don't know about this exchange, attempt to declare it and store that it's up.
	if _, ok := publisher.exchanges[exchange]; !ok {
		if err = publisher.channel.ExchangeDeclare(
			exchange,
			"topic",
			true,
			true,
			false,
			false,
			nil); err != nil {
			return err
		}
		publisher.exchanges[exchange] = true
	}

	//Use the publisher's converter to get the raw bytes to publish:
	if body, err = publisher.converter.ToBytes(event); err != nil {
		return err
	}

	//Finally, send the message to the exchange/topic:
	if err = publisher.channel.Publish(
		exchange,
		topic,
		true,
		false,
		amqp.Publishing{
			Body: body}); err != nil {
		return err
	}

	return nil
}
