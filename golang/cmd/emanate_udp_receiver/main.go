package main

import (
	"os"

	"github.com/EmanateWireless/emanate-udp-tools/golang/udp"
	"github.com/urfave/cli"
)

func main() {
	// create the cli app
	app := cli.NewApp()
	app.Version = "v1.0.0"
	app.Name = "emanate_udp_receiver"
	app.HelpName = "emanate_udp_receiver"
	app.Usage = "Emanate PowerPath UDP CCX packet receiver"
	app.UsageText = "emanate_udp_receiver --port <LISTENING-PORT>"

	// define the cli flags
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 9999,
			Usage: "local udp receiver port number",
		},
	}

	// define the cli execution handler
	app.Action = func(c *cli.Context) error {
		// create a udp reciver instance
		receiver := udp.NewReceiver(&udp.ReceiverOptions{
			Port: c.Int("port"),
		})

		// start receiving packets
		receiver.Run()

		return nil
	}

	// start the cli app
	app.Run(os.Args)
}
