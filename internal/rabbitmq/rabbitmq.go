package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type PublisherInterface interface {
	Publish(exchange, routingKey string, body []byte) error
	Close()
}

type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: conn, ch: ch}, nil
}

func (p *Publisher) Publish(exchange, routingKey string, body []byte) error {
	return p.ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) Close() {
	p.ch.Close()
	p.conn.Close()
}

func Consume(url, queueName string, handler func(d amqp.Delivery)) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	log.Printf(" [*] Waiting for messages on queue '%s'. To exit press CTRL+C", queueName)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf(" [x] Received a message: %s", d.Body)
			handler(d)
		}
	}()

	<-forever
	return nil
}

func SetupConsumer(url, exchangeName, routingKey, queueName string, handler func(d amqp.Delivery)) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Printf("Consumer successfully bound to exchange '%s' with routing key '%s' on queue '%s'", exchangeName, routingKey, queueName)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	<-forever
	return nil
}
