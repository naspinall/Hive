package packets

import "bytes"

type WillProperties struct {
	WillDelayInterval      uint32
	PayloadFormatIndicator bool
	MessageExpiryInterval  uint32
	ContentType            string
	ResponseTopic          string
	CorrelationData        []byte
	UserProperty           *StringPair
}

type ConnectPacket struct {
	Packet
	//Variable Header
	ProtocolName    string
	ProtocolVersion byte
	KeepAlive       uint16

	// Connect Flags
	UsernameFlag   bool
	PasswordFlag   bool
	WillRetainFlag bool
	WillQoSFlag    uint8
	WillFlag       bool
	CleanStartFlag bool

	//Variable Header Properties
	SessionExpiryInterval      uint32
	AuthMethod                 string
	AuthData                   []byte
	RequestResponseInformation bool
	RequestProblemInformation  bool
	RecieveMaximum             uint16
	TopicAliasMaximum          uint16
	UserProperty               *StringPair
	MaximumPacketSize          uint32

	//Payload properties
	ClientID       string
	Username       string
	WillProperties *WillProperties
	Password       []byte
	WillTopic      string
	WillPayload    []byte
}

type ConnackPacket struct {
	Packet
	SessionPresent byte
	ReturnCode     byte
}

func NewConnackPacket(b []byte) (*ConnackPacket, error) {
	p, err := NewMQTTPacket(b)
	if err != nil {
		return nil, err
	}

	rc, err := p.DecodeByte()
	if err != nil {
		return nil, err
	}

	return &ConnackPacket{Packet: *p, ReturnCode: rc}, nil

}

func NewConnectPacket(p *Packet) (*ConnectPacket, error) {
	cp := &ConnectPacket{
		Packet: *p,
	}

	err := cp.DecodeProtocolVersion()
	if err != nil {
		return nil, err
	}
	err = cp.DecodeProtocolName()
	if err != nil {
		return nil, err
	}
	err = cp.DecodeConnectFlags()
	if err != nil {
		return nil, err
	}
	err = cp.DecodeKeepAlive()
	if err != nil {
		return nil, err
	}
	err = cp.DecodePayload()
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func (cp ConnectPacket) DecodeProtocolName() error {
	p := cp.DecodeString()
	cp.ProtocolName = p
	return nil
}

func (cp ConnectPacket) EncodeProtocolName() error {
	return cp.EncodeString(cp.ProtocolName)
}

func (cp ConnectPacket) DecodeProtocolVersion() error {
	v, err := cp.DecodeByte()
	if err != nil {
		return err
	}
	cp.ProtocolVersion = v
	return nil
}

func (cp ConnectPacket) EncodeProtocolVersion() error {
	return cp.EncodeByte(cp.ProtocolVersion)
}

func (cp ConnectPacket) DecodeConnectFlags() error {
	fb, err := cp.DecodeByte()
	if err != nil {
		return err
	}
	cp.UsernameFlag = fb>>6 > 0
	cp.PasswordFlag = (fb&0x40)>>5 > 0
	cp.WillRetainFlag = (fb&0x20)>>4 > 0
	cp.WillQoSFlag = (fb & 0x18) >> 3
	cp.WillFlag = fb>>6 > 0
	cp.CleanStartFlag = fb>>7 > 0
	return nil
}

func (cp ConnectPacket) EncodeConnectFlags(b []byte) ([]byte, error) {
	var flags byte
	if cp.UsernameFlag {
		flags = flags | (uint8(1) << 7)
	}
	if cp.PasswordFlag {
		flags = flags | (uint8(1) << 6)
	}
	if cp.WillRetainFlag {
		flags = flags | (uint8(1) << 5)
	}
	if cp.WillQoSFlag != 0 {
		flags = flags | (uint8(3) << 4)
	}
	if cp.WillFlag {
		flags = flags | (uint8(1) << 2)
	}
	if cp.CleanStartFlag {
		flags = flags | (uint8(1) << 1)
	}

	return append(b, flags), nil
}

func (cp ConnectPacket) DecodeKeepAlive() error {
	ka := cp.DecodeTwoByteInt()
	cp.KeepAlive = ka
	return nil
}

func (cp ConnectPacket) EncodeKeepAlive() error {
	return cp.EncodeTwoByteInt(cp.KeepAlive)
}

func (cp ConnectPacket) DecodePayload() error {

	err := cp.DecodeClientID()
	if err != nil {
		return err
	}
	// If willflag is set to 1, will topic is the next in the payload.
	if cp.WillFlag {
		err = cp.DecodeWillTopic()
		if err != nil {
			return err
		}
		err = cp.DecodeWillMessage()
		if err != nil {
			return err
		}
	}

	// If Username is set, username and password are next in the payload.
	if cp.UsernameFlag {
		err = cp.DecodeUsername()
		if err != nil {
			return err
		}
		err = cp.DecodePassword()
		if err != nil {
			return err
		}
	}

	return nil
}

func (cp ConnectPacket) DecodeWillTopic() error {
	cp.WillTopic = cp.DecodeString()
	return nil
}

func (cp ConnectPacket) DecodeWillMessage() error {
	//cp.WillPayload = cp.DecodeBinaryData() // TODO Implement this
	return nil
}

func (cp ConnectPacket) DecodeUsername() error {
	cp.Username = cp.DecodeString()
	return nil
}

func (cp ConnectPacket) DecodePassword() error {
	//cp.Password = cp.DecodeBinaryData() // TODO Implement this
	return nil
}

func (cp ConnectPacket) DecodeClientID() error {
	cp.ClientID = cp.DecodeString()
	return nil
}

func (cp ConnectPacket) EncodeClientID() error {
	return cp.EncodeString(cp.ClientID)
}

func (cp ConnectPacket) Encode() ([]byte, error) {

	// Starting from the variable header, fixed header is last.

	err := cp.EncodeString("MQTT")
	if err != nil {
		return nil, err
	}

	// Protocol Level revision level used by this client, we are using revision 4
	err = cp.EncodeByte(uint8(4))
	if err != nil {
		return nil, err
	}

	// Keepalive of the packet
	err = cp.EncodeTwoByteInt(cp.KeepAlive)
	if err != nil {
		return nil, err
	}

	// Encoding the Client Identifier
	err = cp.EncodeString(cp.ClientID)
	if err != nil {
		return nil, err
	}

	return cp.EncodeFixedHeader()
}

func (cp ConnackPacket) Encode() ([]byte, error) {
	if err := cp.EncodeByte(cp.SessionPresent); err != nil {
		return nil, err
	}

	if err := cp.EncodeByte(cp.ReturnCode); err != nil {
		return nil, err
	}

	//Connack is just the fixed header and the return code.
	return cp.EncodeFixedHeader()
}

func Accepted() ConnackPacket {
	p := &Packet{
		RemaningLength: 2,
		Type:           CONNACK,
		buff:           &bytes.Buffer{},
	}

	return ConnackPacket{
		Packet:         *p,
		SessionPresent: 1,
		ReturnCode:     ConnectionAccepted,
	}
}
func BadProtocolVersion() ConnackPacket {
	p := &Packet{
		RemaningLength: 2,
		Type:           CONNACK,
		buff:           &bytes.Buffer{},
	}
	return ConnackPacket{Packet: *p,
		SessionPresent: 0,
		ReturnCode:     UnnaceptableProtocolVersion,
	}
}

func InvalidIdentifier() ConnackPacket {
	p := &Packet{
		RemaningLength: 2,
		Type:           CONNACK,
		buff:           &bytes.Buffer{},
	}
	return ConnackPacket{Packet: *p,
		SessionPresent: 0,
		ReturnCode:     IdentifierRejected,
	}
}
func ServiceUnavailable() ConnackPacket {
	p := &Packet{
		RemaningLength: 2,
		Type:           CONNACK,
		buff:           &bytes.Buffer{},
	}
	return ConnackPacket{Packet: *p,
		SessionPresent: 0,
		ReturnCode:     ServerUnavailable,
	}
}
func BadAuth() ConnackPacket {
	p := &Packet{
		RemaningLength: 2,
		Type:           CONNACK,
		buff:           &bytes.Buffer{},
	}
	return ConnackPacket{Packet: *p,
		SessionPresent: 0,
		ReturnCode:     BadUsernameOrPassword,
	}
}
func NotAuth() ConnackPacket {
	p := &Packet{
		RemaningLength: 2,
		Type:           CONNACK,
		buff:           &bytes.Buffer{},
	}
	return ConnackPacket{Packet: *p,
		SessionPresent: 0,
		ReturnCode:     NotAuthorised,
	}
}
