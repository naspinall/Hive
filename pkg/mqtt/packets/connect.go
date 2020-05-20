package packets

type ConnectPacket struct {
	FixedHeader *FixedHeader

	//Variable Header
	ProtocolName    string
	ProtocolVersion byte
	KeepAlive       uint16

	// Connect Flags
	UsernameFlag     bool
	PasswordFlag     bool
	WillRetainFlag   bool
	WillQoSFlag      uint8
	WillFlag         bool
	CleanSessionFlag bool

	//Payload properties
	ClientID    string
	Username    string
	Password    []byte
	WillTopic   string
	WillPayload []byte
}

type ConnackPacket struct {
	FixedHeader    *FixedHeader
	SessionPresent byte
	ReturnCode     byte
}

func NewConnackPacket(fh *FixedHeader, b []byte) (*ConnackPacket, error) {
	rc := b[0]
	return &ConnackPacket{FixedHeader: fh, ReturnCode: rc}, nil

}

func NewConnectPacket(fh *FixedHeader, b []byte) (*ConnectPacket, error) {
	cp := &ConnectPacket{FixedHeader: fh}
	n, err := cp.DecodeProtocolVersion(b, 0)
	if err != nil {
		return nil, err
	}
	n, err = cp.DecodeProtocolName(b, n)
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

func (cp ConnectPacket) DecodeProtocolName(b []byte, n int) (int, error) {
	p, m, err := DecodeString(b[n:])
	cp.ProtocolName = p
	return n + m, err
}

func (cp ConnectPacket) EncodeProtocolName(b []byte) ([]byte, error) {
	return EncodeString(b, cp.ProtocolName)
}

func (cp ConnectPacket) DecodeProtocolVersion(b []byte, n int) (int, error) {
	v, m, err := DecodeByte(b[n:])
	cp.ProtocolVersion = v
	return n + m, err
}

func (cp ConnectPacket) EncodeProtocolVersion(b []byte) ([]byte, error) {
	return EncodeByte(b, cp.ProtocolVersion)
}

func (cp ConnectPacket) DecodeConnectFlags(b []byte, n int) (int, error) {
	fb := b[n]
	cp.UsernameFlag = fb>>6 > 0
	cp.PasswordFlag = (fb&0x40)>>5 > 0
	cp.WillRetainFlag = (fb&0x20)>>4 > 0
	cp.WillQoSFlag = (fb & 0x18) >> 3
	cp.WillFlag = fb>>6 > 0
	cp.CleanSessionFlag = fb>>7 > 0
	return n + 1, nil
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
	if cp.CleanSessionFlag {
		flags = flags | (uint8(1) << 1)
	}

	return append(b, flags), nil
}

func (cp ConnectPacket) DecodeKeepAlive(b []byte, n int) (int, error) {
	ka, m, err := DecodeTwoByteInt(b[n:])
	if err != nil {
		return -1, err
	}
	cp.KeepAlive = ka
	return n + m, nil
}

func (cp ConnectPacket) EncodeKeepAlive(b []byte) ([]byte, error) {
	return EncodeTwoByteInt(b, cp.KeepAlive)
}

func (cp ConnectPacket) DecodePayload(b []byte, n int) error {

	n, err := cp.DecodeClientID(b, n)
	if err != nil {
		return err
	}
	// If willflag is set to 1, fill topic is the next in the payload.
	if cp.WillFlag {
		n, err = cp.DecodeWillTopic(b, n)
		if err != nil {
			return err
		}
		n, err = cp.DecodeWillMessage(b, n)
		if err != nil {
			return err
		}
	}

	// If Username is set, username and password are next in the payload.
	if cp.UsernameFlag {
		n, err = cp.DecodeUsername(b, n)
		if err != nil {
			return err
		}
		n, err = cp.DecodePassword(b, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cp ConnectPacket) DecodeWillTopic(b []byte, n int) (int, error) {
	topic, m, err := DecodeString(b[n:])
	if err != nil {
		return -1, err
	}
	cp.WillTopic = topic
	return m + n, nil

}

func (cp ConnectPacket) DecodeWillMessage(b []byte, n int) (int, error) {

	messageLength, m, err := DecodeTwoByteInt(b[n:])
	n = n + m // Moving past the will message length.

	cp.WillPayload, m, err = DecodeBinaryData(b[n : n+int(messageLength)])
	if err != nil {
		return -1, err
	}
	return n + m, nil
}

func (cp ConnectPacket) DecodeUsername(b []byte, n int) (int, error) {
	username, m, err := DecodeString(b[n:])
	if err != nil {
		return -1, err
	}
	cp.Username = username
	return m + n, nil
}

func (cp ConnectPacket) DecodePassword(b []byte, n int) (int, error) {
	messageLength, m, err := DecodeTwoByteInt(b[n:])
	n = n + m // Moving past the password length.

	cp.Password, m, err = DecodeBinaryData(b[n : n+int(messageLength)])
	if err != nil {
		return -1, err
	}
	return n + m, nil
}

func (cp ConnectPacket) DecodeClientID(b []byte, n int) (int, error) {
	clientId, n, err := DecodeString(b[n:])
	if err != nil {
		return -1, err
	}
	cp.ClientID = clientId
	return n, nil
}

func (cp ConnectPacket) EncodeClientID(b []byte) ([]byte, error) {
	return EncodeString(b, cp.ClientID)
}

func (cp ConnectPacket) Encode() ([]byte, error) {
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

func (cp ConnackPacket) Encode() ([]byte, error) {
	//Connack is just the fixed header and the return code.
	b, err := cp.FixedHeader.EncodeFixedHeader()
	b = append(b, cp.SessionPresent)
	return append(b, cp.ReturnCode), err
}

func Accepted() ConnackPacket {
	fh := &FixedHeader{
		RemaningLength: 2,
		Type:           CONNACK,
	}
	return ConnackPacket{
		FixedHeader:    fh,
		SessionPresent: 1,
		ReturnCode:     ConnectionAccepted}
}
func BadProtocolVersion() ConnackPacket {
	fh := &FixedHeader{
		RemaningLength: 2,
		Type:           CONNACK,
	}
	return ConnackPacket{FixedHeader: fh,
		SessionPresent: 0,
		ReturnCode:     UnnaceptableProtocolVersion}
}

func InvalidIdentifier() ConnackPacket {
	fh := &FixedHeader{
		RemaningLength: 2,
		Type:           CONNACK,
	}
	return ConnackPacket{FixedHeader: fh,
		SessionPresent: 0,
		ReturnCode:     IdentifierRejected}
}
func ServiceUnavailable() ConnackPacket {
	fh := &FixedHeader{
		RemaningLength: 2,
		Type:           CONNACK,
	}
	return ConnackPacket{FixedHeader: fh,
		SessionPresent: 0,
		ReturnCode:     ServerUnavailable}
}
func BadAuth() ConnackPacket {
	fh := &FixedHeader{
		RemaningLength: 2,
		Type:           CONNACK,
	}
	return ConnackPacket{FixedHeader: fh,
		SessionPresent: 0,
		ReturnCode:     BadUsernameOrPassword}
}
func NotAuth() ConnackPacket {
	fh := &FixedHeader{
		RemaningLength: 2,
		Type:           CONNACK,
	}
	return ConnackPacket{FixedHeader: fh,
		SessionPresent: 0,
		ReturnCode:     NotAuthorised}
}

func Connect() (cp ConnectPacket) {
	cp.FixedHeader = &FixedHeader{
		Type: 1,
	}
	cp.ProtocolName = "MQTT"
	cp.ProtocolVersion = 0x4
	cp.KeepAlive = 10
	return
}
