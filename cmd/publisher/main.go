package main

import (
	"context"
	"log/slog"
	"sync"

	"cloud.google.com/go/pubsub"
)

func main() {
	cli, err := pubsub.NewClient(context.Background(), "local-project")
	if err != nil {
		panic(err)
	}

	topicID := "topic"

	topic := cli.Topic(topicID)

	messages := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	var wg sync.WaitGroup

	for _, message := range messages {
		wg.Add(1)
		go func(msg string) {
			defer wg.Done()
			sid, err := topic.Publish(context.Background(), &pubsub.Message{
				Data: []byte(msg),
			}).Get(context.Background())
			if err != nil {
				slog.InfoContext(context.Background(), "failed to publish message: "+msg)
			}
			slog.InfoContext(context.Background(), "published message: "+msg+" with ID: "+sid)
		}(message)
	}

	wg.Wait()
}
