package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	// "io/ioutil"

	"github.com/asaskevich/govalidator"
	"github.com/asdine/storm"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/logger"
)

type Server struct {
	bind   string
	config Config
	router *httprouter.Router

	logger *logger.Logger
}

func (s *Server) ShortenHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		url := r.FormValue("url")
		log.Printf("Request url: %s\n", url)

		if !govalidator.IsURL(url) {
			http.Error(w, "Not a valid URL", http.StatusBadRequest)
			return
		}

		m, _ := regexp.Match("^http(s?)://", []byte(url))
		if !m {
			url = "http://" + url
		}

		u, err := NewURL(url)
		if err != nil {
			log.Printf("error creating new url: %s", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		shortURL := fmt.Sprintf("http://%s/%s", s.config.baseURL, u.ID)

		w.Write([]byte(shortURL))
	}
}

func (s *Server) RedirectHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var u URL

		id := p.ByName("id")
		if id == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err := db.One("ID", id, &u)
		if err != nil && err == storm.ErrNotFound {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("error looking up %s for redirect: %s", id, err)
			http.Error(w, "Iternal Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, u.URL, http.StatusFound)
	}
}

func (s *Server) ListenAndServe() {
	log.Fatal(
		http.ListenAndServe(
			s.bind,
			s.logger.Handler(
				s.router,
			),
		),
	)
}

func (s *Server) initRoutes() {
	s.router.POST("/", s.ShortenHandler())
	s.router.GET("/:id", s.RedirectHandler())
}

func NewServer(bind string, config Config) *Server {
	server := &Server{
		bind:   bind,
		config: config,
		router: httprouter.New(),

		logger: logger.New(logger.Options{
			Prefix:               "shortener",
			RemoteAddressHeaders: []string{"X-Forwarded-For", "Content-Type"},
			OutputFlags:          log.LstdFlags,
		}),
	}

	server.initRoutes()

	return server
}
