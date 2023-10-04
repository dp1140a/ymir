package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"ymir/front"
	"ymir/pkg"
	"ymir/pkg/api"
	"ymir/pkg/api/model"
	"ymir/pkg/api/printer"
	"ymir/pkg/logger/httplogger"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	log "github.com/sirupsen/logrus"
)

const (
	_API_VERSION = "/v1"
)

type Server struct {
	Config     *ServerConfig
	Router     *chi.Mux
	Handlers   []api.HandlerIFace
	tokenAuth  *jwtauth.JWTAuth
	HttpLogger *log.Logger
}

func NewServer() (*Server, error) {
	s := new(Server)
	s.Config = NewServerConfig()
	s.Router = chi.NewRouter()

	//TODO: refactor to pass in algorithm for acceptable algorithms reference: https://github.com/lestrrat-go/jwx/blob/v2/jwa/signature_gen.go
	s.tokenAuth = jwtauth.New("HS256", []byte(s.Config.JWTSecret), nil)

	if s.Config.HttpLogConfig.Enabled == true {
		log.Info("Http Request Logger Enabled.  Initializing")
		s.HttpLogger = log.New()
		s.HttpLogger.Formatter = &log.JSONFormatter{
			// disable, as we set our own
			DisableTimestamp: false,
		}
		s.Router.Use(httplogger.NewStructuredLogger(s.HttpLogger, s.Config.HttpLogConfig))
	}

	//Add and init middleware
	s.Router.Use(middleware.Recoverer)
	//s.Router.Use(middleware.RedirectSlashes)
	s.Router.Use(chiprometheus.NewMiddleware(pkg.APP_NAME))
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Compress(5, "text/html", "application/json"))

	if s.Config.EnableCORS {
		log.Info("Enabling CORS")
		s.Router.Use(corsConfig().Handler)
	}

	//Init Handlers and Services
	s.Router.Mount(fmt.Sprintf("%s/debug", _API_VERSION), middleware.Profiler())
	s.Handlers = append(s.Handlers, model.NewModelHandler())
	s.Handlers = append(s.Handlers, printer.NewPrinterHandler())

	//Append the base and static handlers Last
	s.Handlers = append(s.Handlers, api.NewBaseHandler(s.HttpLogger, s.Router))
	s.Router.Mount(_API_VERSION, s.registerRoutes()) // base

	s.Router.Mount("/", front.StaticHandler("/"))
	//s.Router.Handle("/admin", frontend.SvelteKitHandler("/admin")) //static

	return s, nil
}

func (s *Server) ServeAPI(ctx context.Context) {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%v", s.Config.Hostname, s.Config.Port),
		Handler: s.Router,
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error(err)
		}
		log.Info("Server Shutting down")
		close(done)
	}()

	sChan := make(chan string)
	go func() {
		transport := <-sChan
		log.Infof("Listening %s on %s:%v", transport, s.Config.Hostname, s.Config.Port)
		close(sChan)
	}()

	if s.Config.UseHttps {
		_, err := tls.LoadX509KeyPair(s.Config.TLSCert, s.Config.TLSKey)
		if err != nil { // Fail if we cant load certs
			log.Error("failed to load x509 key pair", err)
			log.Fatal("Stopping")
		}
		srv.TLSConfig = NewTLSConfig(s.Config)
		sChan <- "https"
		if err := srv.ListenAndServeTLS(s.Config.TLSCert, s.Config.TLSKey); err != http.ErrServerClosed {
			log.Error(err)
		}

	} else {
		sChan <- "http"
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error(err)
		}
	}
	<-done
}

func (s *Server) registerRoutes() chi.Router {
	r := chi.NewRouter()
	for _, handler := range s.Handlers {
		log.Info("Initializing Service ", handler.GetPrefix())
		r.Route(handler.GetPrefix(), func(r chi.Router) {
			for _, route := range handler.GetRoutes() {
				log.Info("Adding route ", route.Name)
				if route.Protected { //Protected Route
					r.Group(func(r chi.Router) {
						// Seek, verify and validate JWT tokens
						r.Use(jwtauth.Verifier(s.tokenAuth))
						r.Use(jwtauth.Authenticator)
						r.Method(route.Method, route.Pattern, route.HandlerFunc)
					})
				} else { // Public Route
					r.Method(route.Method, route.Pattern, route.HandlerFunc)
				}
			}
		})
	}
	return r
}

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Sveltekit-Action"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
		//OptionsPassthrough: true,
		//Debug: true,
	})
}
