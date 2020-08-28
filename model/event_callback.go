package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.kardinal.ai/coretech/esther/logging"
	"gitlab.kardinal.ai/coretech/esther/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Payload primitive.M `json:"payload"`
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

// Apply : Apply a modification event at their parent
func (eventCallback EventCallback) Apply() error {
	uri := eventCallback.Parent.URI
	method := eventCallback.Parent.Method
	payload, err := json.Marshal(eventCallback.Parent.Payload)
	if err != nil {
		logging.Logger.Error("Error occured during JSON creation payload")
		return err
	}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(payload))
	if err != nil {
		logging.Logger.WithFields(logging.LogFields{"httpCode": http.StatusInternalServerError, "id": eventCallback.ID, "planId": eventCallback.PlanId, "payload": err}).Warning("Error while constructing publishing request: " + err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.Logger.WithFields(logging.LogFields{"httpCode": http.StatusInternalServerError, "id": eventCallback.ID, "planId": eventCallback.PlanId, "payload": err}).Warning("Error while publishing: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ferr := fmt.Errorf("Error when apply event %s to plan %s\nStatus: %s\nBody : %s", eventCallback.ID, eventCallback.PlanID, resp.Status, err)
			logging.Logger.Error(ferr)
			return ferr
		}
		ferr := fmt.Errorf("Error when apply event %s to plan %s\nStatus: %s\nBody : %s", eventCallback.ID, eventCallback.PlanID, resp.Status, body)
		logging.Logger.Error(ferr)
		return ferr
	}
	return nil
}
