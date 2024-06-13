package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbHost = "localhost"
	dbPort = "27017"
	dbUser = "admin"
	dbPass = "password"
)

type Mongodb struct {
	client *mongo.Client
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func ConnectMongodb() *Mongodb {
	clientOptions := options.Client().ApplyURI("mongodb://" + dbHost + ":" + dbPort)
	clientOptions.SetAuth(options.Credential{
		Username: dbUser,
		Password: dbPass,
	})
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error while connecting to MongoDB: ", err)
		return nil
	}
	log.Println("Connected to MongoDB!")
	return &Mongodb{client: client}
}

func (m *Mongodb) InsertLog(entry LogEntry) error {
	collection := m.client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}
