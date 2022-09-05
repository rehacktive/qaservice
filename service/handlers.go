package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

type createdResponse struct {
	Id string `json:"id"`
}

type errorResponse struct {
	Reason string `json:"error_message"`
}

func (srv *Service) answersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		{
			id, err := srv.useCase.createAnswer(Answer{
				Key:   "question1",
				Value: "this is the answer",
			})
			if err != nil {
				JSON(w, http.StatusConflict, errorResponse{
					Reason: err.Error(),
				})
				return
			}
			JSON(w, http.StatusCreated, createdResponse{
				Id: id,
			})
		}
	case http.MethodPut:
		{
			err := srv.useCase.updateAnswer(Answer{
				Key:   "question1",
				Value: "this is the updated answer",
			})
			if err != nil {
				JSON(w, http.StatusConflict, errorResponse{
					Reason: err.Error(),
				})
				return
			}
			JSON(w, http.StatusAccepted, nil)

		}
	case http.MethodDelete:
		{
			err := srv.useCase.deleteAnswer("question1")
			if err != nil {
				JSON(w, http.StatusConflict, errorResponse{
					Reason: err.Error(),
				})
				return
			}
			JSON(w, http.StatusAccepted, nil)

		}
	}
}

func (srv *Service) answersFetchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	answer, err := srv.useCase.db.getAnswerByKey(vars["key"])
	if err != nil {
		JSON(w, http.StatusNotFound, errorResponse{
			Reason: err.Error(),
		})
		return
	}
	JSON(w, http.StatusOK, answer)
}

func (srv *Service) eventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	events, err := srv.useCase.db.getEventsHistory(vars["key"])
	if err != nil {
		JSON(w, http.StatusNotFound, errorResponse{
			Reason: err.Error(),
		})
		return
	}
	JSON(w, http.StatusOK, events)
}
