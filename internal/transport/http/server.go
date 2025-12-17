package transport

import (
	"errors"
	"go-rummi-q-server/internal/config"
	"net/http"
	"strconv"
	"time"
)

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func StartServer(httpConfig config.HTTPConfig) error {
	mux := NewRouter()

	securedMux := securityHeaders(mux)

	addr := httpConfig.ListenHost + ":" + strconv.Itoa(httpConfig.ListenPort)

	srv := &http.Server{
		Addr:         addr,
		Handler:      securedMux,
		ReadTimeout:  time.Duration(httpConfig.Timeout) * time.Second,
		WriteTimeout: time.Duration(httpConfig.Timeout) * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
