package transport

import (
	"errors"
	"go-rummi-q-server/internal/config"
	"net/http"
	"strconv"
	"time"
)

func StartServer(httpConfig config.HTTPConfig) error {
	mux := NewRouter()

	addr := httpConfig.ListenHost + ":" + strconv.Itoa(httpConfig.ListenPort)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(httpConfig.Timeout) * time.Second,
		WriteTimeout: time.Duration(httpConfig.Timeout) * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
