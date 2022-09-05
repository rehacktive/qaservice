package service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnswersDb interface {
	createAnswer(Answer) (string, error)    // store answer, return generated id or error
	updateAnswer(string, Answer) error      // update by id, with key/value, return error
	getAnswerByKey(string) (*Answer, error) // get answer by key
	deleteAnswer(string) error              // delete answer by id

	getEventsHistory(string) ([]Event, error) // get events by key
	storeEvent(Event) error                   // store event, return error
}

type EventType string

const (
	Create EventType = "CREATE"
	Update EventType = "UPDATE"
	Delete EventType = "DELETE"
)

type Answer struct {
	Id    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Key   string             `json:"key"`
	Value string             `json:"value"`
}

type Event struct {
	Id   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type EventType          `json:"event_type"`
	Data Answer             `json:"data"`
}

// implementation
const (
	database          = "answersdb"
	answersCollection = "answers"
	eventsCollection  = "events"
)

type MongoAnswersDb struct {
	client *mongo.Client
}

func InitDb(URI string) (*MongoAnswersDb, error) {
	log.Println("db: connecting to ", URI)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return &MongoAnswersDb{
		client: client,
	}, nil
}

func (m *MongoAnswersDb) createAnswer(a Answer) (string, error) {
	c := m.client.Database(database).Collection(answersCollection)

	res, err := c.InsertOne(context.TODO(), a)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m *MongoAnswersDb) updateAnswer(id string, a Answer) error {
	c := m.client.Database(database).Collection(answersCollection)

	mongoid, _ := primitive.ObjectIDFromHex(id)
	_, err := c.UpdateOne(
		context.TODO(),
		bson.M{"_id": mongoid},
		bson.D{
			{Key: "$set", Value: bson.M{"key": a.Key}},
			{Key: "$set", Value: bson.M{"value": a.Value}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (m *MongoAnswersDb) getAnswerByKey(key string) (*Answer, error) {
	c := m.client.Database(database).Collection(answersCollection)

	var answer Answer
	if err := c.FindOne(context.TODO(), bson.M{"key": key}).Decode(&answer); err != nil {
		return nil, err
	}
	return &answer, nil
}

func (m *MongoAnswersDb) deleteAnswer(id string) error {
	c := m.client.Database(database).Collection(answersCollection)

	mongoid, _ := primitive.ObjectIDFromHex(id)

	_, err := c.DeleteOne(context.TODO(), bson.M{"_id": mongoid})
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoAnswersDb) getEventsHistory(key string) ([]Event, error) {
	c := m.client.Database(database).Collection(eventsCollection)

	var events []Event
	contactCursor, err := c.Find(context.TODO(), bson.M{"data.key": key})
	if err != nil {
		return nil, err
	}

	if err = contactCursor.All(context.TODO(), &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (m *MongoAnswersDb) storeEvent(event Event) error {
	c := m.client.Database(database).Collection(eventsCollection)

	_, err := c.InsertOne(context.TODO(), event)
	if err != nil {
		return err
	}
	return nil
}
