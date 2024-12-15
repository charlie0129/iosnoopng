package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type dataType struct {
	ProcessStat *ProcessStat
}

//go:embed dist
var static embed.FS

func httpServer() {
	f, err := fs.Sub(static, "dist")
	if err != nil {
		logrus.Fatalf("failed to embed static files: %v", err)
	}

	httpStaticFS := http.FileServerFS(f)

	http.HandleFunc("/api/stats", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		stats := processStat.GetMeta()
		err := json.NewEncoder(w).Encode(stats)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}))

	http.HandleFunc("/api/stats/{exec}", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		exec := r.PathValue("exec")
		stats := processStat.GetDetailsByProcess(exec)
		err := json.NewEncoder(w).Encode(stats)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}))

	http.HandleFunc("/assets/", WithLogging(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
		httpStaticFS.ServeHTTP(w, r)
	}))

	http.HandleFunc("/", WithLogging(func(w http.ResponseWriter, r *http.Request) {
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
