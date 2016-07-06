package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/EmanateWireless/emanate-udp-tools/golang/ccx"
	"github.com/EmanateWireless/emanate-udp-tools/golang/udp"
	"github.com/urfave/cli"
)

const (
	MinSeqNumber = 0
	MaxSeqNumber = 65535
)

func main() {
	// add some console output white-space
	fmt.Println("")

	// create the cli app
	app := cli.NewApp()
	app.Version = "v1.0.0"
	app.Name = "emanate_udp_sender"
	app.HelpName = "emanate_udp_sender"
	app.Usage = "Emanate PowerPath UDP CCX packet transmitter"
	app.UsageText = "emanate_udp_sender --host <IP> --port <PORT> [options]"

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
		cli.BoolTFlag{
			Name:  "all",
			Usage: "sends all possible udp message options for testing",
		},
		cli.IntFlag{
			Name:  "num-dups",
			Value: 0,
			Usage: "number of duplicate udp packets to send",
		},
		cli.IntFlag{
			Name:  "dup-interval-ms",
			Value: 100,
			Usage: "delay interval between duplicate udp packets",
		},
		cli.IntFlag{
			Name:  "seq",
			Value: 1,
			Usage: "sequence number of the emanate udp packet",
		},
		cli.StringFlag{
			Name:  "util-state",
			Value: "unplugged",
			Usage: "utility state of 'unplugged', 'off', 'idle', or 'active'",
		},
		cli.Float64Flag{
			Name:  "temp",
			Value: 12.34,
			Usage: "temperature floating-point value (in celsius)",
		},
		cli.IntFlag{
			Name:  "battery-charge",
			Value: 80,
			Usage: "battery charge percentage remaining (0-100)",
		},
		cli.IntFlag{
			Name:  "battery-days-remaining",
			Value: 100,
			Usage: "number of days remaining for battery charge",
		},
		cli.IntFlag{
			Name:  "battery-age",
			Value: 10,
			Usage: "battery age in days",
		},
		cli.IntFlag{
			Name:  "battery-tolerance",
			Value: 0,
			Usage: "battery prediction tolerance percentage (0-100)",
		},
		cli.BoolTFlag{
			Name:  "button-pressed",
			Usage: "adds the button-pressed telemetry status",
		},
		cli.IntFlag{
			Name:  "door-open-percent",
			Value: 22,
			Usage: "percentage of time the fridge door has been open",
		},
		cli.IntFlag{
			Name:  "high-power-percent",
			Value: 33,
			Usage: "percentage of time the device ran in high-power-mode",
		},
		cli.BoolTFlag{
			Name:  "probe-unplugged",
			Usage: "adds the 'temp probe unplugged' alert telemetry status",
		},
		cli.BoolTFlag{
			Name:  "probe-invalid-value",
			Usage: "adds the 'temp probe invalid value' alert telemetry status",
		},
		cli.IntFlag{
			Name:  "product-type",
			Value: 0,
			Usage: "product-type code set in the ccx 'system group'",
		},
	}

	// define the cli execution handler
	app.Action = func(c *cli.Context) error {
		// create a udp sender instance
		sender := udp.NewSender(&udp.SenderOptions{
			Host: c.String("host"),
			Port: c.Int("port"),
		})

		// check if the kitchen-sink 'all' option is enabled
		sendAll := false
		if c.GlobalIsSet("all") {
			// enable the kitchen-sink 'all' option
			sendAll = true
		}

		// create a ccx packet
		packet := ccx.NewPacket()

		// set the burst length to the number of duplicate udp packets + 1 original
		packet.SetBurstLength(uint8(c.Int("num-dups") + 1))

		// if the sequence number is given
		if sendAll || c.GlobalIsSet("seq") {
			// get the sequence number flag value
			seq := c.Int("seq")

			// validate the given sequence number
			if (seq < MinSeqNumber) || (seq > MaxSeqNumber) {
				msg := fmt.Sprintf("sequence number must be between '%d' and '%d' (inclusive)",
					MinSeqNumber, MaxSeqNumber)
				exitNow(msg)
			}

			// set the packet's sequence number
			packet.SetSequenceNumber(uint16(seq))
		}

		// if the util-state option is given
		if sendAll || c.GlobalIsSet("util-state") {
			// get the utility-state flag value
			util := strings.ToLower(c.String("util-state"))

			// add the util-state telemetry value to the packet
			var err error
			switch util {
			case "unplugged":
				err = packet.SetUtilState(ccx.UtilStateUnplugged)
			case "off":
				err = packet.SetUtilState(ccx.UtilStatePluggedInOff)
			case "idle":
				err = packet.SetUtilState(ccx.UtilStatePluggedInIdle)
			case "active":
				err = packet.SetUtilState(ccx.UtilStatePluggedInActive)
			default:
				// display the cli usage and exit with an error
				exitNow("'util-state' value must be either 'unplugged', 'off', idle', or 'active'")
			}

			// if an error occurred while addeding the util-state telemetry to the packet
			if err != nil {
				msg := fmt.Sprintf("cannot add util-state '%s' to UDP packet", util)
				exitNowWithError(msg, err)
			}
		}

		// set the battery values
		packet.SetBatteryInfo(&ccx.BatteryInfo{
			TolerancePercent: uint8(c.Int("battery-tolerance")),
			PercentRemaining: uint8(c.Int("battery-charge")),
			DaysRemaining:    uint16(c.Int("battery-days-remaining")),
			AgeDays:          uint32(c.Int("battery-age")),
		})

		// if the temperature option is given
		if sendAll || c.GlobalIsSet("temp") {
			// add the temperature telemetry value to the packet
			if err := packet.SetTemperature(float32(c.Float64("temp"))); err != nil {
				exitNowWithError("cannot add temperature to UDP packet", err)
			}
		}

		// if the door-open-percent option is given
		if sendAll || c.GlobalIsSet("door-open-percent") {
			// get the door-open-percent flag value
			percent := c.Int("door-open-percent")

			// validate the value
			if (percent < 0) || (percent > 100) {
				exitNow("door-open-percent must be between 0 and 100")
			}

			// add the door-open telemetry status to the packet
			if err := packet.SetDoorOpenPercent(percent); err != nil {
				exitNowWithError("cannot add 'door-open-percent' to UDP packet", err)
			}
		}

		// if the high-power-percent option is given
		if sendAll || c.GlobalIsSet("high-power-percent") {
			// get the high-power-percent flag value
			percent := c.Int("high-power-percent")

			// validate the value
			if (percent < 0) || (percent > 100) {
				exitNow("high-power-percent must be between 0 and 100")
			}

			// add the high-power telemetry status to the packet
			if err := packet.SetHighPowerPercent(percent); err != nil {
				exitNowWithError("cannot add 'high-power-percent' to UDP packet", err)
			}
		}

		// if the product-type option is given
		if c.GlobalIsSet("product-type") {
			// get the product-type flag value
			pt := uint16(c.Int("product-type"))

			// set the product-type value
			packet.SetProductType(pt)
		}

		// if the button-pressed option is given
		if sendAll || c.GlobalIsSet("button-pressed") {
			// add the button-pressed telemetry status to the packet
			if err := packet.SetButtonPressed(); err != nil {
				exitNowWithError("cannot add 'button-pressed' to UDP packet", err)
			}
		}

		// if the probe-unplugged option is given
		if sendAll || c.GlobalIsSet("probe-unplugged") {
			// add the probe-unplugged telemetry status to the packet
			if err := packet.SetProbeUnplugged(); err != nil {
				exitNowWithError("cannot add 'probe-unplugged' to UDP packet", err)
			}
		}

		// if the probe-invalid-value option is given
		if sendAll || c.GlobalIsSet("probe-invalid-value") {
			// add the probe-invalid-value telemetry status to the packet
			if err := packet.SetProbeInvalidValue(); err != nil {
				exitNowWithError("cannot add 'probe-invalid-value' to UDP packet", err)
			}
		}

		// encode the packet as binary bytes
		data, err := packet.Pack()
		if err != nil {
			exitNowWithError("cannot convert UDP packet into bytes", err)
		}

		// send the first udp ccx packet
		sender.Transmit(data)

		// if the option to send duplicate packets is given
		if c.GlobalIsSet("num-dups") {
			// get the number of duplicate packets to transmit
			ndups := c.Int("num-dups")

			// get the delay interval time between each duplicate packet
			dupDelayMs := c.Int("dup-interval-ms")

			// transmit each duplicate packet and delay the configurable time
			for i := 0; i < ndups; i++ {
				// wait for the configured interval time
				time.Sleep(time.Duration(dupDelayMs) * time.Millisecond)

				// send the next duplicate udp ccx packet
				sender.Transmit(data)
			}
		}

		// log that we are done
		log.Printf("DONE!")
		fmt.Println("")

		// return successfully
		return nil
	}

	// start the cli app
	app.Run(os.Args)
}

func exitNow(msg string) {
	exitNowWithError(msg, nil)
}

func exitNowWithError(msg string, err error) {
	// log the error message and exit with an error
	if err != nil {
		log.Printf("ERROR: %s ('%v')", msg, err)
	} else {
		log.Printf("ERROR: %s", msg)
	}
	os.Exit(1)
}
