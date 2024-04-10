package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
	receiver "github.com/tomekzakrzewski/lottery/number_receiver/client"
	checker "github.com/tomekzakrzewski/lottery/result_checker/client"
)

func main() {
	checkerClient := checker.NewHTTPClient("http://localhost:5000")
	receiverClient := receiver.NewHTTPClient("http://localhost:3000")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redis := NewRedisStore(redisClient)
	svc := NewResultAnnoucerService(*checkerClient, *receiverClient, redis)
	m := NewLogMiddleware(svc)
	srv := NewHttpTransport(m)

	r := chi.NewRouter()
	r.Get("/win/{hash}", srv.handleCheckResult)

	http.ListenAndServe(":6000", r)

}
