package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type amqpConnection struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getURI() string {
	// "amqp://guest:guest@localhost:5672/"

	broker := viper.Get("broker").(string)
	port := viper.Get("port").(string)
	user := viper.Get("user").(string)
	pass := viper.Get("pass").(string)
	virthost := viper.Get("virthost").(string)

	return "amqp://" + user + ":" + pass + "@" + broker + ":" + port + virthost
}

func (a *amqpConnection) connectToBroker() {
	uri := getURI()
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	a.conn = conn

	fmt.Printf("Connected to broker: %s\n", uri)
	// defer conn.Close()
}

func (a *amqpConnection) connectToChannel() {
	ch, err := a.conn.Channel()
	failOnError(err, "Failed to open a channel")
	a.ch = ch

	fmt.Println("Connected to channel")
	// defer ch.Close()
}

func (a *amqpConnection) declareQueue(queueName string) {
	q, err := a.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	a.q = q
	fmt.Printf("Connected to queue %s\n", queueName)
}

func (a *amqpConnection) publishToQueue(exchange string, body []byte) {
	err := a.ch.Publish(
		exchange, // exchange
		a.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func (a *amqpConnection) consumeFromQueue() {
	msgs, err := a.ch.Consume(
		a.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			req := &request{}
			err := json.Unmarshal(d.Body, *req)
			if err != nil {
				return
			}

			rsp := req.processRequest()

			body, err := json.Marshal(rsp)
			if err != nil {
				return
			}

			a.publishToQueue("", body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func getRandomQuoteRate() float64 {
	random := rand.Float64()
	rate := math.Mod(random, 17)

	return rate
}
