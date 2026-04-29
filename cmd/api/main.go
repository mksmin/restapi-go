package main

import (
	"RestAPIGo/internal/database"
	"RestAPIGo/internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:123456test@localhost/go_restapi_db?sslmode=disable"
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8083"
	}
	log.Printf(
		"Starting server on port %s",
		serverPort,
	)
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Database connected")

	taskStore := database.NewTaskStore(db)
	handler := handlers.NewHandlers(taskStore)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/tasks",
		methodHandler(
			handler.GetAllTasks,
			http.MethodGet,
		),
	)
	mux.HandleFunc(
		"/tasks/create",
		methodHandler(
			handler.CreateTask,
			http.MethodPost,
		),
	)
	mux.HandleFunc(
		"/tasks/",
		taskIDHandler(handler),
	)

	loggedMux := loggingMiddleware(mux)

	serverAddr := ":" + serverPort

	err = http.ListenAndServe(
		serverAddr,
		loggedMux,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func methodHandler(
	handlerFunc http.HandlerFunc,
	allowedMethod string,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != allowedMethod {
			http.Error(
				w,
				http.StatusText(
					http.StatusMethodNotAllowed,
				),
				http.StatusMethodNotAllowed,
			)
		}
		handlerFunc(w, r)
	}
}

func taskIDHandler(
	handler *handlers.Handlers,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTask(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(
				w,
				http.StatusText(
					http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed,
			)
		}
	}
}

func loggingMiddleware(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			log.Printf(
				"%s %s %s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
			)
			next.ServeHTTP(w, r)
		})
}
