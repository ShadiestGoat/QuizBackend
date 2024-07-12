package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shadiestgoat/log"
	router "whotfislucy.com/http"
)

func init() {
	log.Init(log.NewLoggerPrint())
}

func main() {
	conf := &opts{}
	sections := Load(conf)

	port := "3000"

	if p := os.Getenv("PORT"); p != "" {
		p = port
	}

	log.Success("Running http server on localhost:%v", port)
	hServer := &http.Server{Addr: ":" + port, Handler: router.Router(sections)}

	go func(s *http.Server) {
		log.FatalIfErr(s.ListenAndServe(), "running http server")
	}(hServer)

	cancel := make(chan os.Signal, 2)
	signal.Notify(cancel, os.Interrupt)

	<- cancel

	ctx, cancelCtx := context.WithTimeout(context.Background(), 25 * time.Second)
	defer cancelCtx()

	err := hServer.Shutdown(ctx)
	log.FatalIfErr(err, "closing http server")
}
