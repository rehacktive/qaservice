package service

import (
	"errors"
	"fmt"
	"log"
)

type AnswerUseCase struct {
	db AnswersDb
}

func (u AnswerUseCase) createAnswer(a Answer) (string, error) {
	// check if key already exists
	existingAnswer, _ := u.db.getAnswerByKey(a.Key)
	if existingAnswer != nil {
		fmt.Println(existingAnswer)
		return "", errors.New("key already exists, conflict")
	}
	// otherwise
	// store the answer
	id, err := u.db.createAnswer(a)
	if err != nil {
		return "", err
	}
	log.Println("created answer with id ", id)
	// and add the event
	err = u.db.storeEvent(Event{
		Type: Create,
		Data: a,
	})
	if err != nil {
		return "", err
	}
	log.Println("created event with id ", id)
	return id, nil
}

func (u AnswerUseCase) updateAnswer(a Answer) error {
	// check if key already exists
	existingAnswer, err := u.db.getAnswerByKey(a.Key)
	if existingAnswer == nil {
		return errors.New("error looking for the existing answer: not found")
	}
	if err != nil {
		return fmt.Errorf("error looking for the existing answer: %v", err)
	}
	// otherwise
	// update the answer
	err = u.db.updateAnswer(existingAnswer.Id.Hex(), a)
	if err != nil {
		return err
	}
	log.Println("updated answer with id ", existingAnswer.Id.Hex())
	// and add the event
	a.Id = existingAnswer.Id
	err = u.db.storeEvent(Event{
		Type: Update,
		Data: a,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u AnswerUseCase) deleteAnswer(key string) error {
	// check if key already exists
	existingAnswer, err := u.db.getAnswerByKey(key)
	if existingAnswer == nil {
		return errors.New("error looking for the existing answer: not found")
	}
	if err != nil {
		return fmt.Errorf("error looking for the existing answer: %v", err)
	}
	return u.db.deleteAnswer(existingAnswer.Id.Hex())
}