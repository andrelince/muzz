package main

import (
	"context"
	"net/http"
	"time"

	"github.com/muzz/api/config"
	"github.com/muzz/api/di"
	"github.com/muzz/api/rest"

	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
)

func main() {
	c, err := di.NewDI()
	if err != nil {
		panic(err)
	}

	if err := c.Invoke(rest.NewRest); err != nil {
		panic(err)
	}

	if err := c.Invoke(start); err != nil {
		panic(err)
	}
}

func start(c config.Config, router *http.ServeMux) error {
	g, _ := errgroup.WithContext(context.Background())

	corss := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler:      corss.Handler(router),
		Addr:         ":" + c.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	g.Go(func() error { return srv.ListenAndServe() })

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
