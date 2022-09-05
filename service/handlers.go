package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
			parsedAnswer, err := parseAnswerPayload(w, r)
			if err != nil {
				JSON(w, http.StatusBadRequest, errorResponse{
					Reason: err.Error(),
				})
				return
			}

			id, err := srv.useCase.createAnswer(*parsedAnswer)
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
			parsedAnswer, err := parseAnswerPayload(w, r)
			if err != nil {
				JSON(w, http.StatusBadRequest, errorResponse{
					Reason: err.Error(),
				})
				return
			}

			err = srv.useCase.updateAnswer(*parsedAnswer)
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

	switch r.Method {
	case http.MethodGet:
		{
			answer, err := srv.useCase.db.getAnswerByKey(vars["key"]) // ADAPT TO USECASE
			if err != nil {
				JSON(w, http.StatusNotFound, errorResponse{
					Reason: err.Error(),
				})
				return
			}
			JSON(w, http.StatusOK, answer)

		}
	case http.MethodDelete:
		{
			err := srv.useCase.deleteAnswer(vars["key"])
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

func (srv *Service) eventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	events, err := srv.useCase.db.getEventsHistory(vars["key"]) // ADAPT TO USECASE
	if err != nil {
		JSON(w, http.StatusNotFound, errorResponse{
			Reason: err.Error(),
		})
		return
	}
	JSON(w, http.StatusOK, events)
}

// utils
func parseAnswerPayload(w http.ResponseWriter, r *http.Request) (*Answer, error) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var answer Answer
	err = json.Unmarshal(data, &answer)
	if err != nil {
		log.Println("The payload is not valid JSON")
		return nil, err
	}
	return &answer, nil
}
