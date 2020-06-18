package packets

import (
	"bytes"
)

type PublishPacket struct {
	Packet
	TopicName        string
	PacketIdentifier uint16
	Payload          []byte
}

type PublishQoSPacket struct {
	Packet
	PacketIdentifier
}

func NewPublishQoSPacket(p *Packet) (*PublishQoSPacket, error) {
	pqp := &PublishQoSPacket{
		Packet: *p,
	}
	err := pqp.DecodePacketIdentifier()
	if err != nil {
		return nil, err
	}
	return pqp, nil
}

func NewPublishPacket(p *Packet) (*PublishPacket, error) {
	pp := &PublishPacket{
		Packet: *p,
	}

	err := pp.DecodeTopicName()
	if err != nil {
		return nil, err
	}

	if pp.Flags.QoS > 0 {
		err = pp.DecodePacketIdentifier()
		if err != nil {
			return nil, err
		}
	}

	pp.Payload = pp.buff.Next(pp.buff.Len())
	return pp, nil
}

func (pp *PublishPacket) DecodeTopicName() error {
	pp.TopicName = pp.DecodeString()
	return nil
}

func (pp *PublishPacket) EncodeTopicName() error {
	pp.EncodeString(pp.TopicName)
	return nil
}

func (pp *PublishPacket) DecodePacketIdentifier() error {
	if pp.Flags.QoS > 0 {
		pp.PacketIdentifier = pp.DecodeTwoByteInt()
	}
	return nil
}

func (pp *PublishPacket) EncodePacketIdentifier() error {
	if pp.Flags.QoS > 0 {
		return pp.EncodeTwoByteInt(pp.PacketIdentifier)
	}
	return nil
}

func (pp *PublishPacket) Encode() ([]byte, error) {
	// Variable header starts with the topic name
	if err := pp.EncodeTopicName(); err != nil {
		return nil, err
	}

	if pp.Flags.QoS > 0 {
		// Packet identifier next
		if err := pp.EncodePacketIdentifier(); err != nil {
			return nil, err
		}
	}

	pp.EncodeBinary(pp.Payload)
	return pp.EncodeFixedHeader()
}

func (pq *PublishQoSPacket) Encode() ([]byte, error) {
	pq.EncodeTwoByteInt(pq.PacketIdentifier.PacketIdentifier)
	return pq.EncodeFixedHeader()
}

func Acknowledge(i uint16) *PublishQoSPacket {
	return &PublishQoSPacket{
		Packet: Packet{
			Type:           4,
			RemaningLength: 2,
			buff:           &bytes.Buffer{},
		},
		// This is a bit ridiculous
		PacketIdentifier: PacketIdentifier{
			PacketIdentifier: i,
		},
	}
}

func Received(i uint16) *PublishQoSPacket {
	return &PublishQoSPacket{
		Packet: Packet{
			Type:           5,
			RemaningLength: 2,
			buff:           &bytes.Buffer{},
		},
		// This is a bit ridiculous
		PacketIdentifier: PacketIdentifier{
			PacketIdentifier: i,
		},
	}
}

func Complete(i uint16) *PublishQoSPacket {
	return &PublishQoSPacket{
		Packet: Packet{
			Type:           6,
			RemaningLength: 2,
			buff:           &bytes.Buffer{},
		},
		// This is a bit ridiculous
		PacketIdentifier: PacketIdentifier{
			PacketIdentifier: i,
		},
	}
}
