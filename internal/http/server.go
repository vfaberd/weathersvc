package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"weathersvc/internal/core"
)

type Server struct {
	router *http.ServeMux
	logger *log.Logger

	storage core.WeatherStorage
}

func NewServer(logger *log.Logger, storage core.WeatherStorage) *Server {
	return &Server{
		router:  http.NewServeMux(),
		logger:  logger,
		storage: storage,
	}
}

func (s *Server) Start(ctx context.Context, addr string) {
	srv := http.Server{
		Addr:    addr,
		Handler: s,
	}
	go func() {
		<-ctx.Done()
		s.logger.Println("Closing http server")
		if err := srv.Close(); err != nil {
			s.logger.Printf("close server:, %v", err)
		}
	}()
	s.registerRoutes()
    if err := srv.ListenAndServe(); err != nil {
        s.logger.Println(err)
    }
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.router.HandleFunc("/weather/latest", s.handleLatest)
}

func (s *Server) handleLatest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetLatest(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGetLatest(w http.ResponseWriter, r *http.Request) {
	loc := r.URL.Query().Get("location")
	if loc == "" {
		encodeError(errors.New("missing required query param: location"), http.StatusBadRequest, w)
		return
	}
	weather, err := s.storage.GetLatest(r.Context(), loc)
	if err != nil {
		encodeError(err, http.StatusInternalServerError, w)
		return
	}
	encodeResponse(w, &weather)
}

func encodeResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		encodeError(err, http.StatusInternalServerError, w)
		return
	}
}

func encodeError(err error, status int, w http.ResponseWriter) {
	resp := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&resp)
}

