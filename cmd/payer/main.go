package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	oapi "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rallinator7/fetch-demo/internal/environment"
	"github.com/rallinator7/fetch-demo/internal/logger"
	"github.com/rallinator7/fetch-demo/internal/payer"
	"github.com/rallinator7/fetch-demo/internal/payer/api"
	"github.com/rallinator7/fetch-demo/internal/pubsub"
	"golang.org/x/sync/errgroup"
)

type Logger interface {
	Info(...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

func main() {
	logs, err := logger.New()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// setup os signal trigger for shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// setup error group for shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)

	serve, server := ServeHTTP(logs)

	logs.Info("starting http server...")

	errGroup.Go(serve)

	// wait for shutdown signals
	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	// will trigger errGroup to shutdown if os signal is what caused shutdown
	cancel()

	logs.Info("attempting to shutdown servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		logs.Fatal(err)
	}

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			logs.Fatal("graceful shutdown timed out.. forcing exit.")
		}
	}()

	// wait for servers to shutdown
	err = errGroup.Wait()
	if err != nil {
		logs.Fatal(err)
		os.Exit(2)
	}

}

func ServeHTTP(logs Logger) (func() error, *http.Server) {
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// environment variables
	pubsubServer := environment.GetVariableFatal("PUBSUB_SERVER", logs)
	payerStream := environment.GetVariableFatal("PAYER_STREAM", logs)
	port, err := strconv.Atoi(environment.GetVariableFatal("PORT", logs))
	if err != nil {
		logs.Fatalf(err.Error())
	}

	streamer, err := pubsub.NewJetstream(pubsubServer)
	if err != nil {
		logs.Fatalf(err.Error())
	}

	pub := pubsub.NewPublisher(streamer)
	if err != nil {
		log.Fatalf(err.Error())
	}

	repo := payer.NewRepository()
	publisher := payer.NewPublisher(&pub, payerStream)

	control := payer.NewController(&repo, &publisher, logs)
	handler := payer.NewHandler(&control)
	router := chi.NewRouter()

	router.Use(logger.Middleware(logs))
	router.Use(middleware.Recoverer)
	router.Use(oapi.OapiRequestValidator(swagger))

	api.HandlerFromMux(&handler, router)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: router,
	}

	return serve(server), server
}

func serve(server *http.Server) func() error {
	return func() error {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	}
}
