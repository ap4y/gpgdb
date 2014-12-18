package main

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/bmizerany/pat"

	"github.com/ap4y/gpgdb/lib"
	"github.com/ap4y/gpgdb/v1"
)

type Server struct {
	Router *pat.PatternServeMux
	es     *lib.EntityStorage
	db     *lib.DB
}

type HandlerWithContext func(http.ResponseWriter, *lib.Request, lib.DBService)

func NewServer(es *lib.EntityStorage, db *lib.DB) *Server {
	server := &Server{pat.New(), es, db}
	server.setupRouting()
	return server
}

func (s *Server) setupRouting() {
	s.Router.Get("/v1/keys", s.handlerWithContext(v1.ListKeys))
	s.Router.Get("/v1/keys/:key", s.handlerWithContext(v1.GetKey))
	s.Router.Put("/v1/keys/:key", s.handlerWithContext(v1.PutKey))
	s.Router.Del("/v1/keys/:key", s.handlerWithContext(v1.DeleteKey))
}

func (s *Server) handlerWithContext(handler HandlerWithContext) http.Handler {
	internalHandler := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %s\n%s", err, debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		req, err := lib.AuthenticatedRequest(r, s.es)
		if err != nil {
			lib.ErrorJSON(w, err.Error(), http.StatusUnauthorized)
			return
		}

		handler(w, req, s.db)
	}

	return http.HandlerFunc(internalHandler)
}
