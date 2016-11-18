package main

import (
	"log"

	"github.com/spf13/pflag"
	"k8s.io/kubernetes/cmd/kubectl/app"
	"k8s.io/kubernetes/cmd/kubelet/app/options"
	"k8s.io/kubernetes/pkg/util/flag"
	"k8s.io/kubernetes/pkg/util/logs"
	"k8s.io/kubernetes/pkg/version/verflag"
)

func main() {
	log.Println("starting nspawnlet...")
	srv := options.NewKubeletServer()
	srv.AddFlags(pflag.CommandLine)
	flag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()
	verflag.PrintAndExitIfRequested()
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
