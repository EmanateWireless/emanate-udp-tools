package ccx

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Parse parses the given packet instance and dumps to the console
func Parse(data []byte) error {
	// parse the packet (everything but the variable length telemetry fields)
	packet := ParsedPacket{}
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, &packet); err != nil {
		return err
	}

	// dump the packet info
	fmt.Printf("  - Sequence = %d\n", packet.Sequence)
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

	return nil
}
