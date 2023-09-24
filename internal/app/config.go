package app

import (
	"context"
	"github.com/go-yaml/yaml"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
	"practice/internal/pkg/users"
	"strings"
)

type Config struct {
	users.User `yaml:",inline"`
	Port       int          `yaml:"port"`
	Host       string       `yaml:"host"`
	Root       string       `yaml:"root"`
	Users      []users.User `yaml:"users"`
	Prefix     string       `yaml:"prefix"`

	middleware.CORSConfig `yaml:"cors"`
}

func NewConfig(path string) (Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to read config file"))
	}

	var cfg Config
	if err = yaml.NewDecoder(strings.NewReader(string(file))).Decode(&cfg); err != nil {
		return Config{}, errors.Wrap(err, "failed to decode config")
	}

	cfg.Handler = &webdav.Handler{
		Prefix:     cfg.Prefix,
		FileSystem: webdav.NewMemFS(),
		LockSystem: webdav.NewMemLS(),
		Logger:     NewLogger(logrus.New()),
	}

	return cfg, nil
}

func (c Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := c.User

	// Checks for user permissions relatively to this PATH.
	noModification := r.Method == "GET" || r.Method == "HEAD" ||
		r.Method == "OPTIONS" || r.Method == "PROPFIND"

	if !u.ACL(r.URL.Path, noModification) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method == "HEAD" {
		w = newResponseWriterNoBody(w)
	}

	if r.Method == "GET" && strings.HasPrefix(r.URL.Path, u.Handler.Prefix) {
		info, err := u.Handler.FileSystem.Stat(context.TODO(), strings.TrimPrefix(r.URL.Path, u.Handler.Prefix))
		if err == nil && info.IsDir() {
			r.Method = "PROPFIND"

			if r.Header.Get("Depth") == "" {
				r.Header.Add("Depth", "1")
			}
		}
	}

	u.Handler.ServeHTTP(w, r)
}

func NewLogger(log *logrus.Logger) func(*http.Request, error) {
	return func(r *http.Request, err error) {
		log.WithFields(logrus.Fields{
			"uri":     r.RequestURI,
			"host":    r.Host,
			"headers": r.Header,
		}).Error(err)
	}
}

type responseWriterNoBody struct {
	http.ResponseWriter
}

func newResponseWriterNoBody(w http.ResponseWriter) *responseWriterNoBody {
	return &responseWriterNoBody{w}
}

func (w responseWriterNoBody) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w responseWriterNoBody) Write(_ []byte) (int, error) {
	return 0, nil
}

func (w responseWriterNoBody) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
