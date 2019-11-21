package packets

type UnsubscribePacket struct {
	SubscribePacket
}
type UnsubAckPacket struct {
	FixedHeader *FixedHeader
	PacketIdentifier
}

func NewUnsubAckPacket(fh *FixedHeader, b []byte) (*UnsubAckPacket, error) {
	usap := &UnsubAckPacket{}
	usap.FixedHeader = fh
	n := 2
	_, err := usap.DecodePacketIdentifier(b, n)
	if err != nil {
		return nil, err
	}
	return usap, nil
}
