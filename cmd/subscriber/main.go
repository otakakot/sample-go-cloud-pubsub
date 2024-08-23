package main

import (
	"context"
	"log/slog"
	"time"

	"cloud.google.com/go/pubsub"
)

func main() {
	cli, err := pubsub.NewClient(context.Background(), "local-project")
	if err != nil {
		panic(err)
	}

	topicID := "topic"

	topic := cli.Topic(topicID)

	if ok, err := topic.Exists(context.Background()); err != nil {
		panic(err)
	} else if !ok {
		if _, err := cli.CreateTopic(context.Background(), topicID); err != nil {
			panic(err)
		}
	}

	subscriptionID := "subscription"

	subscription := cli.Subscription(subscriptionID)

	// subscription.ReceiveSettings.NumGoroutines = 1

	subscription.ReceiveSettings.MaxOutstandingMessages = 1

	if ok, err := subscription.Exists(context.Background()); err != nil {
		panic(err)
	} else if !ok {
		if _, err := cli.CreateSubscription(context.Background(), subscriptionID, pubsub.SubscriptionConfig{
			Topic: topic,
		}); err != nil {
			panic(err)
		}
	}

	slog.InfoContext(context.Background(), "subscribed to topic: "+topicID)

	if err := subscription.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		defer msg.Ack()
		slog.InfoContext(ctx, "received message: "+msg.ID)
		slog.InfoContext(ctx, "do something received message: "+string(msg.Data))
		time.Sleep(3 * time.Second)
		slog.InfoContext(ctx, "finished processing message: "+msg.ID)

	}); err != nil {
		slog.ErrorContext(context.Background(), err.Error())
	}
}
