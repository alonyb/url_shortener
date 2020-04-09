package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	h "github.com/alonyb/url_shortener/api"
	mr "github.com/alonyb/url_shortener/repository/mongodb"
	rr "github.com/alonyb/url_shortener/repository/redis"
	"github.com/alonyb/url_shortener/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", c)
	}()
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortener.RedirectRepository {
	switch "mongo" {
	case "redis":
		//redisUrl := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepository("redis://localhost:6379")
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		//mongoURL := os.Getenv("MONGO_URL")
		//mongodb := os.Getenv("MONGO_DB")
		//mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mr.NewMongoRepository("mongodb://localhost:27017", "shortener", 30)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
