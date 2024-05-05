package consumers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jis4nx/go-ecom/pkg/rabbit"
	"golang.org/x/sync/errgroup"
)

func StartConsumer() {
	conn, err := rabbit.ConnectRabbitMQ("guest", "guest", "rabbit:5672", "/")
	if err != nil {
		log.Fatal("Start Consumer", err.Error())
	}

	client, err := rabbit.NewRabbitClient(conn)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create Exchange
	if err = client.CreateExchange("product_events", "topic", true, false); err != nil {
		log.Println(err)
	}

	// Decalrig New Queue
	if err = client.NewQueue("product_created", true, false); err != nil {
		log.Fatal(err)
	}

	// Create Binding with Queue & Routing key
	if err = client.CreateBinding("product_created", "product.created.*", "product_events"); err != nil {
		log.Fatal(err)
	}

	messageBus, err := client.Consume("product_created", "testConsumer1", false)
	if err != nil {
		log.Fatal(err.Error())
	}

	var blocking chan struct{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(10)

	go func() {
		for message := range messageBus {
			msg := message
			g.Go(func() error {
				if err := msg.Ack(false); err != nil {
					log.Println(err.Error())
					return err
				}
				log.Println("Acknowldeged Message")

				fmt.Println(msg)
				return nil
			})
		}
	}()

	defer cancel()
	<-blocking
}
