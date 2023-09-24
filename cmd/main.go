package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"log"
	"practice/internal/app"
	"practice/internal/pkg/auth"
	"strconv"
)

func main() {
	cfgPath := flag.String("config", "configs/app.yaml", "path to app config")
	flag.Parse()

	if cfgPath == nil {
		log.Fatal("no config path")
	}

	cfg, err := app.NewConfig(*cfgPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to make new config"))
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(cfg.CORSConfig))
	e.Use(middleware.BasicAuth(auth.NewBasicAuth(cfg.Users)))

	h := app.NewHandler(cfg)
	e.Any("*", h.HandlerWebDAV)

	if err = e.Start(cfg.Host + ":" + strconv.Itoa(cfg.Port)); err != nil {
		log.Fatal(errors.Wrap(err, "server stopped with error"))
	}

	log.Println("server stopped successfully")
}
