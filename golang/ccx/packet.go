package ccx

import (
	"bytes"
	"encoding/binary"
)

// packet constants
const (
	EMANATE_PROTOCOL_VERSION      = 0
	DEFAULT_TX_POWER              = 17
	DEFAULT_WIFI_CHANNEL          = 1
	DEFAULT_REGULATORY_CLASS      = 0
	DEFAULT_BURST_LENGTH          = 3
	SYSTEM_GROUP_ID               = 0
	SYSTEM_GROUP_LEN              = 2
	BATTERY_GROUP_ID              = 2
	BATTERY_GROUP_LEN             = 6
	TELEMETRY_GROUP_ID            = 3
	CISCO_PRODUCT_TYPE            = 0
	UTIL_STATE_UNPLUGGED          = "UTIL_STATE=UNPLUGGED"
	UTIL_STATE_PLUGGED_IN_OFF     = "UTIL_STATE=PLUGGED_IN_OFF"
	UTIL_STATE_PLUGGED_IN_IDLE    = "UTIL_STATE=PLUGGED_IN_IDLE"
	UTIL_STATE_PLUGGED_IN_ACTIVE  = "UTIL_STATE=PLUGGED_IN_ACTIVE"
	UTIL_STATE_PLUGGED_IN_UNKNOWN = "UTIL_STATE=PLUGGED_IN_UNKNOWN"
)

// Packet instance struct
type Packet struct {
	Header    Header
	System    SystemGroup
	Battery   BatteryGroup
	Telemetry TelemetryGroup
}

// Header defines Emanate's encapsulating CCX packet header
type Header struct {
	Version         int8
	Power           int8
	Channel         int8
	RegulatoryClass int8
	Burst           int8
}

// SystemGroup defines the CCX system group structure
type SystemGroup struct {
	ID          int8
	Length      int8
	ProductType int8
}

// BatteryGroup defines the CCX battery group structure
type BatteryGroup struct {
	ID      int8
	Length  int8
	Percent int8
	Days    int8
	Age     int8
}

// TelemetryGroup defines a variable set of CCX telemetry types
type TelemetryGroup struct {
	ID     int8
	Length int8
	Data   bytes.Buffer
}

// StatusTelemetry defines the CCX unicode status string telemetry type
type StatusTelemetry struct {
	Type   int8
	Status string
}

// TemperatureTelemetry defines the CCX temperature telemetry type
type TemperatureTelemetry struct {
	Type    int8
	Celsius float32
}

// NewPacket creates a new CCX packet instance
func NewPacket() *Packet {
	// return a new default packet instance
	return &Packet{
		Header: Header{
			Version:         EMANATE_PROTOCOL_VERSION,
			Power:           DEFAULT_TX_POWER,
			Channel:         DEFAULT_WIFI_CHANNEL,
			RegulatoryClass: DEFAULT_REGULATORY_CLASS,
			Burst:           DEFAULT_BURST_LENGTH,
		},
		System: SystemGroup{
			ID:          SYSTEM_GROUP_ID,
			Length:      SYSTEM_GROUP_LEN,
			ProductType: CISCO_PRODUCT_TYPE,
		},
		Battery: BatteryGroup{
			ID:      BATTERY_GROUP_ID,
			Length:  BATTERY_GROUP_LEN,
			Percent: 100,
			Days:    100,
			Age:     1,
		},
		Telemetry: TelemetryGroup{
			ID:     TELEMETRY_GROUP_ID,
			Length: 0,
			Data:   bytes.Buffer{},
		},
	}
}

// SetProtocolVersion sets the protocol version field in the packet
func (p *Packet) SetProtocolVersion(version int8) {
	p.Header.Version = version
}

// SetTransmitPower sets the transmit power field in the packet
func (p *Packet) SetTransmitPower(power int8) {
	p.Header.Power = power
}

// SetRegulatoryClass sets the regulatory class field in the packet
func (p *Packet) SetRegulatoryClass(c int8) {
	p.Header.RegulatoryClass = c
}

// SetBurstLength sets the udp burst length field in the packet
func (p *Packet) SetBurstLength(b int8) {
	p.Header.Burst = b
}

// SetUtilState adds the utility state status telemetry string
func (p *Packet) SetUtilState(state string) error {
	// add the telemetry string to the data
	buf := p.Telemetry.Data
	_, err := buf.Write([]byte(state))
	return err
}

// Bytes returns the packet data structure as a slice of bytes
func (p *Packet) Bytes() ([]byte, error) {
	// encode the packet
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, binary.BigEndian, p.Header); err != nil {
		return []byte{}, err
	}
	if err := binary.Write(buf, binary.BigEndian, p.System); err != nil {
		return []byte{}, err
	}
	if err := binary.Write(buf, binary.BigEndian, p.Battery); err != nil {
		return []byte{}, err
	}

	// encode the telemetry group data
	tData := p.Telemetry.Data.Bytes()
	buf.WriteByte(uint8(TELEMETRY_GROUP_ID))
	buf.WriteByte(uint8(len(tData)))
	if err := binary.Write(buf, binary.BigEndian, tData); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
