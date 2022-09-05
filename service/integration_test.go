package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
)

var jsonPayload = `{"key":"name","value":"john"}`
var jsonPayloadUpdate = `{"key":"name","value":"jack"}`

var httpClient http.Client

// Integration test with Mongo

func setupContainers() (string, error) {
	composeFilePaths := []string{"../docker-compose.yaml"}
	identifier := strings.ToLower(uuid.New().String())

	compose := testcontainers.NewLocalDockerCompose(composeFilePaths, identifier)
	execError := compose.
		WithCommand([]string{"up", "-d"}).
		Invoke()
	err := execError.Error
	if err != nil {
		fmt.Println("error on setupContainers:", err)
		os.Exit(-1)
	}
	return identifier, nil
}

func stopContainers(identifier string) error {
	composeFilePaths := []string{"../docker-compose.yaml"}

	compose := testcontainers.NewLocalDockerCompose(composeFilePaths, identifier)
	execError := compose.Down()
	err := execError.Error
	if err != nil {
		return fmt.Errorf("could not run compose file: %v - %v", composeFilePaths, err)
	}
	return nil
}

func TestIntegration(t *testing.T) {
	httpClient = http.Client{Timeout: time.Duration(5) * time.Second}

	id, err := setupContainers()
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	// do tests
	doTests(t)

	err = stopContainers(id)
	if err != nil {
		panic(err)
	}
}

func doTests(t *testing.T) {
	expected := parseResponse([]byte(jsonPayload))

	// create an answer
	code, _, err := call(http.MethodPost, "http://localhost:8880/answers", jsonPayload)
	checkErr(t, err)
	checkResponseCode(t, "create answer", code, http.StatusCreated)

	// get the answer just created
	code, response, err := call(http.MethodGet, "http://localhost:8880/answers/name", "")
	checkErr(t, err)
	checkResponseCode(t, "get answer", code, http.StatusOK)
	answer := parseResponse(response)
	if answer.Key != expected.Key {
		t.Errorf("Get: Expected %v. Got %v\n", expected.Key, answer.Key)
	}
	if answer.Value != expected.Value {
		t.Errorf("Get: Expected %v. Got %v\n", expected.Value, answer.Value)
	}

	// create conflict
	code, _, err = call(http.MethodPost, "http://localhost:8880/answers", jsonPayload)
	checkErr(t, err)
	checkResponseCode(t, "create answer conflict", code, http.StatusConflict)

	// update
	code, _, err = call(http.MethodPut, "http://localhost:8880/answers", jsonPayloadUpdate)
	checkErr(t, err)
	checkResponseCode(t, "update answer", code, http.StatusAccepted)
	// now fetch again
	code, response, err = call(http.MethodGet, "http://localhost:8880/answers/name", "")
	checkErr(t, err)
	checkResponseCode(t, "get answer", code, http.StatusOK)
	answer = parseResponse(response)
	if answer.Value != "jack" {
		t.Errorf("Get: Expected %v. Got %v\n", "jack", answer.Value)
	}
}

func call(method string, URL string, payload string) (code int, response []byte, err error) {
	req, err := http.NewRequest(method, URL, strings.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return resp.StatusCode, body, nil
}

// utils

func checkResponseCode(t *testing.T, testName string, expected, actual int) {
	if expected != actual {
		t.Errorf("%v: Expected response code %d. Got %d\n", testName, expected, actual)
	}
}

func parseResponse(data []byte) Answer {
	var answer Answer
	_ = json.Unmarshal(data, &answer)
	return answer
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
