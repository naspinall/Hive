package packets

type WillProperties struct {
}

type ConnectPacket struct {
	FixedHeader *FixedHeader

	//Variable Header
	ProtocolName    string
	ProtocolVersion byte
	KeepAlive       uint32

	// Connect Flags
	UsernameFlag   bool
	PasswordFlag   bool
	WillRetainFlag bool
	WillQoSFlag    bool
	WillFlag       bool
	CleanStartFlag bool

	//Variable Header Properties
	SessionExpiryInterval     uint32
	AuthMethod                []byte
	RequestProblemInformation byte
	RecieveMaximum            uint16
	TopicAliasMaximum         uint16
	UserProperty              StringPair
	MaximumPacketSize         uint32

	//Payload properties
	ClientID       string
	Username       string
	WillProperties WillProperties
	Password       string
}
