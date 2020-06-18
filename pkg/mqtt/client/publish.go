package client

import (
	"io"

	"github.com/naspinall/Hive/pkg/mqtt/packets"
)

// Will need to refine.
func Publish(pp *packets.PublishPacket, w io.Writer) error {
	b, err := pp.Encode()
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}
