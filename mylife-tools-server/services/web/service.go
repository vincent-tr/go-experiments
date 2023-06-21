package database

import (
	"context"
	"mylife-tools-server/config"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/io"
	"net/http"
	"sync"
)

var logger = log.CreateLogger("mylife:server:web")

func init() {
	services.Register(&webService{})
}

type webServerConfig struct {
	Address string `mapstructure:"address"`
}

type webService struct {
	server   *http.Server
	exitDone *sync.WaitGroup
	mux      *http.ServeMux
}

// https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve

func (service *webService) Init(arg interface{}) error {
	webServerConfig := webServerConfig{}
	config.BindStructure("webServer", &webServerConfig)

	service.exitDone = &sync.WaitGroup{}
	service.exitDone.Add(1)

	service.mux = http.NewServeMux()

	service.server = &http.Server{
		Addr:    webServerConfig.Address,
		Handler: service.mux,
	}

	service.mux.Handle("/socket.io/", io.GetHandler())

	// TODO
	// service.mux.HandleFunc("/")

	go func() {
		defer service.exitDone.Done()

		// always returns error. ErrServerClosed on graceful close
		if err := service.server.ListenAndServe(); err != http.ErrServerClosed {
			logger.WithError(err).Error("ListenAndServe error")
		}
	}()

	logger.WithField("address", webServerConfig.Address).Info("Listening")

	return nil
}

func (service *webService) Terminate() error {
	if err := service.server.Shutdown(context.TODO()); err != nil {
		return err
	}

	service.exitDone.Wait()

	logger.Info("Stopped")

	return nil
}

func (service *webService) ServiceName() string {
	return "web"
}

func (service *webService) Dependencies() []string {
	return []string{"io"}
}

func getService() *webService {
	return services.GetService[*webService]("web")
}
