package udp

// Encoder instance struct
type Encoder struct {
	options *EncoderOptions
}

// EncoderOptions provide the instance options
type EncoderOptions struct {
}

// NewEncoder creates a new instance
func NewEncoder(options *EncoderOptions) *Encoder {
	// return the new instance
	return &Encoder{
		options: options,
	}
}
