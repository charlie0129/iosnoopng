package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/sirupsen/logrus"
)

//go:embed dist
var static embed.FS

func httpServer() {
	f, err := fs.Sub(static, "dist")
	if err != nil {
		logrus.Fatalf("failed to embed static files: %v", err)
	}

	httpStaticFS := http.FileServerFS(f)

	// Expose the registered metrics at `/metrics` path.
	http.HandleFunc("GET /metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, false)
	})

	http.HandleFunc("GET /api/stats", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		stats := processStat.GetMeta()
		err := json.NewEncoder(w).Encode(stats)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}))

	http.HandleFunc("DELETE /api/stats", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		processStat.Clear()
		w.WriteHeader(http.StatusNoContent)
	}))

	http.HandleFunc("GET /api/stats/{exec}", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		exec := r.PathValue("exec")
		stats := processStat.GetDetailsByProcess(exec)
		err := json.NewEncoder(w).Encode(stats)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}))

	http.HandleFunc("DELETE /api/stats/{exec}", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		exec := r.PathValue("exec")
		processStat.DeleteByProcess(exec)
		w.WriteHeader(http.StatusNoContent)
		return
	}))

	http.HandleFunc("GET /api/stats-full", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		stats := processStat.GetAll()
		err := json.NewEncoder(w).Encode(stats)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}))

	http.HandleFunc("GET /assets/", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
		httpStaticFS.ServeHTTP(w, r)
	}))

	http.HandleFunc("GET /", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/"
		httpStaticFS.ServeHTTP(w, r)
	}))

	logrus.Infof("HTTP Server Listening on %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

func WithLogging(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	loggingFn := func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()

		hf := http.HandlerFunc(h)
		hf.ServeHTTP(rw, req)

		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"uri":      req.RequestURI,
			"method":   req.Method,
			"duration": duration,
		}).Debug("request completed")
	}

	return loggingFn
}
