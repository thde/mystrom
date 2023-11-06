package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"thde.io/mystrom"
)

const (
	DiscoverCommand = "discover"
)

func run() error {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s commands:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), " discover - discover local mystrom devices\n")
		fmt.Fprintf(flag.CommandLine.Output(), " switch [address] (on|off|toggle) - control switch\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	switch flag.Arg(0) {
	case "discover":
		return discover(":7979")
	case "switch":
		return sw(flag.Args()[1:])
	default:
		flag.Usage()
		return nil
	}
}

func discover(address string) error {
	log.Printf("listening on %s", address)

	for {
		discover := mystrom.Discover{Address: address}
		device, err := discover.Device(context.Background())
		if err != nil {
			return fmt.Errorf("error discovering devices: %w", err)
		}

		log.Printf("%s: %+v", address, device)
	}
}

func sw(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("switch requires address and command arguments")
	}

	u, err := url.Parse("http://" + args[0])
	if err != nil {
		return fmt.Errorf("error parsing url for switch %s: %w", args[0], err)
	}

	sw := mystrom.NewSwitch(u)

	switch args[1] {
	case "on":
		return sw.On(context.Background())
	case "off":
		return sw.On(context.Background())
	case "toggle":
		return sw.On(context.Background())
	default:
		return fmt.Errorf("argument '%s' is not defined", args[0])
	}
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
