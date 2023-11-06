# mystrom

The mystrom package provides Go client for the [myStrom REST API](https://api.mystrom.ch/).

## Features

- [x] Switch
- [ ] Button
- [ ] Bulb
- [ ] LED Strip
- [ ] PIR
- [ ] New Button Plus
- [x] Discovery

PR's for additional endpoints are welcome!

## Usage

To install the mystrom package, use the following command:

```go
package main

import (
	"context"
	"log"
	"net"
	"net/url"

	"thde.io/mystrom"
)

func main() {
	client := mystrom.NewClient()

	for {
		discover := mystrom.Discover{}
		device, err := discover.Device(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%+v", device)

		host, _, err := net.SplitHostPort(device.Address.String())
		if err != nil {
			log.Fatal(err)
		}
		if device.Type != mystrom.DeviceTypeSwitchCH && device.Type != mystrom.DeviceTypeSwitchEU {
			continue
		}

		u, err := url.Parse("http://" + host)
		if err != nil {
			log.Fatal(err)
		}
		sw := client.NewSwitch(u)
		temp, err := sw.Temperature(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", temp)
	}
}

```

## CLI

To install the CLI, run:

```shell
go install thde.io/mystrom/cmd/mystrom@latest
```
