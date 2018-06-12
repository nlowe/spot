package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"bitbucket.hylandqa.net/do/spot/pkg/spot"

	arg "github.com/alexflint/go-arg"
	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

type applicationArgs struct {
	Jenkins   []string `arg:"-j,separate" help:"Jenkins Url & credentials in the form of https://jenkins/,username,password"`
	Slack     string   `arg:"-s" help:"Slack-Compatible Incoming Webhook URL"`
	Verbosity string   `arg:"-v" help:"Verbosity [panic, fatal, error, warn, info, debug]"`
	Period    string   `arg:"-p" help:"How long to wait between checks"`
	Once      bool     `arg:"-o" help:"Run checks once and exit"`
}

func (applicationArgs) Description() string {
	return "alerts for disconnected build agents"
}

type dummyNotifier struct{}

func (d *dummyNotifier) Notify(agents map[string][]string) error { return nil }

func initLogrus(level string) {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())

	if level, err := log.ParseLevel(strings.ToLower(level)); err != nil {
		log.Panic(err)
	} else {
		log.SetLevel(level)
	}
}

func watchAllTheThings(interval time.Duration, w *spot.Watchdog, shutdown chan bool) {
	for {
		// do the thing
		if err := w.RunChecks(); err != nil {
			log.WithError(err).Error("Watchdog checks failed")
		}

		// wait for next interval
		select {
		case <-shutdown:
			return
		case <-time.After(interval):
			break
		}
	}
}

func main() {
	args := &applicationArgs{}
	args.Verbosity = "info"
	p := arg.MustParse(args)

	initLogrus(args.Verbosity)
	log.Info("Hello, World!")

	watchdog := &spot.Watchdog{
		Detectors:           []spot.OfflineAgentDetector{},
		NotificationHandler: &dummyNotifier{},
	}

	if args.Slack != "" {
		var err error
		if watchdog.NotificationHandler, err = spot.NewSlackNotifier(args.Slack); err != nil {
			p.Fail(fmt.Sprintf("Invalid slack URL: %s", err.Error()))
		}
	}

	for _, v := range args.Jenkins {
		l := log.WithField("jenkins", v)
		l.Debug("Trying to parse jenkins instance")

		if detector, err := spot.NewJenkinsDetectorFromArg(v); err != nil {
			p.Fail(fmt.Sprintf("Failed to parse jenkins configuration: %s", err.Error()))
		} else {
			watchdog.Detectors = append(watchdog.Detectors, detector)
		}
	}

	if len(watchdog.Detectors) == 0 {
		p.Fail("Provide at least one watchdog configuration")
	}

	if args.Once {
		if err := watchdog.RunChecks(); err != nil {
			panic(err)
		}
	} else {
		period, err := time.ParseDuration(args.Period)
		if err != nil {
			p.Fail(fmt.Sprintf("Failed to parse period: %s", err.Error()))
		}

		shutdown := make(chan bool)
		go watchAllTheThings(period, watchdog, shutdown)

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		shutdown <- true
	}

	log.Info("Goodbye")
}
