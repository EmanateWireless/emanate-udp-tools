package udp

import (
	"fmt"
	"log"
	"net"
)

// Sender is the UDP transmitter instance
type Sender struct {
	options *SenderOptions
}

// SenderOptions provides the instance options
type SenderOptions struct {
	Host string
	Port int
}

// NewSender creates a new instance
func NewSender(options *SenderOptions) *Sender {
	// return the new instance
	return &Sender{
		options: options,
	}
}

// Transmit sends the given message as a UDP packet to the configured destination
func (s *Sender) Transmit(data []byte) {
	dst := fmt.Sprintf("%s:%d", s.options.Host, s.options.Port)
	log.Printf("Sending udp packet to '%s' (%d bytes)", dst, len(data))

	// create the udp socket
	conn, err := net.Dial("udp", dst)
	if err != nil {
		// log the error and return now
		fmt.Printf("Error creating udp socket (error = '%v')\n", err)
		return
	}

	// close the udp socket when finished
	defer conn.Close()

	// send the udp packet
	_, err = conn.Write(data)
	if err != nil {
		// log the error and return now
		log.Printf("Error sending udp packet (error = '%v')", err)
	} else {
		log.Printf("Successfully sent UDP packet to '%s' (%d bytes)\n", dst, len(data))
	}
}
