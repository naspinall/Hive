package packets

type Topic struct {
	Topic string
	QoS   byte
}

type SubscribePacket struct {
	FixedHeader *FixedHeader
	PacketIdentifier

	//Payload Properties
	Topics []Topic
}

type SubAckPacket struct {
	FixedHeader *FixedHeader
	PacketIdentifier
	ReturnCode byte
}

func NewSubAckPacket(fh *FixedHeader, b []byte) (*SubAckPacket, error) {
	sap := &SubAckPacket{}
	n, err := sap.DecodePacketIdentifier(b, 2)
	if err != nil {
		return nil, err
	}
	sap.ReturnCode, n, err = DecodeByte(b[n:])
	return sap, nil

}

func NewSubscribePacket(fh *FixedHeader, b []byte) (*SubscribePacket, error) {
	sp := &SubscribePacket{}
	n, err := sp.DecodePacketIdentifier(b, 2)
	if err = sp.DecodeTopics(b, n); err != nil {
		return nil, err
	}
	return sp, nil
}

func (sp *SubscribePacket) DecodeTopics(b []byte, n int) (err error) {
	for n < int(sp.FixedHeader.RemaningLength) {
		topic, n, err := DecodeString(b[n:])
		if err != nil {
			return err
		}
		qos, n, err := DecodeByte(b[n:])
		sp.Topics = append(sp.Topics, Topic{
			Topic: topic,
			QoS:   qos,
		})
	}

	return nil
}

func (sp *SubscribePacket) EncodeTopics(b []byte) ([]byte, error) {
	for _, topic := range sp.Topics {
		b, err := EncodeString(b, topic.Topic)
		if err != nil {
			return nil, err
		}
		b, err = EncodeByte(b, topic.QoS)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (sp *SubscribePacket) EncodeSubscribePacket() ([]byte, error) {
	var b []byte

	// Packet identifier
	b, err := sp.EncodePacketIdentifier(b)
	if err != nil {
		return nil, err
	}

	// Encode the topics
	b, err = sp.EncodeTopics(b)
	if err != nil {
		return nil, err
	}

	// Encode the QoS byte.
	b, err = EncodeByte(b, sp.FixedHeader.Flags.QoS)

	// Prepending the fixed header
	sp.FixedHeader.RemaningLength = len(b)
	return sp.FixedHeader.PrependFixedHeader(b)
}

func (sp *SubAckPacket) EncodeSubAckPacket() ([]byte, error) {
	var b []byte

	b, err := sp.EncodePacketIdentifier(b)
	if err != nil {
		return nil, err
	}

	b, err = sp.EncodePacketIdentifier(b)
	if err != nil {
		return nil, err
	}

	b, err = EncodeByte(b, sp.ReturnCode)
	if err != nil {
		return nil, err
	}

	// TODO finish this.
	return b, err
}
