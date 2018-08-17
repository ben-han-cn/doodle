package gopacket

type Layer interface {
	LayerType() LayerType
	LayerContents() []byte
	LayerPayload() []byte
}

//TCP/IP layer 1
type LinkLayer interface {
	Layer
	LinkFlow() Flow
}

//TCP/IP layer 2
type NetworkLayer interface {
	Layer
	NetworkFlow() Flow
}

//TCP/IP layer 3
type TransportLayer interface {
	Layer
	TransportFlow() Flow
}

//TCP/IP layer 4
type ApplicationLayer interface {
	Layer
	Payload() []byte
}

type EndpointType int64

const MaxEndpointSize = 16

// Endpoint is the set of bytes used to address packets at various layers.
// See LinkLayer, NetworkLayer, and TransportLayer specifications.
// Endpoints are usable as map keys.
type Endpoint struct {
	typ EndpointType
	len int
	raw [MaxEndpointSize]byte
}

type Flow struct {
	typ        EndpointType
	slen, dlen int
	src, dst   [MaxEndpointSize]byte
}
