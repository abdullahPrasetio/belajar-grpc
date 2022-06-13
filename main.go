package main

import (
	"belajar-grpc/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	config, logger, err := router.Setup()
	if err != nil {
		log.Fatal(err)
	}
	
	logger.Infof("[SERVER] Environment %s is ready", config.APP_ENV)

	// create run group
	g,_:=errgroup.WithContext(context.Background())
	
	var servers []*http.Server
	// goroutine to check for signals to gracefully finish all functions
	g.Go(func() error {
		signalChannel:=make(chan os.Signal, 1)
		signal.Notify(signalChannel,os.Interrupt,syscall.SIGTERM)

		select {
		case sig:= <-signalChannel:
			logger.Infof("received signal: %s\n", sig)

			for i, s:= range servers{
				if err:=s.Shutdown(context.Background()); err!=nil {
					if err!= nil{
						logger.Infof("error shutting down server %d: %v",i,err)
						return err
					}
				}
					
			}
			os.Exit(1)
		}
		return nil
	})


	g.Go(func() error { return router.NewGRPCServer(config, logger) })
	g.Go(func() error { return router.NewHTTPServer(config, logger) })

	if err := g.Wait(); !router.IgnoreErr(err) {
		logger.Fatal(err)
	}
	logger.Infoln("done.")
}