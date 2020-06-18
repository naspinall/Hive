package packets

type DisconnectPacket struct {
	Packet
}

func NewDisconnectPacket(p *Packet) (*DisconnectPacket, error) {
	dp := &DisconnectPacket{*p}
	return dp, nil
}

func (dp *DisconnectPacket) Encode() ([]byte, error) {
	return dp.EncodeFixedHeader()
}
