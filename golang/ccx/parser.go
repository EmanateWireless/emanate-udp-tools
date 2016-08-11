package ccx

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/EmanateWireless/emanate-udp-tools/golang/util"
)

// Parse parses the given packet instance and dumps to the console
func Parse(data []byte) error {
	// parse the packet (everything but the variable length telemetry fields)
	packet := ParsedPacket{}
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, &packet); err != nil {
		return err
	}

	// dump the static packet info
	dumpStaticInfo(&packet)

	// get the start of the telemetry data
	telemetryData := data[TelemetryDataOffset:]

	// dump the telemetry info
	err := dumpTelemetryInfo(telemetryData)

	// return whether an error occurred
	return err
}

func dumpStaticInfo(packet *ParsedPacket) {
	// dump the static packet info
	fmt.Printf("  - UDP Version = %d\n", packet.EmanateHeader.UDPVersion)
	fmt.Printf("  - Tag MAC = %s\n", util.MACBytesToString(packet.EmanateHeader.TagMACAddr))
	fmt.Printf("  - AP MAC = %s\n", util.MACBytesToString(packet.EmanateHeader.APMACAddr))
	fmt.Printf("  - Sequence = %d\n", packet.EmanateHeader.Sequence)
	fmt.Printf("  - Header\n")
	fmt.Printf("    - Protocol Version = %d\n", packet.Header.Version)
	fmt.Printf("    - Transmit Power = %d\n", packet.Header.Power)
	fmt.Printf("    - Wifi Channel = %d\n", packet.Header.Channel)
	fmt.Printf("    - Burst Length = %d\n", packet.Header.Burst)
	fmt.Printf("  - System Group\n")
	fmt.Printf("    - ID = %d\n", packet.System.ID)
	fmt.Printf("    - Length = %d\n", packet.System.Length)
	fmt.Printf("    - Product Type = %d\n", packet.System.ProductType)
	fmt.Printf("  - Battery Group\n")
	fmt.Printf("    - ID = %d\n", packet.Battery.ID)
	fmt.Printf("    - Length = %d\n", packet.Battery.Length)
	fmt.Printf("    - Tolerance = %d %%\n", (packet.Battery.Percent&0x07)*10)
	fmt.Printf("    - Charge = %d %%\n", ((packet.Battery.Percent>>3)&0x0F)*10)
	fmt.Printf("    - Days Remaining = %d\n", packet.Battery.Days)
	fmt.Printf("    - Age = %d days\n", packet.Battery.Age)
}

func dumpTelemetryInfo(data []byte) error {
	//fmt.Printf("TELEMETRY: data len = %d, tel data = %v\n", len(data), data)
	cursor := 0

	// iterate through all of the telemetry entries
	for cursor < len(data) {
		// if not enough data is available (groupID + groupLength + telemetryType)
		if len(data) < 3 {
			// log the error and return now
			msg := "Truncated of malformed Emanate CCX UDP packet"
			fmt.Printf("\nERROR: %s\n\n", msg)
			return fmt.Errorf(msg)
		}

		// get the group id
		groupID := data[cursor]
		cursor = cursor + 1

		// determine the group id
		switch groupID {
		// telemetry group id (3)
		case TelemetryGroupID:
			// get the group length
			groupLength := int(data[cursor])
			cursor = cursor + 1

			// get the telemetry type
			telemetryType := data[cursor]
			cursor = cursor + 1

			// determine the telemetry type
			switch telemetryType {
			case TemperatureTelemetryType:
				// if not enough data is available (tempC)
				if len(data)-cursor < 4 {
					// log the error and return now
					msg := "Truncated of malformed temperature group in Emanate CCX UDP packet"
					fmt.Printf("\nERROR: %s\n\n", msg)
					return fmt.Errorf(msg)
				}

				tempC := util.Float32FromBytes(data[cursor:])
				cursor = cursor + 4

				fmt.Printf("  - Temperature Group\n")
				fmt.Printf("    - Group ID = %d\n", groupID)
				fmt.Printf("    - Group Length = %d\n", groupLength)
				fmt.Printf("    - Type = %d\n", telemetryType)
				fmt.Printf("    - Temperature = %.2f C\n", tempC)

			case StatusTelemetryType:
				// if not enough data is available (tempC)
				if len(data)-cursor < 1 {
					// log the error and return now
					msg := "Truncated of malformed status telemetry group in Emanate CCX UDP packet"
					fmt.Printf("\nERROR: %s\n\n", msg)
					return fmt.Errorf(msg)
				}

				statusLength := int(data[cursor])
				cursor = cursor + 1

				// if not enough data is available (tempC)
				if len(data)-cursor < statusLength {
					// log the error and return now
					msg := "Truncated of malformed status telemetry string in Emanate CCX UDP packet"
					fmt.Printf("\nERROR: %s\n\n", msg)
					return fmt.Errorf(msg)
				}

				// extract the utf-16 status string and advance the cursor
				statusUTF16 := string(data[cursor : cursor+statusLength])
				cursor = cursor + statusLength

				// convert the utf-16 string into an ascii string
				statusASCII := util.UTF16ToASCII(statusUTF16)

				fmt.Printf("  - Status Group\n")
				fmt.Printf("    - Group ID = %d\n", groupID)
				fmt.Printf("    - Group Length = %d\n", groupLength)
				fmt.Printf("    - Type = %d\n", telemetryType)
				fmt.Printf("    - Status Length = %d\n", statusLength)
				fmt.Printf("    - Status = '%s'\n", statusASCII)

			default:
				fmt.Printf("  TELEMETRY TYPE ERROR! (%d)", telemetryType)
			}

		default:
			fmt.Printf("  TELEMETRY GROUP ID ERROR! (%d)", groupID)
		}
	}

	return nil
}
