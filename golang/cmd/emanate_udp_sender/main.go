package main

import (
	"os"

	"github.com/EmanateWireless/emanate-udp-tools/udp"
	"github.com/urfave/cli"
)

func main() {
	// create the cli app
	app := cli.NewApp()
	app.Version = "v1.0.0"
	app.Name = "emanate_udp_sender"
	app.Usage = "Emanate PowerPath UDP CCX packet transmitter"

	// define the cli flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "127.0.0.1",
			Usage: "udp target hostname or ip-address",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 9999,
			Usage: "udp target port number",
		},
	}

	// define the cli execution handler
	app.Action = func(c *cli.Context) error {
		// create a udp sender instance
		sender := udp.NewSender(&udp.SenderOptions{
			Host: c.String("host"),
			Port: c.Int("port"),
		})

		// send the packets
		sender.Transmit([]byte("Howdy there!"))

		return nil
	}

	// start the cli app
	app.Run(os.Args)
}
