package main

import (
	"fmt"
	"os"

	"github.com/EmanateWireless/emanate-udp-tools/golang/ccx"
	"github.com/EmanateWireless/emanate-udp-tools/golang/udp"
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

		// create a ccx packet
		packet := ccx.NewPacket()
		if err := packet.SetUtilState(ccx.UTIL_STATE_UNPLUGGED); err != nil {
			fmt.Printf("Error: cannot add util state to UDP packet ('%v')\n", err)
			return nil
		}

		// encode the packet as bytes
		data, err := packet.Bytes()
		if err != nil {
			fmt.Printf("Error: cannot convert UDP packet into bytes ('%v')\n", err)
			return nil
		}

		// send the ccx packet
		sender.Transmit(data)

		return nil
	}

	// start the cli app
	app.Run(os.Args)
}
