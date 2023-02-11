package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/KarolosLykos/ptask/docs"
	"github.com/KarolosLykos/ptask/internal/api"
	"github.com/KarolosLykos/ptask/internal/constants"
	"github.com/KarolosLykos/ptask/internal/logger/log"
	"github.com/KarolosLykos/ptask/internal/ptask/usecase"
)

var (
	host, port string
	debug      bool
)

//	@title			Periodic Task Api
//	@version		1.0
//	@description	JSON/HTTP service in Golang, that returns the matching timestamps of a periodic task.

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	ctx := context.Background()

	flag.StringVar(&host, "host", "0.0.0.0", "-host localhost")
	flag.StringVar(&port, "port", "8080", "-port 8080")
	flag.BoolVar(&debug, "debug", false, "-debug")
	flag.Parse()

	docs.SwaggerInfo.Host = host + ":" + port

	// init logger.
	logger := log.Default(debug, constants.LoggerFormat)

	// init periodic task useCase.
	useCase := usecase.NewPeriodicTaskUC(logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	// init server.
	s := api.New(logger, host+":"+port, useCase)

	// start server.
	s.Start(ctx)

	event := <-quit
	logger.Info(ctx, fmt.Sprintf("received signal: %v", event))

	// shutdown server.
	s.Shutdown(ctx)
}
