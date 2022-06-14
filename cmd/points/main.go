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
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rallinator7/fetch-demo/internal/environment"
	"github.com/rallinator7/fetch-demo/internal/logger"
	"github.com/rallinator7/fetch-demo/internal/points"
	"github.com/rallinator7/fetch-demo/internal/points/api"
	"github.com/rallinator7/fetch-demo/internal/points/payer"
	"github.com/rallinator7/fetch-demo/internal/points/user"
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

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf(err.Error())
	}

	pubsubServer := environment.GetVariableFatal("PUBSUB_SERVER", logs)
	payerStream := environment.GetVariableFatal("PAYER_STREAM", logs)
	payerQueue := environment.GetVariableFatal("PAYER_QUEUE", logs)
	userStream := environment.GetVariableFatal("USER_STREAM", logs)
	userQueue := environment.GetVariableFatal("USER_QUEUE", logs)
	port, err := strconv.Atoi(environment.GetVariableFatal("PORT", logs))
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

	// Repositories
	payerRepo := payer.NewRepository()
	userRepo := user.NewRepository()

	streamer, err := pubsub.NewJetstream(pubsubServer)
	if err != nil {
		logs.Fatalf(err.Error())
	}

	queuer := pubsub.NewSubscriber(streamer)
	payerSub := payer.NewSubscriber(&queuer, &payerRepo, logs)
	userSub := user.NewSubscriber(&queuer, &userRepo, logs)

	subscriberCtx, subscriberCancel := context.WithCancel(context.Background())
	defer subscriberCancel()

	serve, server := ServeHTTP(logs, &payerRepo, &userRepo, swagger, port)

	errGroup.Go(serve)
	errGroup.Go(payerSub.PayerAddedEvent(subscriberCtx, payerStream, payerQueue))
	errGroup.Go(userSub.UserAddedEvent(subscriberCtx, userStream, userQueue))

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

	serverCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = server.Shutdown(serverCtx)
	if err != nil {
		logs.Fatal(err)
	}

	subscriberCancel()

	go func() {
		<-serverCtx.Done()
		if serverCtx.Err() == context.DeadlineExceeded {
			logs.Fatal("graceful shutdown timed out.. forcing exit.")
		}
		<-subscriberCtx.Done()
		if subscriberCtx.Err() == context.DeadlineExceeded {
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

func ServeHTTP(logs logger.Logger, payerRepo *payer.Repository, userRepo *user.Repository, swagger *openapi3.T, port int) (func() error, *http.Server) {
	control := points.NewController(userRepo, payerRepo)
	handler := points.NewHandler(&control)
	router := chi.NewRouter()

	router.Use(logger.Middleware(logs))
	router.Use(middleware.Recoverer)
	router.Use(oapi.OapiRequestValidator(swagger))

	api.HandlerFromMux(&handler, router)

	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Handler: router}

	serverFunc := func() error {
		logs.Infow(fmt.Sprintf("starting http server on port %d", port))

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	}

	return serverFunc, server
}
