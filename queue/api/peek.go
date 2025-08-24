package api

import (
	"log"
	"net/http"

	"github.com/nicograef/go-playground/queue/queue"
)

type PeekResponse struct {
	QueueSize int           `json:"queueSize"`
	Message   queue.Message `json:"message"`
}

func NewPeekHandler(q *queue.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validateMethod(w, r, http.MethodPost)

		message := q.Peek()
		queueSize := q.Size()

		if message == nil && queueSize != 0 {
			log.Printf("WARN: Queue is NOT empty, but cannot get next message!")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, PeekResponse{
			QueueSize: queueSize,
			Message:   *message,
		})
	}
}
