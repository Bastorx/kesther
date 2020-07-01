package persistence

import (
	"context"
	"fmt"
	"os"
	"time"

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
	mongoDbServiceHost string
	mongoDbPort        string
	mongoDbURI         string
	mongoDbClient      *mongo.Client
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

func init() {
	if errors := initEnv(); len(errors) == 0 {
		initConnection()
	}
}

func initEnv() []string {
	errors := []string{}
	envDbServiceHost := "MONGODB_SERVICE_HOST"
	mongoDbServiceHost = os.Getenv(envDbServiceHost)
	if mongoDbServiceHost == "" {
		errors = append(errors, fmt.Sprintf("%s is not set", envDbServiceHost))
	} else {
		envDbPort := "MONGODB_PORT"
		mongoDbPort = os.Getenv(envDbPort)
		if mongoDbPort == "" {
			errors = append(errors, fmt.Sprintf("%s is not set", envDbPort))
		} else {
			mongoDbURI = fmt.Sprintf("mongodb://%s:%s/", mongoDbServiceHost, mongoDbPort)
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
	ctx, cancel := getContext()
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

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), contextTimeout)
}

func checkConnection() error {
	if mongoDbClient == nil {
		return fmt.Errorf("The database connection seems to be not initialized")
	}
	ctx, cancel := getContext()
	defer cancel()
	if err := mongoDbClient.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}
