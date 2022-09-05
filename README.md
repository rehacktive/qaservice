# QuestionAnswers Service

## Run it

Run the entire project (service + mongo + mongo express):

```sh
docker-compose up -d --build qaservice
```

or just run Mongo on docker:

```sh
docker run -d -p 27017:27017 --name test-mongo mongo:latest
```
and then just:
```sh
go run main.go
```

## Manual examples (using curl)

```sh
# create
> curl -X POST -d '{"key":"name","value":"john"}' http://localhost:8880/answers | jq .
{
  "id": "631617b8481428005dd1a982"
}

# get
> curl http://localhost:8880/answers/name | jq .
{
  "id": "631617b8481428005dd1a982",
  "key": "name",
  "value": "john"
}

# error on conflict
> curl -X POST -d '{"key":"name","value":"john"}' http://localhost:8880/answers | jq .
{
  "error_message": "key already exists, conflict"
}

# get events
> curl http://localhost:8880/events/name | jq .
[
  {
    "id": "631617b8481428005dd1a983",
    "event_type": "CREATE",
    "data": {
      "id": "631617b8481428005dd1a982",
      "key": "name",
      "value": "john"
    }
  }
]
# update
> curl -X PUT -d '{"key":"name","value":"jack"}' http://localhost:8880/answers | jq .

# fetch updated
> curl http://localhost:8880/answers/name | jq .
{
  "id": "631617b8481428005dd1a982",
  "key": "name",
  "value": "jack"
}
```

## A sample unit test

A sample unit test for the usecase has been added.

## Some integration tests

Some integration tests using docker compose with testcontainers are implemented.
You can find them in `service/integration_test.go` and run `TestIntegration` (it could take a while).
