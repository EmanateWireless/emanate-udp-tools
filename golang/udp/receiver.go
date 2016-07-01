package udp

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Receiver is the UDP server instance
type Receiver struct {
	options *ReceiverOptions
}

// ReceiverOptions provides the instance options
type ReceiverOptions struct {
	Port int
}

// NewReceiver creates a new instance
func NewReceiver(options *ReceiverOptions) *Receiver {
	// return the new instance
	return &Receiver{
		options: options,
	}
}

// Run starts the UDP receiver instance
func (r *Receiver) Run() {
	fmt.Printf("Starting UDP receiver listening on port '%d'\n", r.options.Port)

	// start listening to udp packet
	socket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: r.options.Port,
	})

	// if an error occurred
	if err != nil {
		// log the error and exit the process now
		fmt.Printf("\nError starting UDP receiver listening to UDP port '%d' (error = '%v')\n\n",
			r.options.Port, err)
		os.Exit(1)
	}

	// close the udp socket when done
	defer socket.Close()

	// create the udp receive buffer
	buf := make([]byte, 2048)

	// wait for received udp packets until commanded to quit
	for {
		// read the next udp packet
		_, remoteaddr, err := socket.ReadFromUDP(buf)
		if err != nil {
			// log the error and exit the process now
			fmt.Printf("\nError occurred while receiving UDP packet (error = '%v'", err)
			os.Exit(1)
		}

		// get the current time
		now := time.Now()

		fmt.Printf("[%s] %s: '%s'\n", now, remoteaddr, buf)
	}
}
