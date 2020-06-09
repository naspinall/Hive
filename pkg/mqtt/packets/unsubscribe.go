package packets

type UnsubAckPacket struct {
	*Packet
	PacketIdentifier
}

func NewUnsubAckPacket(p *Packet, b []byte) (*UnsubAckPacket, error) {
	usap := &UnsubAckPacket{
		Packet: p,
	}
	if err := usap.DecodePacketIdentifier(); err != nil {
		return nil, err
	}
	return usap, nil
}

func (uap *UnsubAckPacket) Encode() ([]byte, error) {
	if err := uap.EncodePacketIdentifier(); err != nil {
		return nil, err
	}
	return uap.EncodeFixedHeader()
}
