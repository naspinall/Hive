package packets

type DisconnectPacket struct {
	FixedHeader *FixedHeader
}

func NewDisconnectPacket(fh *FixedHeader, b []byte) (*DisconnectPacket, error) {
	dp := &DisconnectPacket{FixedHeader: fh}
	return dp, nil
}
