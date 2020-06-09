package packets

type PingPacket struct {
	*Packet
}

func NewPingPacket() *PingPacket {
	return &PingPacket{}
}

func (prp PingPacket) Encode() ([]byte, error) {
	return prp.EncodeFixedHeader()
}

func PingRequest() *PingPacket {
	return &PingPacket{
		Packet: &Packet{
			Type:           12,
			RemaningLength: 0,
		},
	}
}

func PingResponse() *PingPacket {
	return &PingPacket{
		Packet: &Packet{
			Type:           13,
			RemaningLength: 0,
		},
	}
}
