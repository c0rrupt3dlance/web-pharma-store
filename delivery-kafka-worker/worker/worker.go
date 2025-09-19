package worker

import (
	"context"
	"delivery-kafka-worker/producer"
	"delivery-kafka-worker/repository"
	"log"
	"time"
)

type Worker struct {
	repo     *repository.Repository
	producer *producer.Producer
	interval time.Duration
}

func NewWorker(repo *repository.Repository, producer *producer.Producer, interval time.Duration) *Worker {
	return &Worker{repo: repo, producer: producer, interval: interval}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processBatch(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) {
	events, err := w.repo.FetchUnprocessed(ctx, 10)
	if err != nil {
		log.Printf("failed to fetch outbox: %v", err)
		return
	}

	for _, e := range events {
		topic := e.AggregateType
		log.Printf("Sending event to Kafka topic=%s key=%s payload=%s\n", topic, e.EventType, string(e.Payload))
		if err := w.producer.Send(ctx, topic, e.EventType, e.Payload); err != nil {
			log.Printf("failed to send event %d: %v", e.Id, err)
			continue
		}

		if err := w.repo.MarkProcessed(ctx, e.Id); err != nil {
			log.Printf("failed to mark processed %d: %v", e.Id, err)
		}
	}
}
