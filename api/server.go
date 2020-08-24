package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	httpClient *http.Client
	waitFor    time.Duration
	httpPort   int
	routesFile string
)

func init() {
	waitFor = 5 * time.Second
	envWaitFor, err := strconv.Atoi(os.Getenv("SHARD_ROUTER_RESTART_WAIT"))
	if err == nil && envWaitFor > 0 {
		waitFor = time.Duration(envWaitFor) * time.Second
	}

	httpPort = 8080
	envHTTPPort, err := strconv.Atoi(os.Getenv("SHARD_ROUTER_HTTP_PORT"))
	if err == nil && envHTTPPort > 0 {
		httpPort = envHTTPPort
	}

	routesFile = os.Getenv("SHARD_ROUTER_ROUTES_FILE")
	if routesFile == "" {
		routesFile = "routes.json" // try a default
	}

	httpTimeout := 10 * time.Second
	envHTTPTimeout, err := strconv.Atoi(os.Getenv("SHARD_ROUTER_HTTP_TIMEOUT"))
	if err == nil && envHTTPTimeout > 0 {
		httpTimeout = time.Duration(envHTTPTimeout) * time.Second
	}

	// setup http client
	tr := &http.Transport{DisableCompression: true}
	httpClient = &http.Client{Transport: tr, Timeout: httpTimeout}
}

func StartServer(c *cli.Context) (err error) {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGQUIT)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", httpPort))
	if err != nil {
		return
	}

	server := http.Server{Handler: routes()}

	serveResult := make(chan error)
	go func() {
		fmt.Fprintf(os.Stdout, "listening on :%d\n", httpPort)
		serveResult <- server.Serve(listener)
	}()

	// wait for either a signal to be received or for the server to
	// error out for some other reason
	signal := <-stop
	fmt.Fprintf(os.Stdout, "Received signal %s\n", signal)
	fmt.Fprintf(os.Stdout, "Waiting %s for requests to finish\n", waitFor)
	ctx, cancel := context.WithTimeout(context.Background(), waitFor)
	defer cancel()
	err = server.Shutdown(ctx)

	srerr := <-serveResult // always capture this
	close(serveResult)
	if err == nil {
		// only emit the result of .Serve if .Shutdown doesn't result in an error
		err = srerr
	}
	return err
}
