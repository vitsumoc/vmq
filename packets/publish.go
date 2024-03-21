package packets

import (
	"bytes"
	"io"

	t "github.com/vitsumoc/vmq/types"
)

type PUBLISH_FLAG_QOS byte

const PUBLISH_FLAG_RETAIN byte = 0x01
const PUBLISH_FLAG_QOS_0 PUBLISH_FLAG_QOS = 0x00
const PUBLISH_FLAG_QOS_1 PUBLISH_FLAG_QOS = 0x02
const PUBLISH_FLAG_QOS_2 PUBLISH_FLAG_QOS = 0x04
const PUBLISH_FLAG_DUP byte = 0x09

type PUBLISH_PACKET struct {
	FixHeader      PUBLISH_FIX_HEADER
	VariableHeader PUBLISH_VARIABLE_HEADER
	Payload        PUBLISH_PAYLOAD
}

type PUBLISH_FIX_HEADER struct {
	PacketType      t.MQTT_BYTE // pak type dup qos retain
	RemainingLength t.MQTT_VAR_INT
}

type PUBLISH_VARIABLE_HEADER struct {
	TopicName         t.MQTT_UTF8
	PacketId          t.MQTT_U16
	PublishProperties PROPERTIES
}

type PUBLISH_PAYLOAD struct {
	ApplicationMessage []byte
}

func NewPublishPacketFromConf(pc *PublishConf) *PUBLISH_PACKET {
	packet := PUBLISH_PACKET{}
	// fix header
	packetType := byte(PACKET_TYPE_PUBLISH)
	if pc.retain {
		packetType |= PUBLISH_FLAG_RETAIN
	}
	packetType |= byte(pc.qos)
	if pc.dup {
		packetType |= PUBLISH_FLAG_DUP
	}
	packet.FixHeader.PacketType.FromValue(packetType)
	remainingLen := 0

	// variable header
	packet.VariableHeader.TopicName.FromValue(pc.topic)
	remainingLen += packet.VariableHeader.TopicName.Length()
	// only qos1 or qos2, there is packetId
	if packet.Qos() > PUBLISH_FLAG_QOS_0 {
		// PID is 0, then vmq will set it
		packet.VariableHeader.PacketId.FromValue(0)
		remainingLen += packet.VariableHeader.PacketId.Length()
	}
	packet.VariableHeader.PublishProperties = pc.publishProperties
	remainingLen += packet.VariableHeader.PublishProperties.Length()

	// payload
	packet.Payload.ApplicationMessage = pc.applicationMessage
	remainingLen += len(packet.Payload.ApplicationMessage)

	// set remainingLen
	packet.FixHeader.RemainingLength.FromValue(remainingLen)

	return &packet
}

func (p *PUBLISH_PACKET) ToStream(output io.Writer) (int, error) {
	var err error
	var buffer = bytes.NewBuffer(nil)
	// fix header
	_, err = p.FixHeader.PacketType.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	_, err = p.FixHeader.RemainingLength.ToStream(buffer)
	if err != nil {
		return 0, err
	}

	// variable header
	_, err = p.VariableHeader.TopicName.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	if p.Qos() > PUBLISH_FLAG_QOS_0 {
		_, err = p.VariableHeader.PacketId.ToStream(buffer)
		if err != nil {
			return 0, err
		}
	}
	_, err = p.VariableHeader.PublishProperties.ToStream(buffer)
	if err != nil {
		return 0, err
	}

	// payload
	_, err = buffer.Write(p.Payload.ApplicationMessage)
	if err != nil {
		return 0, err
	}

	// to stream
	return output.Write(buffer.Bytes())
}

func NewPublishPacket(packetType *t.MQTT_BYTE, remainingLength *t.MQTT_VAR_INT) *PUBLISH_PACKET {
	return &PUBLISH_PACKET{
		FixHeader: PUBLISH_FIX_HEADER{
			PacketType:      *packetType,
			RemainingLength: *remainingLength,
		},
		VariableHeader: PUBLISH_VARIABLE_HEADER{
			PublishProperties: *NewProperties(),
		},
		Payload: PUBLISH_PAYLOAD{
			ApplicationMessage: make([]byte, 0),
		},
	}
}

func (p *PUBLISH_PACKET) FromStream(input io.Reader) (int, error) {
	length := 0
	var n int
	var err error

	// variable header
	n, err = p.VariableHeader.TopicName.FromStream(input)
	if err != nil {
		return 0, err
	}
	length += n
	if p.Qos() > PUBLISH_FLAG_QOS_0 {
		n, err = p.VariableHeader.PacketId.FromStream(input)
		if err != nil {
			return 0, err
		}
		length += n
	}
	n, err = p.VariableHeader.PublishProperties.FromStream(input)
	if err != nil {
		return 0, err
	}
	length += n

	p.Payload.ApplicationMessage = make([]byte, p.FixHeader.RemainingLength.ToValue()-length)
	n, err = io.ReadAtLeast(input, p.Payload.ApplicationMessage, p.FixHeader.RemainingLength.ToValue()-length)
	if err != nil {
		return 0, err
	}
	length += n

	return length, nil
}

func (p *PUBLISH_PACKET) Qos() PUBLISH_FLAG_QOS {
	if p.FixHeader.PacketType.ToValue()&byte(PUBLISH_FLAG_QOS_2) > 0 {
		return PUBLISH_FLAG_QOS_2
	}
	if p.FixHeader.PacketType.ToValue()&byte(PUBLISH_FLAG_QOS_1) > 0 {
		return PUBLISH_FLAG_QOS_1
	}
	return PUBLISH_FLAG_QOS_0
}
