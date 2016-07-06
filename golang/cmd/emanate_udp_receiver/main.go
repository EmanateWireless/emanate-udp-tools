package main

import (
	"fmt"
	"os"

	"github.com/EmanateWireless/emanate-udp-tools/golang/ccx"
	"github.com/EmanateWireless/emanate-udp-tools/golang/udp"
	"github.com/urfave/cli"
)

func main() {
	// add some console output white-space
	fmt.Println("")

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
		// create a udp receiver instance
		receiver := udp.NewReceiver(&udp.ReceiverOptions{
			Port: c.Int("port"),
		})

		// register the data handler
		receiver.DataHandler(func(du *udp.DataUpdate) {
			fmt.Printf("\nUDP PACKET RECEIVED\n")
			fmt.Printf("===================\n\n")
			fmt.Printf("  - Total Bytes = %d\n", len(du.Data))
			fmt.Printf("  - Remote Addr = %s:%d\n", du.RemoteIP, du.RemotePort)

			// parse the udp data as a ccx packet and dump to the console
			if err := ccx.Parse(du.Data); err != nil {
				fmt.Printf("Error receiving UDP data ('%v')\n", err)
			}
		})

		// start receiving packets
		receiver.Run()

		return nil
	}

	// start the cli app
	app.Run(os.Args)
}
