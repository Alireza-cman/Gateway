package queue

import "github.com/adjust/rmq"

func Queue() {
	connection := rmq.OpenConnection("my service", "tcp", "localhost:6379", 1)
	taskQueue := connection.OpenQueue("tasks")
	delivery := "task payload"
	taskQueue.Publish(delivery)
	// taskConsumer := &TaskConsumer{}
	// taskQueue.AddConsumer("task consumer", taskConsumer)
}
