package service

import "testing"

func TestCreateAnswerSuccess(t *testing.T) {
	usecase := AnswerUseCase{
		db: &MockedDatabase{
			currentAnswers: make(map[string]*Answer),
			events:         make([]Event, 0),
		},
	}

	id, err := usecase.createAnswer(Answer{
		Key:   "name",
		Value: "john",
	})

	if id != defaultId {
		t.Errorf("Expected %s Got %s", "1", id)
	}
	if err != nil {
		t.Errorf("Expected no errors Got %s", err)
	}

	storedAnswer, _ := usecase.getAnswerByKey("name")
	if storedAnswer.Value != "john" {
		t.Errorf("Expected %s Got %s", "jack", storedAnswer.Value)
	}
}

func TestCreateAnswerTwiceFail(t *testing.T) {
	usecase := AnswerUseCase{
		db: &MockedDatabase{
			currentAnswers: make(map[string]*Answer),
			events:         make([]Event, 0),
		},
	}

	_, _ = usecase.createAnswer(Answer{
		Key:   "name",
		Value: "john",
	})

	_, err := usecase.createAnswer(Answer{
		Key:   "name",
		Value: "john",
	})

	if err == nil {
		t.Errorf("Expected an error Got no errors")
	}
}

func TestUpdateAnswerSuccess(t *testing.T) {
	usecase := AnswerUseCase{
		db: &MockedDatabase{
			currentAnswers: make(map[string]*Answer),
			events:         make([]Event, 0),
		},
	}

	_, _ = usecase.createAnswer(Answer{
		Key:   "name",
		Value: "john",
	})

	err := usecase.updateAnswer(Answer{
		Key:   "name",
		Value: "jack",
	})
	if err != nil {
		t.Errorf("Expected no errors Got %s", err)
	}

	storedAnswer, _ := usecase.getAnswerByKey("name")
	if storedAnswer.Value != "jack" {
		t.Errorf("Expected %s Got %s", "jack", storedAnswer.Value)
	}
}

func TestUpdateAnswerNotExistingFailure(t *testing.T) {
	usecase := AnswerUseCase{
		db: &MockedDatabase{
			currentAnswers: make(map[string]*Answer),
			events:         make([]Event, 0),
		},
	}

	err := usecase.updateAnswer(Answer{
		Key:   "name",
		Value: "jack",
	})
	if err == nil {
		t.Errorf("Expected an error Got no errors")
	}
}

func TestDeleteAnswerSuccess(t *testing.T) {
	usecase := AnswerUseCase{
		db: &MockedDatabase{
			currentAnswers: make(map[string]*Answer),
			events:         make([]Event, 0),
		},
	}

	_, _ = usecase.createAnswer(Answer{
		Key:   "name",
		Value: "john",
	})

	err := usecase.deleteAnswer("name")
	if err != nil {
		t.Errorf("Expected no errors Got %s", err)
	}

	_, err = usecase.getAnswerByKey("name")
	if err == nil {
		t.Errorf("Expected an error Got no errors")
	}
}

func TestDeleteAnswerNotExistingFailure(t *testing.T) {
	usecase := AnswerUseCase{
		db: &MockedDatabase{
			currentAnswers: make(map[string]*Answer),
			events:         make([]Event, 0),
		},
	}

	err := usecase.deleteAnswer("name")

	if err == nil {
		t.Errorf("Expected an error Got no errors")
	}
}

func TestEventsHistory(t *testing.T) {
	usecase := AnswerUseCase{
		db: &MockedDatabase{
			currentAnswers: make(map[string]*Answer),
			events:         make([]Event, 0),
		},
	}

	_, _ = usecase.createAnswer(Answer{
		Key:   "name",
		Value: "john",
	})

	_ = usecase.updateAnswer(Answer{
		Key:   "name",
		Value: "jack",
	})

	events, err := usecase.getEventsHistory("name")
	if err != nil {
		t.Errorf("Expected no errors Got %s", err)
	}

	// check events
	if events[0].Type != Create {
		t.Errorf("Expected %v Got %v", Create, events[1].Type)
	}
	if events[0].Data.Value != "john" {
		t.Errorf("Expected %v Got %v", "john", events[0].Data.Value)
	}
	if events[1].Type != Update {
		t.Errorf("Expected %v Got %v", Create, events[1].Type)
	}
	if events[1].Data.Value != "jack" {
		t.Errorf("Expected %v Got %v", "jack", events[0].Data.Value)
	}
}
