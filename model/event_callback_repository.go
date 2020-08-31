package model

import (
	"fmt"

	"github.com/imdario/mergo"
	"gitlab.kardinal.ai/coretech/esther/logging"
	"gitlab.kardinal.ai/coretech/esther/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindEventCallbacksByPlanId Find all events by planId
func FindEventCallbacksByPlanId(planID string) ([]EventCallback, error) {
	collection := persistence.GetCollection(EventCallback{})
	ctx, _ := persistence.GetContext()
	cur, err := collection.Find(ctx, bson.D{primitive.E{Key: "planid", Value: planID}})
	if err != nil {
		err := fmt.Sprintf("Can't find EventCallback with plan-id: %s", planID)
		logging.Logger.Error(err)
		return []EventCallback{}, fmt.Errorf(err)
	}
	eventCallbacks := make([]EventCallback, 0)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var eventCallback EventCallback
		err := cur.Decode(&eventCallback)
		if err != nil {
			logging.Logger.Error(err)
		}
		eventCallbacks = append(eventCallbacks, eventCallback)
	}

	return eventCallbacks, nil
}

// FindEventCallbackById Find one event by his ID
func FindEventCallbackById(planID string, eventID string) (EventCallback, error) {
	eventCallbackInterface := persistence.FindOne(EventCallback{ID: eventID, PlanID: planID})
	if eventCallbackInterface == nil {
		return EventCallback{}, fmt.Errorf("Can't retrieve the event-callback %s in plan %s", eventID, planID)
	}
	eventCallback, _ := (*eventCallbackInterface).(EventCallback)
	return eventCallback, nil
}

// CreateEventCallback Create one event
func CreateEventCallback(planID string, eventCallback EventCallback) (EventCallback, error) {
	eventCallback.PlanID = planID
	id := persistence.InsertOne(eventCallback)
	if id == "" {
		return EventCallback{}, fmt.Errorf("Can't create the event-callback in plan %s", eventCallback.PlanId())
	}
	eventCallback.ID = id
	return eventCallback, nil
}

// UpdateEventCallback Create one event
func UpdateEventCallback(planID string, eventID string, eventCallback EventCallback) (EventCallback, error) {
	ec, err := FindEventCallbackById(planID, eventID)
	if err != nil {
		return EventCallback{}, err
	}
	if err := mergo.Merge(&ec, eventCallback, mergo.WithOverride); err != nil {
		return EventCallback{}, err
	}
	if persistence.ReplaceOne(ec) == false {
		return EventCallback{}, fmt.Errorf("Can't update the event-callback %s in plan %s", eventCallback.Id(), eventCallback.PlanId())
	}
	return ec, nil
}

// DeleteEventCallbacksById Delete one event
func DeleteEventCallbacksById(planID string, eventID string) error {
	if persistence.DeleteOne(EventCallback{ID: eventID, PlanID: planID}) == false {
		return fmt.Errorf("Can't delete the event-callback %s in plan %s", eventID, planID)
	}
	return nil
}

// DeleteEventCallbacksByPlanId Delete many events by their planId
func DeleteEventCallbacksByPlanId(planID string) error {
	collection := persistence.GetCollection(EventCallback{})
	ctx, _ := persistence.GetContext()
	dr, err := collection.DeleteMany(ctx, bson.D{
		primitive.E{Key: "planid", Value: planID},
	})
	if err != nil && dr.DeletedCount >= 1 {
		err := fmt.Sprintf("Can't delete EventCallback with plan-id: %s", planID)
		logging.Logger.Error(err)
		return fmt.Errorf(err)
	}

	return nil
}
