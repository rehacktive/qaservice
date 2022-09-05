package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const maxPayloadSize = 1048576

func JSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
	}
}

func parseAnswerPayload(w http.ResponseWriter, r *http.Request) (*Answer, error) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, maxPayloadSize)
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
