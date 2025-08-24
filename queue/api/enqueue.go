package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nicograef/go-playground/queue/queue"
)

type EnqueueRequest struct {
	Payload any `json:"payload"`
}
type EnqueueResponse struct {
	MessageId string `json:"messageId"`
	QueueSize int    `json:"queueSize"`
}

// read the payload from the request body, create a new message with a new UUID and the current timestamp, add it to the queue, and return the ID of the new message in the response.
func NewEnqueueHandler(q *queue.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateMethod(w, r, http.MethodPost) {
			return
		}

		requestBody := EnqueueRequest{}
		if !readJSONRequest(w, r, &requestBody) {
			return
		}

		message := queue.Message{
			ID:        uuid.New(),
			Timestamp: time.Now(),
			Payload:   requestBody.Payload,
		}

		q.Enqueue(message)

		sendJSONResponse(w, EnqueueResponse{
			MessageId: message.ID.String(),
			QueueSize: q.Size(),
		})
	}
}
