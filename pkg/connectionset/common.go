package connectionset

type Protocol interface {
	// InverseDirection returns the response expected for a request made using this protocol
	InverseDirection() Protocol
}

type AnyProtocol struct{}

func (t AnyProtocol) InverseDirection() Protocol { return AnyProtocol{} }

type ProtocolStr string

const (
	ProtocolStringTCP  ProtocolStr = "TCP"
	ProtocolStringUDP  ProtocolStr = "UDP"
	ProtocolStringICMP ProtocolStr = "ICMP"
)
