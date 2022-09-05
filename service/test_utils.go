package service

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mock database

const defaultId = "6315d3b328cb4edb2bff8204"

type MockedDatabase struct {
	currentAnswers map[string]*Answer
	events         []Event
}

func (m *MockedDatabase) createAnswer(a Answer) (string, error) {
	actualId, _ := primitive.ObjectIDFromHex(defaultId)
	a.Id = actualId
	m.currentAnswers[defaultId] = &a
	return defaultId, nil

}

func (m *MockedDatabase) updateAnswer(id string, a Answer) error {
	m.currentAnswers[id] = &a
	return nil
}

func (m *MockedDatabase) getAnswerByKey(key string) (*Answer, error) {
	for _, a := range m.currentAnswers {
		if a != nil && a.Key == key {
			return a, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockedDatabase) deleteAnswer(id string) error {
	m.currentAnswers[id] = nil
	return nil
}

func (m *MockedDatabase) getEventsHistory(key string) ([]Event, error) {
	return m.events, nil
}

func (m *MockedDatabase) storeEvent(event Event) error {
	m.events = append(m.events, event)
	return nil
}
