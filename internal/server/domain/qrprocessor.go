package domain

// QRDecoder is an interface for any kind of QR decoders
type QRDecoder interface {
	Decode([]byte) ([]byte, error)
}

// QREncoder is an interface for any kind of QR encoders
type QREncoder interface {
	Encode([]byte) ([]byte, error)
}
