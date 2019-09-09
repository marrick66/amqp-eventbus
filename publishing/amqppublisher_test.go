package publishing_test

import (
	"testing"

	"github.com/marrick66/amqp-eventbus/converters"
	"github.com/marrick66/amqp-eventbus/publishing"
	"github.com/streadway/amqp"
)

var connStr = "amqp://guest:guest@localhost:5672"

func CreateLocalPublisher() (*amqp.Connection, publishing.Publisher, error) {
	var conn *amqp.Connection
	var channel *amqp.Channel
	var pub publishing.Publisher
	var err error

	if conn, err = amqp.Dial(connStr); err != nil {
		return nil, nil, err
	}

	if channel, err = conn.Channel(); err != nil {
		return nil, nil, err
	}

	if pub, err = publishing.New(channel, &converters.JSONConverter{}, nil); err != nil {
		return nil, nil, err
	}

	return conn, pub, nil
}

func TestPublishWithLocalIntegration(t *testing.T) {
	var pub publishing.Publisher
	var err error

	if _, pub, err = CreateLocalPublisher(); err != nil {
		t.Errorf("TestPublishWithLocalIntegration: Failed to get publisher %s", err)
		return
	}

	message := "Hello!"
	exchange := "TestExchange"
	topic := "TestTopic"
	if err = pub.Publish(exchange, topic, message); err != nil {
		t.Errorf("Publish: Unable to publish to %s/%s with message %s: %s", exchange, topic, message, err)
		return
	}
}

func TestPublishWithClosedConnection(t *testing.T) {
	var err error
	var pub publishing.Publisher
	var conn *amqp.Connection

	if conn, pub, err = CreateLocalPublisher(); err != nil {
		t.Errorf("TestPublishWithClosedConnection: Failed to get publisher %s", err)
		return
	}

	conn.Close()

	if err = pub.Publish("test", "test", "test"); err == nil {
		t.Errorf("TestPublishWithClosedConnection: Expected an error, but received none")
		return
	}
}
