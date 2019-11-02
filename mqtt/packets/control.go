package packets

// TODO Add more error handling

import (
	"encoding/binary"
	"errors"
	"log"
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

type FixedHeaderFlags struct {
	Duplicate bool
	QoS       int
	Retain    bool
}

type FixedHeader struct {
	Type           int
	Flags          FixedHeaderFlags
	RemaningLength int
}

type StringPair struct {
	name  string
	value string
}

func NewFixedHeader(typeAndFlags byte, remainingLength byte) (*FixedHeader, error) {
	fh := &FixedHeader{}
	if err := fh.GetTypeAndFlags(typeAndFlags); err != nil {
		log.Fatal("Invalid control packet")
	}
	rl, _, err := DecodeVariableByteInteger([]byte{remainingLength})
	if err != nil {
		return nil, err
	}
	fh.RemaningLength = rl
	return fh, nil
}

func (fh *FixedHeader) GetTypeAndFlags(b byte) error {
	fh.Type = int(b >> 4)
	fh.Flags = FixedHeaderFlags{}
	fh.Flags.Duplicate = (b >> 3 & 0x01) > 0
	fh.Flags.QoS = int(b >> 2 & 0x03)
	fh.Flags.Retain = b&0x01 > 0
	return nil
}

func EncodeVariableByteInteger(x int) []byte {
	var b []byte

	for {
		eb := byte(x % 128)
		x /= 128
		if x > 0 {
			eb = eb | 128
		}
		b = append(b, eb)
		if x == 0 {
			break
		}
	}
	return b
}

// All functions return the value expected and the number of bytes traversed
func DecodeVariableByteInteger(b []byte) (int, int, error) {
	m := 1
	v := 0
	n := 0
	for eb := range b {
		v += (eb & 0x7F) * m
		if m > 128*128*128 {
			return -1, 0, errors.New("Malformed byte")
		}
		m *= 128
		n++
		if eb&0x80 == 0 {
			break
		}
	}

	return v, n, nil
}

func DecodeByte(b []byte) (byte, int, error) {
	return b[0], 1, nil
}

func DecodeFourByteInt(b []byte) (uint32, int, error) {
	return binary.BigEndian.Uint32(b[0:4]), 4, nil
}

func DecodeTwoByteInt(b []byte) (uint16, int, error) {
	return binary.BigEndian.Uint16(b[0:2]), 2, nil
}

func DecodeString(b []byte) (string, int, error) {
	length, _, err := DecodeTwoByteInt(b)
	if err != nil {
		return "", 0, err
	}
	return string(b[2 : length+2]), int(length + 2), nil
}

func DecodeBinaryData(b []byte) ([]byte, int, error) {
	length, _, err := DecodeTwoByteInt(b)
	if err != nil {
		return nil, 0, err
	}
	return b[2 : length+2], int(length + 2), nil
}

func DecodeStringPair(b []byte) (*StringPair, int, error) {
	name, n, err := DecodeString(b)
	if err != nil {
		return nil, 0, err
	}
	value, m, err := DecodeString(b[n:])
	if err != nil {
		return nil, 0, err
	}
	return &StringPair{
		name:  name,
		value: value,
	}, m + n, nil
}
