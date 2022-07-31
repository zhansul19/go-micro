package main

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/zhansul19/go-micro/listener-service/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//connect to rabbitmq
	conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Println("connected to rabbitmq")

	//start listening
	log.Println("start listening for and consu,ing rabbitmq messages...")

	//create consumer
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		log.Panic(err)
	}
	//watch the queue and consume events

	err = consumer.Listen([]string{"log.WARNING", "log.ERROR", "log.INFO"})
	if err != nil {
		log.Println(err)
	}

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("connection to rabbitmq not ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
