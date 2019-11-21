package packets

type PublishPacket struct {
	FixedHeader      *FixedHeader
	TopicName        string
	PacketIdentifier uint16
	Payload          []byte
}

type PublishAcknowledgmentPacket struct {
	FixedHeader *FixedHeader
	PacketIdentifier
}

type PublishRecievedPacket struct {
	PublishAcknowledgmentPacket
}

type PublishCompletePacket struct {
	PublishAcknowledgmentPacket
}

func NewPublishAcknowledgmentPacket(fh *FixedHeader, b []byte) (*PublishAcknowledgmentPacket, error) {
	pap := &PublishAcknowledgmentPacket{
		FixedHeader: fh,
	}
	_, err := pap.DecodePacketIdentifier(b, 2)
	if err != nil {
		return nil, err
	}
	return pap, nil
}

func NewPublishPacket(fh *FixedHeader, b []byte) (*PublishPacket, error) {
	pp := &PublishPacket{
		FixedHeader: fh,
	}
	n, err := pp.DecodeTopicName(b, 0)
	if err != nil {
		return nil, err
	}
	n, err = pp.DecodePacketIdentifier(b, n)
	if err != nil {
		return nil, err
	}
	payLength := n - int(pp.FixedHeader.RemaningLength) - 1
	pp.Payload = (b[n:payLength])
	return pp, nil
}

func (pp *PublishPacket) DecodeTopicName(b []byte, n int) (m int, err error) {
	pp.TopicName, m, err = DecodeString(b[n:])
	return
}

func (pp *PublishPacket) EncodeTopicName(b []byte) ([]byte, error) {
	return EncodeString(b, pp.TopicName)
}

func (pp *PublishPacket) DecodePacketIdentifier(b []byte, n int) (m int, err error) {
	if pp.FixedHeader.Flags.QoS > 0 {
		pp.PacketIdentifier, m, err = DecodeTwoByteInt(b[n:])
	}
	return
}

func (pp *PublishPacket) EncodePacketIdentifier(b []byte) ([]byte, error) {
	if pp.FixedHeader.Flags.QoS > 0 {
		return EncodeTwoByteInt(b, pp.PacketIdentifier)
	}
	return b, nil
}

func (pp *PublishPacket) Encode(b []byte) ([]byte, error) {
	// Variable header starts with the topic name
	b, err := pp.EncodeTopicName(b)
	if err != nil {
		return nil, err
	}
	// Packet identifier next
	b, err = pp.EncodePacketIdentifier(b)
	if err != nil {
		return nil, err
	}

	// Payload next
	b = append(b, pp.Payload...)

	fhb, err := pp.FixedHeader.EncodeFixedHeader()
	if err != nil {
		return nil, err
	}
	return append(fhb, b...), err
}

func (pa *PublishAcknowledgmentPacket) Encode(b []byte) ([]byte, error) {
	b, err := pa.FixedHeader.EncodeFixedHeader()
	if err != nil {
		return nil, err
	}
	return pa.EncodePacketIdentifier(b)
}
