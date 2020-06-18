package packets

import (
	"bytes"
)

type Topic struct {
	Topic string
	QoS   byte
}

type SubscribePacket struct {
	PacketIdentifier

	//Payload Properties
	Topics []Topic
}

type SubAckPacket struct {
	PacketIdentifier
	ReturnCode byte
}

func NewSubAckPacket(p *Packet) (*SubAckPacket, error) {
	sap := &SubAckPacket{
		PacketIdentifier: PacketIdentifier{
			Packet: *p,
		},
	}
	err := sap.DecodePacketIdentifier()
	if err != nil {
		return nil, err
	}
	sap.ReturnCode, err = sap.DecodeByte()
	return sap, nil

}

func NewSubscribePacket(p *Packet) (*SubscribePacket, error) {
	sp := &SubscribePacket{
		PacketIdentifier: PacketIdentifier{
			Packet: *p,
		},
	}
	err := sp.DecodePacketIdentifier()
	if err = sp.DecodeTopics(); err != nil {
		return nil, err
	}
	return sp, nil
}

func (sp *SubscribePacket) DecodeTopics() error {
	var i int
	topicsLen := sp.Packet.RemaningLength - 2
	// Reconsider this.
	for i < topicsLen {

		length := sp.DecodeTwoByteInt()
		b := sp.buff.Next(int(length))

		topic := string(b)
		qos, err := sp.DecodeByte()
		if err != nil {
			return err
		}
		sp.Topics = append(sp.Topics, Topic{
			Topic: topic,
			QoS:   qos,
		})
		i += int(length + 1)

	}

	return nil
}

func (sp *SubscribePacket) EncodeTopics() error {
	for _, topic := range sp.Topics {
		if err := sp.EncodeString(topic.Topic); err != nil {

			return err
		}
		if err := sp.EncodeByte(topic.QoS); err != nil {
			return err
		}

	}
	return nil
}

func (sp *SubscribePacket) Encode() ([]byte, error) {

	// Packet identifier
	if err := sp.EncodePacketIdentifier(); err != nil {
		return nil, err
	}

	// Encode the topics
	if err := sp.EncodeTopics(); err != nil {
		return nil, err
	}

	// Encode the QoS byte.
	sp.EncodeByte(sp.Flags.QoS)

	return sp.EncodeFixedHeader()
}

func (sp *SubAckPacket) Encode() ([]byte, error) {
	if err := sp.EncodePacketIdentifier(); err != nil {
		return nil, err
	}

	if err := sp.EncodeByte(sp.ReturnCode); err != nil {
		return nil, err
	}

	return sp.EncodeFixedHeader()
}

func SubAck(packetIdentifier uint16, rc byte) *SubAckPacket {
	p := &Packet{
		buff:           &bytes.Buffer{},
		Type:           SUBACK,
		RemaningLength: 3,
	}
	return &SubAckPacket{
		PacketIdentifier: PacketIdentifier{
			Packet:           *p,
			PacketIdentifier: packetIdentifier,
		},
	}
}
