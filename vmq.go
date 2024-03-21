package vmq

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"

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
	// packet id
	muPID    sync.Mutex
	packetId uint16
	getPID   func() uint16
	// packets with conf info
	packetConn    *p.CONNECT_PACKET
	packetConnAck *p.CONNACK_PACKET
	packetDisconn *p.DISCONNECT_PACKET
	// callbacks
	onConnect    func(*vmq)
	onDisConnect func(*vmq)
	onError      func(*vmq)
	onMessage    func(*vmq, string, []byte)
}

func New(name string, network string, addr string, port int) *vmq {
	v := &vmq{
		name:    name,
		network: network,
		addr:    addr,
		port:    port,
		status:  STATUS_IDLE,
		conn:    nil,
		// callbacks
		onConnect:    func(*vmq) {},
		onDisConnect: func(*vmq) {},
		onError:      func(*vmq) {},
		onMessage:    func(*vmq, string, []byte) {},
	}
	v.getPID = func() uint16 {
		v.muPID.Lock()
		defer v.muPID.Unlock()

		v.packetId += 1
		return v.packetId
	}
	return v
}

// actions
func (v *vmq) Connect(cc *p.ConnectConf) error {
	if v.status != STATUS_IDLE {
		return errors.New("vmq status error: has connected")
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
	go handleData(v)

	// send connect
	v.packetConn = p.NewConnectPacketFromConf(cc)
	_, err = v.packetConn.ToStream(v.conn)
	if err != nil {
		v.conn.Close()
		return err
	}

	return nil
}

func (v *vmq) Disconnect(dc *p.DisconnectConf) error {
	if v.status != STATUS_CONNECTED {
		return errors.New("vmq status error: not connected")
	}

	v.packetDisconn = p.NewDisconnectPacketFromConf(dc)
	_, err := v.packetDisconn.ToStream(v.conn)
	if err != nil {
		v.conn.Close()
		return err
	}
	err = v.setStatus(STATUS_DISCONNECTED)
	if err != nil {
		v.conn.Close()
		return err
	}

	return nil
}

func (v *vmq) Subscribe(sc *p.SubscribeConf) error {
	if v.status != STATUS_CONNECTED {
		return errors.New("vmq status error: not connected")
	}

	// sub packet and id
	sub := p.NewSubscribePacketFromConf(sc)
	sub.VariableHeader.PacketId.FromValue(v.getPID())

	// just send
	_, err := sub.ToStream(v.conn)
	if err != nil {
		return err
	}
	return nil
}

func (v *vmq) Publish(pc *p.PublishConf) error {
	if v.status != STATUS_CONNECTED {
		return errors.New("vmq status error: not connected")
	}

	publish := p.NewPublishPacketFromConf(pc)
	// packetId if qos > 0
	if publish.Qos() > p.PUBLISH_FLAG_QOS_0 {
		publish.VariableHeader.PacketId.FromValue(v.getPID())
	}
	_, err := publish.ToStream(v.conn)
	if err != nil {
		return err
	}
	return nil
}

// set callbacks
func (v *vmq) OnConnect(f func(v *vmq)) {
	v.onConnect = f
}

func (v *vmq) OnDisConnect(f func(v *vmq)) {
	v.onDisConnect = f
}

func (v *vmq) OnError(f func(v *vmq)) {
	v.onError = f
}

func (v *vmq) OnMessage(f func(v *vmq, topic string, message []byte)) {
	v.onMessage = f
}

// handle data
func handleData(v *vmq) {
	packetType := t.NewByte()
	remainingLength := t.NewVarInt()
	for {
		// read fix header: PacketType + RemainingLength
		_, err := packetType.FromStream(v.conn)
		// EOF means the server is disconnected
		if err == io.EOF {
			v.conn.Close()
			v.setStatus(STATUS_DISCONNECTED)
			break
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
		err = handlePacket(v, packetType, remainingLength)
		if err != nil {
			v.conn.Close()
			v.setStatus(STATUS_IDLE)
			break
		}
	}
}

func handlePacket(v *vmq, packetType *t.MQTT_BYTE, remainingLength *t.MQTT_VAR_INT) error {
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
		return handleConnAck(v, ca)
	}
	if packetType.ToValue() == byte(p.PACKET_TYPE_DISCONNECT) {
		dc := p.NewDisconnectPacket(packetType, remainingLength)
		_, err := dc.FromStream(v.conn)
		if err != nil {
			return err
		}
		// handle connACK
		return handleDisconn(v, dc)
	}
	if packetType.ToValue() == byte(p.PACKET_TYPE_SUBACK) {
		sa := p.NewSubackPacket(packetType, remainingLength)
		_, err := sa.FromStream(v.conn)
		if err != nil {
			return err
		}
		return handleSuback(v, sa)
	}
	if packetType.ToValue() == byte(p.PACKET_TYPE_PUBLISH) {
		pub := p.NewPublishPacket(packetType, remainingLength)
		_, err := pub.FromStream(v.conn)
		if err != nil {
			return err
		}
		return handlePublish(v, pub)
	}
	return errors.New("can't metch packet type")
}

func handleConnAck(v *vmq, ca *p.CONNACK_PACKET) error {
	if v.status != STATUS_CONNECTING {
		return errors.New("vmq status error: not connecting")
	}
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

func handleDisconn(v *vmq, dc *p.DISCONNECT_PACKET) error {
	if v.status != STATUS_CONNECTED {
		return errors.New("vmq status error: not connected")
	}
	v.packetDisconn = dc
	err := v.conn.Close()
	if err != nil {
		return err
	}
	err = v.setStatus(STATUS_DISCONNECTED)
	if err != nil {
		return err
	}
	return nil
}

func handleSuback(v *vmq, sa *p.SUBACK_PACKET) error {
	// TODO
	fmt.Println(sa)
	return nil
}

func handlePublish(v *vmq, pub *p.PUBLISH_PACKET) error {
	v.onMessage(v, pub.VariableHeader.TopicName.ToValue(), pub.Payload.ApplicationMessage)
	return nil
}
