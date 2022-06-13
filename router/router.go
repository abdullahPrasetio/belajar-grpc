package router

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"belajar-grpc/configs"
	"belajar-grpc/package/connections"

	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
)
  



func Setup() (*configs.Config,*logrus.Logger, error) {
	config,err:=configs.LoadConfig(".")
	if err != nil{
		return nil,nil,err
	}

	logger := logrus.New()
	// remove all writes to regular log to logger
	log.SetOutput(logger.Writer())
	log.SetFlags(0)

	// configure logging for environment
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		ForceQuote:    true,
		FullTimestamp: true,
	}

	logger.Info("[CONFIG] Setup complete")

	// Setup MySQL
	_,err= connections.GetConnection(config)
	if err != nil {
		fmt.Println(err)
		logger.Error("[CONFIG] Setup mysql error")
		return nil,nil,err
	}

	logger.Info("[CONFIG] Setup mysql complete")
	
	return &config,logger,err
}

// IgnoreErr returns true when err can be safely ignored.
func IgnoreErr(err error) bool {
	switch {
	case err == nil || err == http.ErrServerClosed || err == cmux.ErrListenerClosed:
		return true
	}
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err.Error() == "use of closed network connection"
	}
	return false
}