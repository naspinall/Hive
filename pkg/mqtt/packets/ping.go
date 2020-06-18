package packets

import "bytes"

type PingPacket struct {
	Packet
}

func NewPingPacket() *PingPacket {
	return &PingPacket{}
}

func (prp PingPacket) Encode() ([]byte, error) {
	return prp.EncodeFixedHeader()
}

func PingRequest() *PingPacket {
	return &PingPacket{
		Packet{
			Type:           12,
			RemaningLength: 0,
			buff:           &bytes.Buffer{},
		},
	}
}

func PingResponse() *PingPacket {
	return &PingPacket{
		Packet{
			Type:           13,
			RemaningLength: 0,
			buff:           &bytes.Buffer{},
		},
	}
}
