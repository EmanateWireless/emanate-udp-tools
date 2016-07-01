package udp

// Decoder instance struct
type Decoder struct {
	options *DecoderOptions
}

// DecoderOptions provide the instance options
type DecoderOptions struct {
}

// NewDecoder creates a new instance
func NewDecoder(options *DecoderOptions) *Decoder {
	// return the new instance
	return &Decoder{
		options: options,
	}
}
