package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/iomatters/config"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

var (
	BuildVersion = "(build)"
)

func init() {
	_ = flag.Set("logtostderr", "true")
}

func main() {
	filename := pflag.String("config", "", "path to config file")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	defer glog.Flush()

	params, err := config.NewAppConfig(*filename)
	if err != nil {
		glog.Fatalln(err)
	}

	glog.Infoln("Starting collector version", BuildVersion)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	collector := newCollector("cryptocompare", params.DB.Host, params.DB.Port, params.DB.DBName, params.DB.User, params.DB.Password, params.Main.Fsyms, params.Main.Tsyms, c)
	collector.run()

}
