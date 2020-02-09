package packets

import "fmt"

type PublishPacket struct {
	FixedHeader      *FixedHeader
	TopicName        string
	PacketIdentifier uint16
	Payload          []byte
}

type PublishQoSPacket struct {
	FixedHeader *FixedHeader
	PacketIdentifier
}

func NewPublishQoSPacket(fh *FixedHeader, b []byte) (*PublishQoSPacket, error) {
	pqp := &PublishQoSPacket{
		FixedHeader: fh,
	}
	_, err := pqp.DecodePacketIdentifier(b, 2)
	if err != nil {
		return nil, err
	}
	return pqp, nil
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
	pqyLength := int(pp.FixedHeader.RemaningLength) - n - 1
	fmt.Println(pqyLength)
	pp.Payload = (b[n:pqyLength])
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

func (pq *PublishQoSPacket) Encode(b []byte) ([]byte, error) {
	b, err := pq.FixedHeader.EncodeFixedHeader()
	if err != nil {
		return nil, err
	}
	return pq.EncodePacketIdentifier(b)
}

func Acknowledge(i uint16) *PublishQoSPacket {
	return &PublishQoSPacket{
		FixedHeader: &FixedHeader{
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
		FixedHeader: &FixedHeader{
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
		FixedHeader: &FixedHeader{
			Type:           6,
			RemaningLength: 2,
		},
		// This is a bit ridiculous
		PacketIdentifier: PacketIdentifier{
			PacketIdentifier: i,
		},
	}
}
