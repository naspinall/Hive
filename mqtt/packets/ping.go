package packets

type PingRequestPacket struct {
	FixedHeader *FixedHeader
}

type PingResponsePacket struct {
	PingRequestPacket
}

func NewPingRequestPacket(fh *FixedHeader) *PingRequestPacket {
	return &PingRequestPacket{FixedHeader: fh}
}

func NewPingResponsePacket(fh *FixedHeader) *PingResponsePacket {
	return &PingResponsePacket{}
}

func (prp *PingRequestPacket) EncodePingRequestPacket() ([]byte, error) {
	return prp.FixedHeader.EncodeFixedHeader()
}
func (prp *PingResponsePacket) EncodePingResponsePacket() ([]byte, error) {
	return prp.FixedHeader.EncodeFixedHeader()
}
