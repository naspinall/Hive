package packets

type Topic struct {
	Topic string
	QoS   byte
}

type SubscribePacket struct {
	*Packet
	PacketIdentifier

	//Payload Properties
	Topics []Topic
}

type SubAckPacket struct {
	*Packet
	PacketIdentifier
	ReturnCode byte
}

func NewSubAckPacket(p *Packet) (*SubAckPacket, error) {
	sap := &SubAckPacket{
		Packet: p,
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
		Packet: p,
	}
	err := sp.DecodePacketIdentifier()
	if err = sp.DecodeTopics(); err != nil {
		return nil, err
	}
	return sp, nil
}

func (sp *SubscribePacket) DecodeTopics() error {
	// Reconsider this.
	for sp.buff.Len() > 0 {
		topic := sp.DecodeString()
		qos, err := sp.DecodeByte()
		if err != nil {
			return err
		}
		sp.Topics = append(sp.Topics, Topic{
			Topic: topic,
			QoS:   qos,
		})
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

	if err := sp.EncodePacketIdentifier(); err != nil {
		return nil, err
	}

	if err := sp.EncodeByte(sp.ReturnCode); err != nil {
		return nil, err
	}

	return sp.EncodeFixedHeader()
}
