package queue

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
)

var ErrInvalidState = errors.New("invalid state")
var ErrQueueEmpty = errors.New("queue is empty")

type Queue struct {
	messages []Message
}

// QueueMarshal is a helper struct for JSON marshalling/unmarshalling
type QueueMarshal struct {
	Messages []Message `json:"messages"`
}

func New() *Queue {
	return &Queue{
		messages: make([]Message, 0),
	}
}

// Peek returns the first message without removing it from the queue.
// If the queue is empty, it returns nil.
func (q *Queue) Peek() *Message {
	if q.Size() == 0 {
		return nil
	}

	return &q.messages[0]
}

// Add a new message to the end of the queue.
func (q *Queue) Enqueue(message Message) {
	q.messages = append(q.messages, message)
}

func (q *Queue) Dequeue(messageId uuid.UUID) error {
	if q.Size() == 0 {
		return ErrQueueEmpty
	}

	message := q.messages[0]
	if message.ID != messageId {
		return ErrInvalidState
	}

	// Remove the first message from the queue
	q.messages = q.messages[1:]

	return nil
}

func (q *Queue) Size() int {
	return len(q.messages)
}

func (q *Queue) PersistToJsonFile() error {
	data, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("queue.json", data, 0644)
}

func LoadQueueFromJsonFile() (*Queue, error) {
	data, err := os.ReadFile("queue.json")
	if err != nil {
		return nil, err
	}

	var q Queue
	if err := json.Unmarshal(data, &q); err != nil {
		return nil, err
	}

	return &q, nil
}

func (q *Queue) MarshalJSON() ([]byte, error) {
	return json.Marshal(QueueMarshal{
		Messages: q.messages,
	})
}

func (q *Queue) UnmarshalJSON(data []byte) error {
	tempQueue := QueueMarshal{}
	if err := json.Unmarshal(data, &tempQueue); err != nil {
		return err
	}

	q.messages = tempQueue.Messages

	return nil
}
