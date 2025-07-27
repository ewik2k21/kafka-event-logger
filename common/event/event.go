package event

import (
	"encoding/json"
	"time"
)

type EventType int

const (
	View EventType = iota
	AddToCart
	Purchase
)

func (et *EventType) ToString() string {
	switch *et {
	case View:
		return "view"
	case AddToCart:
		return "add_to_cart"
	case Purchase:
		return "purchase"
	default:
		return "unknown"
	}
}

func ToEventType(tp string) EventType {
	switch tp {
	case "view":
		return View
	case "add_to_cart":
		return AddToCart
	case "purchase":
		return Purchase
	default:
		return -1
	}
}

type Event struct {
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	EventType EventType `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

func NewEvent(userID, eventID, productID string, eventType EventType) Event {
	return Event{
		UserID:    userID,
		EventID:   eventID,
		ProductID: productID,
		EventType: eventType,
		Timestamp: time.Now(),
	}
}

func (e *Event) ToJson() ([]byte, error) {
	return json.Marshal(e)
}

func FromJson(data []byte) (Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	return event, err
}
