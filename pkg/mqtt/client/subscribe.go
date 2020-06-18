package client

import (
	"io"

	"github.com/naspinall/Hive/pkg/mqtt/packets"
)

// Will need to refine.
func SubAck(pi uint16, rc byte, w io.Writer) error {
	sa := packets.SubAck(pi, rc)
	b, err := sa.Encode()
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}
