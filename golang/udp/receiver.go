package udp

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Receiver is the UDP server instance
type Receiver struct {
	options     *ReceiverOptions
	dataHandler DataUpdateFunc
}

// ReceiverOptions provides the instance options
type ReceiverOptions struct {
	Port int
}

// DataUpdate defines the UDP data update structure passed to the registered data handler
type DataUpdate struct {
	TS         time.Time
	RemoteIP   string
	RemotePort int
	Data       []byte
}

// DataUpdateFunc is the callback function type used to notify when UDP data is received
type DataUpdateFunc func(du *DataUpdate)

// NewReceiver creates a new instance
func NewReceiver(options *ReceiverOptions) *Receiver {
	// return the new instance
	return &Receiver{
		options: options,
	}
}

// DataHandler registers the update handler to call when UDP data is received
func (r *Receiver) DataHandler(handler DataUpdateFunc) {
	// save the data handler
	r.dataHandler = handler
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
		fmt.Printf("Error starting UDP receiver listening to UDP port '%d' (error = '%v')\n\n",
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
		_, remoteAddr, err := socket.ReadFromUDP(buf)
		if err != nil {
			// log the error and exit the process now
			fmt.Printf("Error occurred while receiving UDP packet (error = '%v'", err)
			os.Exit(1)
		}

		// get the current time
		now := time.Now()

		// call the update handler if registered
		if r.dataHandler != nil {
			r.dataHandler(&DataUpdate{
				TS:         now,
				RemoteIP:   remoteAddr.IP.String(),
				RemotePort: remoteAddr.Port,
				Data:       buf,
			})
		}
	}
}
