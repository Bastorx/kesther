package persistence

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"gitlab.citodi.com/coretech/esther/logging"
)

const (
	contextTimeout = 10 * time.Second
	mongoDbName    = "esther"
)

var (
	mongoDbHost   string
	mongoDbPort   string
	mongoDbURI    string
	mongoDbClient *mongo.Client
)

// ReadyCheck checks if the package is ready to work
func ReadyCheck() []string {
	checkErrors := []string{}
	if err := checkConnection(); err != nil {
		checkErrors = append(checkErrors, fmt.Sprintf("Not connected to the database: %s", err))
	}
	if len(checkErrors) > 0 {
		return checkErrors
	}
	return nil
}

// Reset will reload the package and clean the data
func Reset() []string {
	errors := initEnv()
	if len(errors) == 0 {
		errors = initConnection()
	}
	if len(errors) == 0 {
		ctx, cancel := GetContext()
		defer cancel()
		if err := mongoDbClient.Database(mongoDbName).Drop(ctx); err != nil {
			errors = append(errors, err.Error())
		} else {
			logging.Logger.WithFields(logging.LogFields{
				"uri":  mongoDbURI,
				"name": mongoDbName,
			}).Info("The database has been reset")
		}
	}
	if len(errors) > 0 {
		logging.Logger.WithField("errors", errors).Error("The reset has failed")
	}
	return errors
}

func init() {
	if errors := initEnv(); len(errors) == 0 {
		initConnection()
	}
}

func initEnv() []string {
	errors := []string{}
	envDbServiceHost := "MONGODB_SERVICE_HOST"
	mongoDbHost = os.Getenv(envDbServiceHost)
	if mongoDbHost == "" {
		errors = append(errors, fmt.Sprintf("%s is not set", envDbServiceHost))
	} else {
		envDbPort := "MONGODB_PORT"
		mongoDbPort = os.Getenv(envDbPort)
		if mongoDbPort == "" {
			errors = append(errors, fmt.Sprintf("%s is not set", envDbPort))
		} else {
			mongoDbURI = fmt.Sprintf("mongodb://%s:%s/", mongoDbHost, mongoDbPort)
			logging.Logger.WithFields(logging.LogFields{
				"uri": mongoDbURI,
			}).Info("The database connection URI has been set")
		}
	}
	if len(errors) > 0 {
		logging.Logger.WithField("errors", errors).Error("The persistence environment is not properly set")
	}
	return errors
}

func initConnection() []string {
	errors := []string{}
	var err error
	ctx, cancel := GetContext()
	defer cancel()
	mongoDbClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoDbURI))
	if err != nil {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		logging.Logger.WithFields(logging.LogFields{
			"dbname": mongoDbName,
			"errors": errors,
		}).Error("Failed to connect to the database")
	} else {
		logging.Logger.WithFields(logging.LogFields{
			"dbname": mongoDbName,
		}).Info("Connected to the database")
	}

	return errors
}

func checkConnection() error {
	if mongoDbClient == nil {
		return fmt.Errorf("The database connection seems to be not initialized")
	}
	ctx, cancel := GetContext()
	defer cancel()
	if err := mongoDbClient.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}

// GetContext Retrieve Database ctx
func GetContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	return ctx, cancel
}

// GetCollection Retrieve Collection
func GetCollection(p Persistable) *mongo.Collection {
	return mongoDbClient.Database(mongoDbName).Collection(p.EntityName())
}

type Persistable interface {
	Id() string
	ResetId(id string) Persistable
	PlanId() string
	EntityName() string
	FromBson(sr *mongo.SingleResult) Persistable
}

func ToBson(p Persistable) []byte {
	p = p.ResetId("")
	bsonSolution, err := bson.Marshal(p)
	if err != nil {
		logging.Logger.WithFields(logging.LogFields{
			"error": err.Error(),
		}).Error(fmt.Sprintf("Cant bsonMarshal %s", p.EntityName()))
	}
	return bsonSolution
}

// FindOne: Find one entity
func FindOne(p Persistable) *Persistable {
	id, err := primitive.ObjectIDFromHex(p.Id())
	if err != nil {
		logging.Logger.Error(fmt.Sprintf("Can't retrieve HEX ObjectId : %s", p.Id()))
		return nil
	}
	ctx, cancel := GetContext()

	sr := GetCollection(p).FindOne(ctx, bson.M{"_id": id, "planid": p.PlanId()})
	if sr.Err() != nil {
		cancel()
		logging.Logger.WithFields(logging.LogFields{
			"stacktrace": sr.Err().Error(),
		}).Error(fmt.Sprintf("Can't retrieve %s with id : %s and planId : %s", p.EntityName(), p.Id(), p.PlanId()))
		return nil
	}
	instance := p.FromBson(sr)
	return &instance
}

// InsertOne : Insert one entity
func InsertOne(p Persistable) string {
	bson := ToBson(p)
	ctx, cancel := GetContext()
	res, errCol := GetCollection(p).InsertOne(ctx, bson)
	if errCol != nil {
		cancel()
		logging.Logger.WithFields(logging.LogFields{
			"collection": p.EntityName(),
			"id":         p.Id(),
			"planId":     p.PlanId(),
			"error":      errCol.Error(),
		}).Error("Can't get collection instance")
		return ""
	}
	logging.Logger.WithFields(logging.LogFields{"id": p.Id(), "planId": p.PlanId()}).Info(fmt.Sprintf("%s persisted", p.EntityName()))
	return res.InsertedID.(primitive.ObjectID).Hex()
}

// ReplaceOne : Replace one entity
func ReplaceOne(p Persistable) bool {
	id, err := primitive.ObjectIDFromHex(p.Id())
	if err != nil {
		logging.Logger.Error(fmt.Sprintf("Can't retrieve HEX ObjectId : %s", p.Id()))
		return false
	}
	toBson := ToBson(p)
	ctx, cancel := GetContext()
	_, errCol := GetCollection(p).ReplaceOne(ctx, bson.M{"_id": id, "planid": p.PlanId()}, toBson)
	if errCol != nil {
		cancel()
		logging.Logger.WithFields(logging.LogFields{
			"collection": p.EntityName(),
			"id":         p.Id(),
			"planId":     p.PlanId(),
			"error":      errCol.Error(),
		}).Error("Can't get collection instance")
		return false
	}
	logging.Logger.WithFields(logging.LogFields{"id": p.Id(), "planId": p.PlanId()}).Info(fmt.Sprintf("%s persisted", p.EntityName()))
	return true
}

// DeleteOne : Delete one entity
func DeleteOne(p Persistable) bool {
	id, err := primitive.ObjectIDFromHex(p.Id())
	if err != nil {
		logging.Logger.Error(fmt.Sprintf("Can't retrieve HEX ObjectId : %s", p.Id()))
		return false
	}
	ctx, cancel := GetContext()
	dr, err := GetCollection(p).DeleteOne(ctx, bson.M{"_id": id, "planid": p.PlanId()})
	if err != nil {
		cancel()
		logging.Logger.WithFields(logging.LogFields{
			"stacktrace": err.Error(),
		}).Error(fmt.Sprintf("Can't delete %s with id : %s", p.EntityName(), p.Id()))
	}

	return dr.DeletedCount == 1
}
