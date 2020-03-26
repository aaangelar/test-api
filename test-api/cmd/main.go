package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	 test "test-api"

	"github.com/go-kit/kit/log"
)

func main() {

	fmt.Println("test-api")

	cnf := test.NewConfig()

	var (
		httpAddr = flag.String("http.addr", ":"+cnf.APIPort, "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	var ctx context.Context
	{
		ctx = context.Background()
	}

	var s test.Service
	{
		s = test.NewDataExportService()
		s = test.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = test.MakeHTTPHandler(ctx, s, log.NewContext(logger).With("component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
