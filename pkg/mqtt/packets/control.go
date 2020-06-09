package packets

// TODO Add more error handling

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	*Packet
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

// All functions return the value expected and the number of bytes traversed
func DecodeVariableByteInteger(b []byte) (int, int, error) {
	m := 1
	v := 0
	n := 0
	for _, eb := range b {
		v += (int(eb) & 0x7F) * m
		m *= 128
		if m > 128*128*128 {
			return -1, 0, errors.New("Malformed byte")
		}
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

func EncodeByte(b []byte, nb byte) ([]byte, error) {
	return append(b, nb), nil
}

func EncodeFourByteInt(b []byte, ni uint32) ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, ni)
	return append(b, buf...), nil
}

func EncodeTwoByteInt(b []byte, ni uint16) ([]byte, error) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, ni)
	return append(b, buf...), nil
}

func EncodeString(b []byte, ns string) ([]byte, error) {
	stringBuf := bytes.NewBufferString(ns).Bytes()
	strLen := uint16(len(stringBuf))
	b, err := EncodeTwoByteInt(b, strLen)
	return append(b, stringBuf...), err
}

func EncodeBinary(b []byte, bin []byte) ([]byte, error) {
	b, err := EncodeTwoByteInt(b, uint16(len(bin)))
	return append(b, bin...), err

}
