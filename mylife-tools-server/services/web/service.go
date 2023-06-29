package web

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"mylife-tools-server/config"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	sio "mylife-tools-server/services/io"
	"net/http"
	"os"
	"path"
	"strings"
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

	service.mux.Handle("/socket.io/", sio.GetHandler())

	if arg != nil {
		fs := arg.(fs.FS)

		list, err := getAllFilenames(fs)
		if err != nil {
			return err
		}

		for _, name := range list {
			logger.Tracef("Serving embedded file '%s'", name)
		}

		service.mux.Handle("/", spaHandler(fs))
	}

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

func spaHandler(fsys fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}
		upath = path.Clean(upath)[1:] // remove leading /

		// Serve static file, if not found serve index.html
		file, err := openFile(fsys, r)
		if err != nil {
			logger.WithError(err).WithField("path", upath).Error("Error serving SPA")
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			logger.WithError(err).WithField("path", upath).Error("Error serving SPA")
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		content, ok := file.(io.ReadSeeker)
		if !ok {
			logger.WithField("path", upath).Error("Error serving SPA: file is not seekable")
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.WithField("path", upath).Tracef("Serving path")

		http.ServeContent(w, r, stat.Name(), stat.ModTime(), content)
	})
}

func openFile(fsys fs.FS, r *http.Request) (fs.File, error) {
	const defaultFileName = "index.html"

	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	upath = path.Clean(upath)[1:] // remove leading /

	if upath == "" {
		return fsys.Open(defaultFileName)
	}

	file, err := fsys.Open(upath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return fsys.Open(defaultFileName)
	}

	return file, err
}

func getAllFilenames(efs fs.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}
