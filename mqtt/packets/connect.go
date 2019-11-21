package packets

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
	FixedHeader *FixedHeader

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
	Password       string
	WillTopic      string
	WillPayload    []byte
}

type ConnackPacket struct {
	FixedHeader *FixedHeader
	ReturnCode  byte
}

func NewConnackPacket(b []byte) (*ConnackPacket, error) {
	fh, err := NewFixedHeader(b)
	if err != nil {
		return nil, err
	}
	rc := b[3]
	return &ConnackPacket{FixedHeader: fh, ReturnCode: rc}, nil

}

func NewConnectPacket(fh *FixedHeader, b []byte) (*ConnectPacket, error) {
	cp := &ConnectPacket{FixedHeader: fh}
	n, err := cp.DecodeProtocolName(b, 2)
	if err != nil {
		return nil, err
	}
	n, err = cp.DecodeProtocolVersion(b, n)
	if err != nil {
		return nil, err
	}
	n, err = cp.DecodeConnectFlags(b, n)
	if err != nil {
		return nil, err
	}
	n, err = cp.DecodeKeepAlive(b, n)
	if err != nil {
		return nil, err
	}
	err = cp.DecodePayload(b, n)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func (cp *ConnectPacket) DecodeProtocolName(b []byte, n int) (int, error) {
	p, m, err := DecodeString(b[n:])
	cp.ProtocolName = p
	return n + m, err
}

func (cp *ConnectPacket) EncodeProtocolName(b []byte) ([]byte, error) {
	return EncodeString(b, cp.ProtocolName)
}

func (cp *ConnectPacket) DecodeProtocolVersion(b []byte, n int) (int, error) {
	v, m, err := DecodeByte(b[n:])
	cp.ProtocolVersion = v
	return n + m, err
}

func (cp *ConnectPacket) EncodeProtocolVersion(b []byte) ([]byte, error) {
	return EncodeByte(b, cp.ProtocolVersion)
}

func (cp *ConnectPacket) DecodeConnectFlags(b []byte, n int) (int, error) {
	fb := b[n]
	cp.UsernameFlag = fb>>6 > 0
	cp.PasswordFlag = (fb&0x40)>>5 > 0
	cp.WillRetainFlag = (fb&0x20)>>4 > 0
	cp.WillQoSFlag = (fb & 0x18) >> 3
	cp.WillFlag = fb>>6 > 0
	cp.CleanStartFlag = fb>>7 > 0
	return n + 1, nil
}

func (cp *ConnectPacket) EncodeConnectFlags(b []byte) ([]byte, error) {
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

func (cp *ConnectPacket) DecodeKeepAlive(b []byte, n int) (int, error) {
	ka, m, err := DecodeTwoByteInt(b[n:])
	if err != nil {
		return -1, err
	}
	cp.KeepAlive = ka
	return n + m, nil
}

func (cp *ConnectPacket) EncodeKeepAlive(b []byte) ([]byte, error) {
	return EncodeTwoByteInt(b, cp.KeepAlive)
}

func (cp *ConnectPacket) DecodePayload(b []byte, n int) error {

	err := cp.DecodeClientID(b, n)
	if err != nil {
		return err
	}
	return nil
}

func (cp *ConnectPacket) DecodeClientID(b []byte, n int) error {
	clientId, n, err := DecodeString(b[n:])
	if err != nil {
		return err
	}
	cp.ClientID = clientId
	return nil
}

func (cp *ConnectPacket) EncodeClientID(b []byte) ([]byte, error) {
	return EncodeString(b, cp.ClientID)
}

func (cp *ConnectPacket) EncodeConnectPacket() ([]byte, error) {
	var b []byte

	// Starting from the variable header, fixed header is last.

	b, err := EncodeString(b, "MQTT")
	if err != nil {
		return nil, err
	}

	// Protocol Level revision level used by this client, we are using revision 4
	b, err = EncodeByte(b, uint8(4))
	if err != nil {
		return nil, err
	}

	// Keepalive of the packet
	b, err = EncodeTwoByteInt(b, cp.KeepAlive)
	if err != nil {
		return nil, err
	}

	// Encoding the Client Identifier
	b, err = EncodeString(b, cp.ClientID)
	if err != nil {
		return nil, err
	}

	// TODO Implement all the will bullshit.

	//Encoding the FixedHeader
	cp.FixedHeader.RemaningLength = len(b)
	fhb, err := cp.FixedHeader.EncodeFixedHeader()
	if err != nil {
		return nil, err
	}

	return append(fhb, b...), nil
}

func (cp *ConnackPacket) EncodeConnackPacket() ([]byte, error) {
	//Connack is just the fixed header and the return code.
	b, err := cp.FixedHeader.EncodeFixedHeader()
	return append(b, cp.ReturnCode), err
}
