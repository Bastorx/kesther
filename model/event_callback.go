package model

import (
	"fmt"

	"gitlab.citodi.com/coretech/esther/logging"
	"gitlab.citodi.com/coretech/esther/persistence"
	"go.mongodb.org/mongo-driver/mongo"
)

// EventCallback represents a JSON input/output payload of an event callback
type EventCallback struct {
	ID     string `bson:"_id,omitempty" json:"id"`
	PlanID string `json:"planId"`
	Title  string `json:"title"`
	Undo   Action `json:"undo"`
	Parent Action `json:"parent"`
}

// Action represents a possible action (apply to parent, or undo)
type Action struct {
	PlanID  string      `json:"planId"`
	Method  string      `json:"method"`
	URI     string      `json:"uri"`
	Payload interface{} `json:"payload"`
}

// Id : Get id
func (eventCallback EventCallback) Id() string {
	return eventCallback.ID
}

// ResetId : Reset id (Has to create a new Instance)
func (eventCallback EventCallback) ResetId(ID string) persistence.Persistable {
	eventCallback.ID = ID
	return eventCallback
}

// PlanId : Get Plan Id
func (eventCallback EventCallback) PlanId() string {
	return eventCallback.PlanID
}

// EntityName : Get entity name
func (eventCallback EventCallback) EntityName() string {
	return "event_callback"
}

// FromBson : Transform to BSON
func (eventCallback EventCallback) FromBson(sr *mongo.SingleResult) persistence.Persistable {
	err := sr.Decode(&eventCallback)
	if err != nil {
		logging.Logger.Error(fmt.Sprintf("Cannot unmarshal pair : %s", err.Error()))
	}
	return eventCallback
}
