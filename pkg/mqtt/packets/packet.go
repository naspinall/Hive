package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Control packet types
const (
	Reserved    = 0  //Reserved
	CONNECT     = 1  //Connection Request
	CONNACK     = 2  //Connect Acknowledgment
	PUBLISH     = 3  //Publish Message
	PUBACK      = 4  //Publish Acknowledgment
	PUBREC      = 5  //Publish Receieved
	PUBREL      = 6  //Publish Release
	PUBCOMP     = 7  //Publish Complete
	SUBSCRIBE   = 8  //Subscribe Request
	SUBACK      = 9  //Subscribe Acknowlegement
	UNSUBSCRIBE = 10 //Unsubscribe Request
	UNSUBACK    = 11 //Unsubscribe Acknowledgement
	PINGREQ     = 12 //Ping Request
	PINGRESP    = 13 //Ping Response
	DISCONNECT  = 14 //Disconnect Notification
	AUTH        = 15 //Authentication Exchange
)

// Connect Return Code Values
const (
	ConnectionAccepted          = 0x00 //Connection accepted
	UnnaceptableProtocolVersion = 0x01 //The Server does not support the level of the MQTT protocol requested by the Client
	IdentifierRejected          = 0x02 //The Client identifier is correct UTF-8 but not allowed by the Server
	ServerUnavailable           = 0x03 //The Network Connection has been made but the MQTT service is unavailable
	BadUsernameOrPassword       = 0x04 //The data in the user name or password is malformed
	NotAuthorised               = 0x05 //The Client is not authorized to connect
)

type PacketIdentifier struct {
	Packet
	PacketIdentifier uint16
}

type StringPair struct {
	name  string
	value string
}

func (pi *PacketIdentifier) DecodePacketIdentifier() error {
	pi.PacketIdentifier = pi.DecodeTwoByteInt()
	return nil
}

func (pi *PacketIdentifier) EncodePacketIdentifier() error {
	pi.EncodeTwoByteInt(pi.PacketIdentifier)
	return nil
}

type FixedHeaderFlags struct {
	Duplicate bool
	QoS       uint8
	Retain    bool
}

type Packet struct {
	Type           uint8
	Flags          FixedHeaderFlags
	RemaningLength int
	buff           *bytes.Buffer
}

// Gets how many bytes the remaining length has taken
func (p *Packet) RemainingLengthByteLength() int {
	if p.RemaningLength < 128 {
		return 1
	} else if p.RemaningLength < 16384 {
		return 2
	} else if p.RemaningLength < 2097152 {
		return 3
	} else {
		return 4
	}
}

func FromReader(reader io.Reader) (*Packet, error) {
	b := make([]byte, 5)
	_, err := reader.Read(b) // TOOD Bad Packet Read

	Type, Flags := DecodeTypeAndFlags(b[0])
	rl, err := DecodeVariableByteInteger(b[1:])

	if err != nil {
		return nil, err
	}

	p := &Packet{
		Type:           Type,
		Flags:          Flags,
		RemaningLength: rl,
	}

	offset := 5 - p.RemainingLengthByteLength()
	buff := make([]byte, p.RemaningLength-offset)
	// Getting bytes from the start of the
	buff = append(b[offset:], buff...)
	p.buff = bytes.NewBuffer(buff)

	return p, nil
}

func NewMQTTPacket(b []byte) (*Packet, error) {

	p := &Packet{
		buff: bytes.NewBuffer(b),
	}

	if err := p.DecodeTypeAndFlags(); err != nil {
		return nil, err
	}

	if err := p.DecodeRemainingLength(); err != nil {
		return nil, err
	}

	return p, nil

}

func (p *Packet) Write(b []byte) (int, error) {
	return p.buff.Write(b)
}

func (p *Packet) DecodeRemainingLength() error {
	rl, err := p.DecodeVariableByteInteger()
	if err != nil {
		return err
	}
	p.RemaningLength = rl
	return nil
}

func (p *Packet) DecodeTypeAndFlags() error {
	b, err := p.buff.ReadByte()
	if err != nil {
		return err
	}
	p.Type = b >> 4
	p.Flags = FixedHeaderFlags{}
	p.Flags.Duplicate = (b >> 3 & 0x01) > 0
	p.Flags.QoS = uint8(b >> 1 & 0x03)
	p.Flags.Retain = b&0x01 > 0
	return nil
}

func DecodeTypeAndFlags(b byte) (byte, FixedHeaderFlags) {
	Type := b >> 4
	Flags := FixedHeaderFlags{}
	Flags.Duplicate = (b >> 3 & 0x01) > 0
	Flags.QoS = uint8(b >> 1 & 0x03)
	Flags.Retain = b&0x01 > 0
	return Type, Flags
}

func (p *Packet) DecodeVariableByteInteger() (int, error) {
	m := 1
	v := 0
	n := 0
	for {
		// Reading next byte from bufer.
		eb, err := p.buff.ReadByte()
		if err != nil {
			return 0, err
		}
		v += (int(eb) & 0x7F) * m
		m *= 128
		if m > 128*128*128 {
			return -1, errors.New("Malformed byte")
		}
		n++
		if eb&0x80 == 0 {
			break
		}
	}

	return v, nil
}

func DecodeVariableByteInteger(b []byte) (int, error) {
	m := 1
	v := 0
	n := 0
	for _, eb := range b {

		v += (int(eb) & 0x7F) * m
		m *= 128
		if m > 128*128*128 {
			return -1, errors.New("Malformed byte")
		}
		n++
		if eb&0x80 == 0 {
			break
		}
	}

	return v, nil
}

func (p *Packet) DecodeByte() (byte, error) {
	return p.buff.ReadByte()
}

func (p *Packet) DecodeFourByteInt() uint32 {
	b := p.buff.Next(4)
	return binary.BigEndian.Uint32(b)
}

func (p *Packet) DecodeTwoByteInt() uint16 {
	b := p.buff.Next(2)
	return binary.BigEndian.Uint16(b)
}

func (p *Packet) DecodeString() string {
	length := p.DecodeTwoByteInt()
	b := p.buff.Next(int(length))
	return string(b)
}

func (p *Packet) DecodeBinaryData(length int) []byte {
	return p.buff.Next(length)
}

func (p *Packet) DecodeStringPair() *StringPair {
	name := p.DecodeString()
	value := p.DecodeString()

	return &StringPair{
		name:  name,
		value: value,
	}
}

func (p *Packet) EncodeVariableByteInteger(x int) error {
	var vbi []byte

	for {
		eb := byte(x % 128)
		x /= 128
		if x > 0 {
			eb = eb | 128
		}
		vbi = append(vbi, eb)
		if x == 0 {
			break
		}
	}
	_, err := p.buff.Write(vbi)
	return err
}

func (p *Packet) EncodeByte(nb byte) error {
	return p.buff.WriteByte(nb)
}

func (p *Packet) EncodeFourByteInt(ni uint32) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, ni)
	_, err := p.buff.Write(buf)
	return err
}

func (p *Packet) EncodeTwoByteInt(ni uint16) error {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, ni)
	_, err := p.buff.Write(buf)
	return err
}

func (p *Packet) EncodeString(ns string) error {
	stringBytes := []byte(ns)
	length := uint16(len(stringBytes))
	p.EncodeTwoByteInt(length)
	_, err := p.buff.Write(stringBytes)
	return err
}

func (p *Packet) EncodeBinary(bin []byte) error {
	err := p.EncodeTwoByteInt(uint16(len(bin)))
	if err != nil {
		return err
	}
	_, err = p.buff.Write(bin)
	return err
}

func (p *Packet) EncodeTypeAndFlags() byte {
	tf := byte(p.Type) << 4

	if p.Flags.Duplicate {
		tf |= 0x08
	}
	tf |= p.Flags.QoS << 2
	if p.Flags.Retain {
		tf |= 0x01
	}
	return tf
}

// Performed at the end of each packet write.
func (p *Packet) EncodeFixedHeader() ([]byte, error) {

	tf := p.EncodeTypeAndFlags()
	rl, err := EncodeVariableByteInteger(p.buff.Len()) // Will be set when encoding the entire thing.
	if err != nil {
		return nil, err
	}
	b := p.buff.Bytes()
	// Adding remaning length
	b = append(rl, b...)
	// Adding Type and Flags
	b = append([]byte{tf}, b...)

	return b, nil

}

func EncodeVariableByteInteger(x int) ([]byte, error) {
	var vbi []byte

	for {
		eb := byte(x % 128)
		x /= 128
		if x > 0 {
			eb = eb | 128
		}
		vbi = append(vbi, eb)
		if x == 0 {
			break
		}
	}
	return vbi, nil
}
