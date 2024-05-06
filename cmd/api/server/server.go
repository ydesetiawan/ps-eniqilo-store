package server

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"ps-eniqilo-store/internal/shared"
	bhandler "ps-eniqilo-store/pkg/base/handler"
	"time"

	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

type Server struct {
	baseHandler *bhandler.BaseHTTPHandler
	router      *muxtrace.Router
	port        int
}

func NewServer(
	bHandler *bhandler.BaseHTTPHandler,
	port int,
) Server {
	return Server{
		baseHandler: bHandler,
		router:      muxtrace.NewRouter(muxtrace.WithServiceName(shared.ServiceName)),
		port:        port,
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log details about the incoming request
		log.Printf("[%s] %s %s", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)

		// Create a custom response writer to intercept the response status code
		crw := &customResponseWriter{ResponseWriter: w, buf: bytes.NewBuffer(nil)}

		// Call the next handler in the chain
		next.ServeHTTP(crw, r)

		// Log details about the outgoing response
		log.Printf("[%s] Response Status: %d, Response Body: %s", time.Now().Format("2006-01-02 15:04:05"), crw.status, crw.buf.String())
	})
}

// Custom ResponseWriter to intercept the response status code
type customResponseWriter struct {
	http.ResponseWriter
	status int
	buf    *bytes.Buffer
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.status = code
	crw.ResponseWriter.WriteHeader(code)
}

func (crw *customResponseWriter) Write(b []byte) (int, error) {
	crw.buf.Write(b)
	return crw.ResponseWriter.Write(b)
}

func (s *Server) Run() error {
	slog.Info(fmt.Sprintf("Starting HTTP server at :%d ...", s.port))
	s.router.Use(otelmux.Middleware(shared.ServiceName))
	s.setupRouter()

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(s.router)

	srv := &http.Server{
		Handler:      loggingMiddleware(handler),
		Addr:         fmt.Sprintf(":%d", s.port),
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	return srv.ListenAndServe()
}
