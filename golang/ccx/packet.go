package ccx

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/EmanateWireless/emanate-udp-tools/golang/util"
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
	ProbeUnpluggedTelemetry     = "TEMP_PROBE_ERROR=UNPLUGGED"
	ProbeInvalidValueTelemetry  = "TEMP_PROBE_ERROR=INVALID_VALUE"
	TelemetryDataOffset         = 34
)

// Packet instance struct
type Packet struct {
	EmanateHeader EmanateHeader
	Header        Header
	System        SystemGroup
	Battery       BatteryGroup
	TelemetryData bytes.Buffer
}

// ParsedPacket includes the fields that can be parsed by the encoding/binary package
type ParsedPacket struct {
	EmanateHeader EmanateHeader
	Header        Header
	System        SystemGroup
	Battery       BatteryGroup
}

// EmanateHeader provides an encapsulating header to the CCX-formatted packet
type EmanateHeader struct {
	UDPVersion uint16
	TagMACAddr [6]uint8
	APMACAddr  [6]uint8
	Sequence   uint16
}

// Header defines the CCX packet header
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
	ProductType uint16
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
	ID           uint8
	Length       uint8
	Type         uint8
	StatusLength uint8
	Status       string
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
	// get the tag and AP mac-address bytes
	tagMACAddr, _ := util.MACAddrToBytes("11:22:33:44:55:66")
	apMACAddr, _ := util.MACAddrToBytes("66:55:44:33:22:11")

	// return a new default packet instance
	return &Packet{
		EmanateHeader: EmanateHeader{
			UDPVersion: uint16(0),
			TagMACAddr: tagMACAddr,
			APMACAddr:  apMACAddr,
			Sequence:   uint16(1),
		},
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

	// get the status string length
	statusLen := uint8(len(encodedStr))

	// create and return the status telemetry instance
	return StatusTelemetry{
		ID:           TelemetryGroupID,    // 3
		Length:       statusLen + 2,       // variable group length
		Type:         StatusTelemetryType, // 8
		StatusLength: statusLen,           // variable status length
		Status:       encodedStr,          // utf-16 status string
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

// SetUDPVersion sets the UDP version in the encapsulating Emanate header
func (p *Packet) SetUDPVersion(v uint16) {
	p.EmanateHeader.UDPVersion = v
}

// SetTagMACAddress sets the tag MAC address in the encapsulating Emanate header
func (p *Packet) SetTagMACAddress(mac string) error {
	// convert the mac-address string into a 6-byte array
	macBytes, err := util.MACAddrToBytes(mac)

	// if conversion successful
	if err == nil {
		p.EmanateHeader.TagMACAddr = macBytes
	}

	// return whether an error occurred
	return err
}

// SetAPMACAddress sets the wifi AP MAC address in the encapsulating Emanate header
func (p *Packet) SetAPMACAddress(mac string) error {
	// convert the mac-address string into a 6-byte array
	macBytes, err := util.MACAddrToBytes(mac)

	// if conversion successful
	if err == nil {
		p.EmanateHeader.APMACAddr = macBytes
	}

	// return whether an error occurred
	return err
}

// SetSequenceNumber sets the packet sequence number
func (p *Packet) SetSequenceNumber(seq uint16) {
	p.EmanateHeader.Sequence = seq
}

// IncSequenceNumber increments the packet sequence number
func (p *Packet) IncSequenceNumber(seq uint16) {
	// increment the packet sequence number by 1
	p.IncSequenceNumberByValue(1)
}

// IncSequenceNumberByValue increments the packet sequence number by the given value
func (p *Packet) IncSequenceNumberByValue(v uint16) {
	p.EmanateHeader.Sequence = p.EmanateHeader.Sequence + v
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

// SetProductType sets the system group product-type with the given values
func (p *Packet) SetProductType(pt uint16) {
	// set the product-type
	p.System.ProductType = pt
}

// SetBatteryInfo sets the battery group with the given values
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

// SetDoorOpenPercent adds the door-open telemetry string
func (p *Packet) SetDoorOpenPercent(percent int) error {
	// build the telemetry status string
	status := fmt.Sprintf("DOOR_OPEN_PERCENT=%d", percent)

	// create the status telemetry instance
	t := NewStatusTelemetry(status)

	// write the button-pressed telemetry to the telemetry buffer
	err := p.WriteStatusTelemetry(t)

	// return whether an error occurred
	return err
}

// SetHighPowerPercent adds the high-power-mode telemetry string
func (p *Packet) SetHighPowerPercent(percent int) error {
	// build the telemetry status string
	status := fmt.Sprintf("HIGH_POWER_MODE_PERCENT=%d", percent)

	// create the status telemetry instance
	t := NewStatusTelemetry(status)

	// write the button-pressed telemetry to the telemetry buffer
	err := p.WriteStatusTelemetry(t)

	// return whether an error occurred
	return err
}

// SetButtonPressed adds the button-pressed status telemetry string
func (p *Packet) SetButtonPressed() error {
	// create the status telemetry instance
	t := NewStatusTelemetry(ButtonPressedTelemetry)

	// write the status telemetry entry to the telemetry buffer
	err := p.WriteStatusTelemetry(t)

	// return whether an error occurred
	return err
}

// SetProbeUnplugged adds the temperature probe unplugged error status telemetry string
func (p *Packet) SetProbeUnplugged() error {
	// create the status telemetry instance
	t := NewStatusTelemetry(ProbeUnpluggedTelemetry)

	// write the status telemetry entry to the telemetry buffer
	err := p.WriteStatusTelemetry(t)

	// return whether an error occurred
	return err
}

// SetProbeInvalidValue adds the temperature probe read error status telemetry string
func (p *Packet) SetProbeInvalidValue() error {
	// create the status telemetry instance
	t := NewStatusTelemetry(ProbeInvalidValueTelemetry)

	// write the status telemetry entry to the telemetry buffer
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
	// add each status telemetry field and return any errors
	if err := p.TelemetryData.WriteByte(t.ID); err != nil {
		return err
	}
	if err := p.TelemetryData.WriteByte(t.Length); err != nil {
		return err
	}
	if err := p.TelemetryData.WriteByte(t.Type); err != nil {
		return err
	}
	if err := p.TelemetryData.WriteByte(t.StatusLength); err != nil {
		return err
	}
	if _, err := p.TelemetryData.WriteString(t.Status); err != nil {
		return err
	}

	// return successfully
	return nil
}

// Pack returns the binary packet data structure as a slice of bytes
func (p *Packet) Pack() ([]byte, error) {
	// create the encoding buffer
	buf := &bytes.Buffer{}

	// pack the udp packet version
	if err := binary.Write(buf, binary.BigEndian, p.EmanateHeader.UDPVersion); err != nil {
		return []byte{}, err
	}

	// pack the tag mac address
	if err := binary.Write(buf, binary.BigEndian, p.EmanateHeader.TagMACAddr); err != nil {
		return []byte{}, err
	}

	// pack the wifi AP mac address
	if err := binary.Write(buf, binary.BigEndian, p.EmanateHeader.APMACAddr); err != nil {
		return []byte{}, err
	}

	// pack the sequence number
	if err := binary.Write(buf, binary.BigEndian, p.EmanateHeader.Sequence); err != nil {
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
	if err := binary.Write(buf, binary.BigEndian, tData); err != nil {
		return []byte{}, err
	}

	// return the assembled bytes
	return buf.Bytes(), nil
}
