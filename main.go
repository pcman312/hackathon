package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
	"github.com/vrecan/death"
	"google.golang.org/grpc"
	"gopkg.in/olivere/elastic.v3"

	"github.com/pcman312/hackathon/conf"
	"github.com/pcman312/hackathon/protos"
	"github.com/pcman312/hackathon/services"
)

func main() {
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")
	checkErr(err, "unable to load seelog config")

	err = log.UseLogger(logger)
	checkErr(err, "unable to use logger from config")
	defer log.Flush()

	opts, err := conf.LoadOpts("env")
	checkErr(err, "unable to load config")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.Port))
	checkErrf(err, "failed to listen on port [%d]", opts.Port)

	grpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(opts.GRPCConnTimeout),
	)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: opts.ESSkipVerify,
			},
			IdleConnTimeout: 5 * time.Minute,
			MaxIdleConns:    1000,
		},
		Timeout: opts.ESTimeout,
	}

	log.Debugf("Creating ES client to [%s]...", opts.ESHost)
	esClient, err := elastic.NewSimpleClient(
		elastic.SetHttpClient(httpClient),
		elastic.SetURL(opts.ESHost.String()),
		elastic.SetScheme(opts.ESHost.Scheme),
		elastic.SetErrorLog(ESErrorLogger{}),
		// elastic.SetInfoLog(ESInfoLogger{}),
		// elastic.SetTraceLog(ESTraceLogger{}),
		// elastic.SetSniff(b.sniff),
	)
	checkErr(err, "unable to create elasticsearch client")
	log.Debugf("Done creating ES client")

	svc, err := services.NewHackathonService(esClient, "hackathon")
	checkErr(err, "unable to create company setting service")

	hackathon.RegisterHackathonServiceServer(grpcServer, svc)

	log.Infof("Hackathon service listening on port %d...", opts.Port)

	d := death.NewDeath(syscall.SIGTERM, syscall.SIGINT)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintln(rw, "Hello, World!")
	})

	go func() {
		err = http.ListenAndServe(":80", mux)
		if err != nil {
			log.Criticalf("Failed to serve HTTP endpoint: %s", err)
		}
		d.FallOnSword()
	}()

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Criticalf("Failed to serve GRPC endpoint: %s", err)
		}
		d.FallOnSword()
	}()

	d.WaitForDeathWithFunc(func() {
		done := make(chan bool, 0)
		go func(done chan bool) {
			log.Info("Gracefully stopping hackathon service...")
			grpcServer.GracefulStop()
			log.Info("Goodbye")
			close(done)
		}(done)

		timer := time.NewTimer(opts.ShutdownWait)
		select {
		case <-timer.C:
			log.Info("Hackathon service didn't stop in %s - trying to force stop")
			grpcServer.Stop()
		case <-done:
			timer.Stop()
			break
		}
	})
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Critical(errors.Wrap(err, msg))
		panic(errors.Wrap(err, msg))
	}
}

func checkErrf(err error, format string, vals ...interface{}) {
	checkErr(err, fmt.Sprintf(format, vals...))
}
