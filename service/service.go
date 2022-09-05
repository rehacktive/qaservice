package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Service struct {
	router  *mux.Router
	useCase *AnswerUseCase
}

func InitService(dbInstance AnswersDb) *Service {
	return &Service{
		router: mux.NewRouter(),
		useCase: &AnswerUseCase{
			db: dbInstance,
		},
	}
}

func (srv Service) Start() {
	srv.router.HandleFunc("/answers", srv.answersHandler).Methods(http.MethodPost, http.MethodPut, http.MethodDelete)
	srv.router.HandleFunc("/answers/{key}", srv.answersFetchHandler).Methods(http.MethodGet)
	srv.router.HandleFunc("/events/{key}", srv.eventsHandler).Methods(http.MethodGet)

	server := &http.Server{Addr: ":8880", Handler: srv.router}

	go server.ListenAndServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("unable to stop gracefully server: %v", err)
	}

	log.Println("service stopped.")
}
