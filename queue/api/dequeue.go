package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nicograef/go-playground/queue/queue"
)

type DequeueRequest struct {
	MessageId uuid.UUID `json:"messageId" validate:"required"` // required
}
type DequeueResponse struct {
	QueueSize int `json:"queueSize"`
}

// read the payload from the request body, create a new message with a new UUID and the current timestamp, add it to the queue, and return the ID of the new message in the response.
func NewDequeueHandler(q *queue.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateMethod(w, r, http.MethodPost) {
			return
		}

		requestBody := DequeueRequest{}
		if !readJSONRequest(w, r, &requestBody) {
			return
		}

		err := q.Dequeue(requestBody.MessageId)
		if err != nil {
			if err == queue.ErrQueueEmpty {
				http.Error(w, "Queue is empty", http.StatusBadRequest)
				return
			}

			if err == queue.ErrInvalidState {
				http.Error(w, "Message ID does not match the next message in the queue", http.StatusBadRequest)
				return
			}

			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, DequeueResponse{
			QueueSize: q.Size(),
		})
	}
}
