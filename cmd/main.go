package main

import (
	"k8s.io/klog/v2"

	"github.com/crazytaxii/fake-dcgm-exporter/cmd/server"
)

func main() {
	klog.InitFlags(nil)

	cmd := server.NewCommand()
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
