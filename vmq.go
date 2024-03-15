package vmq

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	p "github.com/vitsumoc/vmq/packets"
	t "github.com/vitsumoc/vmq/types"
)

type vmq struct {
	name    string
	network string
	addr    string
	port    int
	status  VMQ_STATUS
	conn    net.Conn
	// packets with conf info
	packetConn    *p.CONNECT_PACKET
	packetConnAck *p.CONNACK_PACKET
}

func New(name string, network string, addr string, port int) *vmq {
	return &vmq{
		name:    name,
		network: network,
		addr:    addr,
		port:    port,
		status:  STATUS_IDLE,
		conn:    nil,
	}
}

func (v *vmq) Connect(cc *p.ConnectConf) error {
	if v.status != STATUS_IDLE {
		return errors.New("vmq has connected")
	}
	// dial network
	var err error
	v.conn, err = net.DialTimeout(v.network, v.addr+":"+strconv.Itoa(v.port), DIAL_TIMEOUT)
	if err != nil {
		return err
	}
	err = v.setStatus(STATUS_CONNECTING)
	if err != nil {
		v.conn.Close()
		return err
	}

	// listen from server
	go onData(v)

	// send connect
	v.packetConn = p.NewConnectPacket(cc)
	_, err = v.packetConn.ToStream(v.conn)
	if err != nil {
		v.conn.Close()
		return err
	}

	return nil
}

// handle server data
func onData(v *vmq) {
	packetType := t.NewByte()
	remainingLength := t.NewVarInt()
	for {
		time.Sleep(LISTEN_INTERVAL)
		// read fix header: PacketType + RemainingLength
		_, err := packetType.FromStream(v.conn)
		if err == io.EOF {
			continue
		} else if err != nil {
			v.conn.Close()
			v.setStatus(STATUS_IDLE)
			break
		}
		_, err = remainingLength.FromStream(v.conn)
		if err != nil {
			v.conn.Close()
			v.setStatus(STATUS_IDLE)
			break
		}
		// handle packet
		err = recPacket(v, packetType, remainingLength)
		if err != nil {
			v.conn.Close()
			v.setStatus(STATUS_IDLE)
			break
		}
	}
}

func recPacket(v *vmq, packetType *t.MQTT_BYTE, remainingLength *t.MQTT_VAR_INT) error {
	// slect packet by type
	// TODO need more packet type
	if packetType.ToValue() == byte(p.PACKET_TYPE_CONNACK) {
		// stream to packet
		ca := p.NewConnackPacket(packetType, remainingLength)
		_, err := ca.FromStream(v.conn)
		if err != nil {
			return err
		}
		// handle connACK
		return onConnAck(v, ca)
	}
	return errors.New("can't metch packet type")
}

func onConnAck(v *vmq, ca *p.CONNACK_PACKET) error {
	// error check
	if ca.VariableHeader.ConnectReasonCode.ToValue() >= byte(p.RC_UNSPECIFIED_ERROR) {
		return fmt.Errorf("connack error, rc is:%v", ca.VariableHeader.ConnectReasonCode.ToValue())
	}
	// save packet and turn to CONNTED
	v.packetConnAck = ca
	err := v.setStatus(STATUS_CONNECTED)
	if err != nil {
		return err
	}
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
