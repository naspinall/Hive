package packets

type PublishPacket struct {
	*Packet
	TopicName        string
	PacketIdentifier uint16
	Payload          []byte
}

type PublishQoSPacket struct {
	*Packet
	PacketIdentifier
}

func NewPublishQoSPacket(p *Packet) (*PublishQoSPacket, error) {
	pqp := &PublishQoSPacket{
		Packet: p,
	}
	err := pqp.DecodePacketIdentifier()
	if err != nil {
		return nil, err
	}
	return pqp, nil
}

func NewPublishPacket(p *Packet) (*PublishPacket, error) {
	pp := &PublishPacket{
		Packet: p,
	}
	err := pp.DecodeTopicName()
	if err != nil {
		return nil, err
	}

	err = pp.DecodePacketIdentifier()
	if err != nil {
		return nil, err
	}

	pp.Payload = pp.DecodeBinaryData()
	return pp, nil
}

func (pp *PublishPacket) DecodeTopicName() error {
	pp.TopicName = pp.DecodeString()
	return nil
}

func (pp *PublishPacket) EncodeTopicName(b []byte) ([]byte, error) {
	return EncodeString(b, pp.TopicName)
}

func (pp *PublishPacket) DecodePacketIdentifier() error {
	if pp.Flags.QoS > 0 {
		pp.PacketIdentifier = pp.DecodeTwoByteInt()
	}
	return nil
}

func (pp *PublishPacket) EncodePacketIdentifier(b []byte) ([]byte, error) {
	if pp.Flags.QoS > 0 {
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

	return pp.EncodeFixedHeader()
}

func (pq *PublishQoSPacket) Encode() ([]byte, error) {
	pq.EncodeTwoByteInt(pq.PacketIdentifier.PacketIdentifier)
	return pq.EncodeFixedHeader()
}

func Acknowledge(i uint16) *PublishQoSPacket {
	return &PublishQoSPacket{
		Packet: &Packet{
			Type:           4,
			RemaningLength: 2,
		},
		// This is a bit ridiculous
		PacketIdentifier: PacketIdentifier{
			PacketIdentifier: i,
		},
	}
}

func Received(i uint16) *PublishQoSPacket {
	return &PublishQoSPacket{
		Packet: &Packet{
			Type:           5,
			RemaningLength: 2,
		},
		// This is a bit ridiculous
		PacketIdentifier: PacketIdentifier{
			PacketIdentifier: i,
		},
	}
}

func Complete(i uint16) *PublishQoSPacket {
	return &PublishQoSPacket{
		Packet: &Packet{
			Type:           6,
			RemaningLength: 2,
		},
		// This is a bit ridiculous
		PacketIdentifier: PacketIdentifier{
			PacketIdentifier: i,
		},
	}
}
