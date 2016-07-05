package ccx

import (
	"bytes"
	"encoding/binary"
)

// packet constants
const (
	EmanateProtocolVersion      = 0
	DefaultTxPower              = 17
	DefaultWifiChannel          = 1
	DefaultRegulatoryClass      = 0
	DefaultBurstLength          = 3
	SystemGroupID               = 0
	SystemGroupLength           = 2
	BatteryGroupID              = 2
	BatteryGroupLength          = 7
	DefaultBatteryTolerance     = 10
	DefaultBatteryCharge        = 80
	DefaultBatteryDaysRemaining = 100
	DefaultBatteryAgeDays       = 10
	TelemetryGroupID            = 3
	TemperatureTelemetryType    = 1
	TemperatureTelemetryLength  = 5
	StatusTelemetryType         = 8
	CiscoProductType            = 0
	UtilStateUnplugged          = "UTIL_STATE=UNPLUGGED"
	UtilStatePluggedInOff       = "UTIL_STATE=PLUGGED_IN_OFF"
	UtilStatePluggedInIdle      = "UTIL_STATE=PLUGGED_IN_IDLE"
	UtilStatePluggedInActive    = "UTIL_STATE=PLUGGED_IN_ACTIVE"
	ButtonPressedTelemetry      = "BUTTON=PRESSED"
)

// Packet instance struct
type Packet struct {
	Sequence      uint16
	Header        Header
	System        SystemGroup
	Battery       BatteryGroup
	TelemetryData bytes.Buffer
}

// ParsedPacket includes the fields that can be parsed by the encoding/binary package
type ParsedPacket struct {
	Sequence uint16
	Header   Header
	System   SystemGroup
	Battery  BatteryGroup
}

// Header defines Emanate's encapsulating CCX packet header
type Header struct {
	Version         uint8
	Power           uint8
	Channel         uint8
	RegulatoryClass uint8
	Burst           uint8
}

// SystemGroup defines the CCX system group structure
type SystemGroup struct {
	ID          uint8
	Length      uint8
	ProductType uint8
}

// BatteryGroup defines the CCX battery group structure
type BatteryGroup struct {
	ID      uint8
	Length  uint8
	Percent uint8
	Days    uint16
	Age     uint32
}

// BatteryInfo defines the info struct when setting the battery group info
type BatteryInfo struct {
	TolerancePercent uint8
	PercentRemaining uint8
	DaysRemaining    uint16
	AgeDays          uint32
}

// StatusTelemetry defines the CCX unicode status string telemetry type
type StatusTelemetry struct {
	ID     uint8
	Length uint8
	Type   uint8
	Status string
}

// TemperatureTelemetry defines the CCX temperature telemetry type
type TemperatureTelemetry struct {
	ID      uint8
	Length  uint8
	Type    uint8
	Celsius float32
}

// NewPacket creates a new CCX packet instance
func NewPacket() *Packet {
	// return a new default packet instance
	return &Packet{
		Sequence: uint16(1),
		Header: Header{
			Version:         EmanateProtocolVersion,
			Power:           DefaultTxPower,
			Channel:         DefaultWifiChannel,
			RegulatoryClass: DefaultRegulatoryClass,
			Burst:           DefaultBurstLength,
		},
		System: SystemGroup{
			ID:          SystemGroupID,
			Length:      SystemGroupLength,
			ProductType: CiscoProductType,
		},
		Battery: BatteryGroup{
			ID:      BatteryGroupID,
			Length:  BatteryGroupLength,
			Percent: ((DefaultBatteryCharge / 10) << 3) | (DefaultBatteryTolerance & 0x07),
			Days:    DefaultBatteryDaysRemaining,
			Age:     DefaultBatteryAgeDays,
		},
		TelemetryData: bytes.Buffer{},
	}
}

// NewStatusTelemetry creates a new status telemetry string instance
func NewStatusTelemetry(status string) StatusTelemetry {
	// encode the given string as UTF-16
	encodedBytes := []byte{}
	for _, b := range []byte(status) {
		encodedBytes = append(encodedBytes, 0x00, b)
	}
	encodedStr := string(encodedBytes)

	// create and return the status telemetry instance
	return StatusTelemetry{
		ID:     TelemetryGroupID,           // 3
		Length: uint8(len(encodedStr) + 1), // variable
		Type:   StatusTelemetryType,        // 8
		Status: encodedStr,
	}
}

// NewTemperatureTelemetry creates a new temperature telemetry instance
func NewTemperatureTelemetry(v float32) TemperatureTelemetry {
	return TemperatureTelemetry{
		ID:      TelemetryGroupID,           // 3
		Length:  TemperatureTelemetryLength, // 5
		Type:    TemperatureTelemetryType,   // 1
		Celsius: v,
	}
}

// SetSequenceNumber sets the packet sequence number
func (p *Packet) SetSequenceNumber(seq uint16) {
	p.Sequence = seq
}

// IncSequenceNumber increments the packet sequence number
func (p *Packet) IncSequenceNumber(seq uint16) {
	// increment the packet sequence number by 1
	p.IncSequenceNumberByValue(1)
}

// IncSequenceNumberByValue increments the packet sequence number by the given value
func (p *Packet) IncSequenceNumberByValue(v uint16) {
	p.Sequence = p.Sequence + v
}

// SetProtocolVersion sets the protocol version field in the packet
func (p *Packet) SetProtocolVersion(version uint8) {
	p.Header.Version = version
}

// SetTransmitPower sets the transmit power field in the packet
func (p *Packet) SetTransmitPower(power uint8) {
	p.Header.Power = power
}

// SetRegulatoryClass sets the regulatory class field in the packet
func (p *Packet) SetRegulatoryClass(c uint8) {
	p.Header.RegulatoryClass = c
}

// SetBurstLength sets the udp burst length field in the packet
func (p *Packet) SetBurstLength(b uint8) {
	p.Header.Burst = b
}

// SetBatteryInfo set the battery group with the given values
func (p *Packet) SetBatteryInfo(b *BatteryInfo) {
	// calculate the bit fields
	tolerancePercent := uint8(b.TolerancePercent/10) & 0x07
	percentRemaining := uint8(b.PercentRemaining/10) & 0x0F

	// set the battery info
	p.Battery.Percent = (percentRemaining << 3) | tolerancePercent
	p.Battery.Days = b.DaysRemaining
	p.Battery.Age = b.AgeDays
}

// SetTemperature adds the given value as a temperature telemetry value
func (p *Packet) SetTemperature(v float32) error {
	// create the temperature telemetry struct
	t := NewTemperatureTelemetry(v)

	// write the temperature telemetry struct to the telemetry group buffer
	err := binary.Write(&p.TelemetryData, binary.BigEndian, t)

	// return whether an error occurred
	return err
}

// SetButtonPressed adds the 'BUTTON=PRESSED' status telemetry string
func (p *Packet) SetButtonPressed() error {
	// create the status telemetry instance
	t := NewStatusTelemetry(ButtonPressedTelemetry)

	// write the button-pressed telemetry to the telemetry buffer
	err := p.WriteStatusTelemetry(t)

	// return whether an error occurred
	return err
}

// SetUtilState adds the utility state status telemetry string
func (p *Packet) SetUtilState(state string) error {
	// create the status telemetry instance
	t := NewStatusTelemetry(state)

	// write the util-state telemetry to the telemetry buffer
	err := p.WriteStatusTelemetry(t)

	// return whether an error occurred
	return err
}

// WriteStatusTelemetry writes the given status telemetry instance to the telemetry group buffer
func (p *Packet) WriteStatusTelemetry(t StatusTelemetry) error {
	// get the telemetry group buffer
	buf := p.TelemetryData

	// add each status telemetry field and return any errors
	if err := buf.WriteByte(t.ID); err != nil {
		return err
	}
	if err := buf.WriteByte(t.Length); err != nil {
		return err
	}
	if err := buf.WriteByte(t.Type); err != nil {
		return err
	}
	if _, err := buf.WriteString(t.Status); err != nil {
		return err
	}

	// return successfully
	return nil
}

// Pack returns the binary packet data structure as a slice of bytes
func (p *Packet) Pack() ([]byte, error) {
	// create the encoding buffer
	buf := &bytes.Buffer{}

	// pack the sequence number
	if err := binary.Write(buf, binary.BigEndian, p.Sequence); err != nil {
		return []byte{}, err
	}

	// pack the header bytes
	if err := binary.Write(buf, binary.BigEndian, p.Header); err != nil {
		return []byte{}, err
	}

	// pack the 'system group' bytes
	if err := binary.Write(buf, binary.BigEndian, p.System); err != nil {
		return []byte{}, err
	}

	// pack the 'battery group' bytes
	if err := binary.Write(buf, binary.BigEndian, p.Battery); err != nil {
		return []byte{}, err
	}

	// pack the 'telemetry group' bytes
	tData := p.TelemetryData.Bytes()
	buf.WriteByte(uint8(TelemetryGroupID))
	buf.WriteByte(uint8(len(tData)))
	if err := binary.Write(buf, binary.BigEndian, tData); err != nil {
		return []byte{}, err
	}

	// return the assembled bytes
	return buf.Bytes(), nil
}
