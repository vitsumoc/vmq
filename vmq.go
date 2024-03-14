package vmq

import (
	"net"

	"github.com/vitsumoc/vmq/packets"
)

type vmq struct {
	name      string
	addr      net.Addr
	conn      net.Conn
	onPublish func(*packets.PUBLISH_PACKET)
}

func New(name string, addr net.Addr) *vmq {
	return &vmq{
		name:      name,
		addr:      addr,
		onPublish: func(p *packets.PUBLISH_PACKET) {},
	}
}

func (*vmq) Connect(cc *packets.ConnectConf) error {
	return nil
}

// func (*vmq) Disconnect(*packets.DISCONNECT_PACKET) {

// }

// func (*vmq) Publish(*packets.PUBLISH_PACKET) {

// }

// func (*vmq) Subscribe(*packets.SUBSCRIBE_PACKET) {

// }

// func (*vmq) Unsubscribe(*packets.UNSUBSCRIBE_PACKET) {

// }

// func (*vmq) Pingreq(*packets.PINGREQ_PACKET) {

// }
