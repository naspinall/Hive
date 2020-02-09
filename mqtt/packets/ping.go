package packets

type PingPacket struct {
	FixedHeader *FixedHeader
}

func NewPingPacket(fh *FixedHeader) *PingPacket {
	return &PingPacket{FixedHeader: fh}
}

func (prp PingPacket) Encode() ([]byte, error) {
	return prp.FixedHeader.EncodeFixedHeader()
}

func PingRequest() *PingPacket {
	return &PingPacket{
		FixedHeader: &FixedHeader{
			Type:           12,
			RemaningLength: 0,
		},
	}
}

func PingResponse() *PingPacket {
	return &PingPacket{
		FixedHeader: &FixedHeader{
			Type:           13,
			RemaningLength: 0,
		},
	}
}
