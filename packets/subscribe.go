package packets

import (
	"bytes"
	"io"

	t "github.com/vitsumoc/vmq/types"
)

type SUBSCRIBE_OPTION_QOS byte
type SUBSCRIBE_OPTION_RETAIN_HANDLING byte

const SUBSCRIBE_OPTION_QOS_0 SUBSCRIBE_OPTION_QOS = 0x00
const SUBSCRIBE_OPTION_QOS_1 SUBSCRIBE_OPTION_QOS = 0x01
const SUBSCRIBE_OPTION_QOS_2 SUBSCRIBE_OPTION_QOS = 0x02
const SUBSCRIBE_OPTION_NO_LOCAL byte = 0x04
const SUBSCRIBE_OPTION_RETAIN_AS_PUBLISHED byte = 0x08
const SUBSCRIBE_OPTION_RETAIN_HANDLING_0 SUBSCRIBE_OPTION_RETAIN_HANDLING = 0x00
const SUBSCRIBE_OPTION_RETAIN_HANDLING_1 SUBSCRIBE_OPTION_RETAIN_HANDLING = 0x10
const SUBSCRIBE_OPTION_RETAIN_HANDLING_2 SUBSCRIBE_OPTION_RETAIN_HANDLING = 0x20

type SUBSCRIBE_PACKET struct {
	FixHeader      SUBSCRIBE_FIX_HEADER
	VariableHeader SUBSCRIBE_VARIABLE_HEADER
	Payload        SUBSCRIBE_PAYLOAD
}

type SUBSCRIBE_FIX_HEADER struct {
	PacketType      t.MQTT_BYTE
	RemainingLength t.MQTT_VAR_INT
}

type SUBSCRIBE_VARIABLE_HEADER struct {
	PacketId            t.MQTT_U16
	SubscribeProperties PROPERTIES
}

type SUBSCRIBE_PAYLOAD struct {
	SubPairs []SUB_PAIR
}

type SUB_PAIR struct {
	TopicFilter         t.MQTT_UTF8
	SubscriptionOptions t.MQTT_BYTE
}

func NewSubscribePacketFromConf(sc *SubscribeConf) *SUBSCRIBE_PACKET {
	packet := SUBSCRIBE_PACKET{}
	// fix header
	packet.FixHeader.PacketType.FromValue(byte(PACKET_TYPE_SUBSCRIBE))
	remainingLen := 0

	// variable header
	// packetId is allocated by vmq, so there is no need to configure
	remainingLen += packet.VariableHeader.PacketId.Length()
	packet.VariableHeader.SubscribeProperties = sc.subscribeProperties
	remainingLen += packet.VariableHeader.SubscribeProperties.Length()

	// payload
	packet.Payload.SubPairs = sc.subPairs
	for x := 0; x < len(packet.Payload.SubPairs); x++ {
		remainingLen += packet.Payload.SubPairs[x].TopicFilter.Length()
		remainingLen += packet.Payload.SubPairs[x].SubscriptionOptions.Length()
	}

	// set remainingLen
	packet.FixHeader.RemainingLength.FromValue(remainingLen)

	return &packet
}

func (s *SUBSCRIBE_PACKET) ToStream(output io.Writer) (int, error) {
	var err error
	var buffer = bytes.NewBuffer(nil)
	// fix header
	_, err = s.FixHeader.PacketType.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	_, err = s.FixHeader.RemainingLength.ToStream(buffer)
	if err != nil {
		return 0, err
	}

	// variable header
	_, err = s.VariableHeader.PacketId.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	_, err = s.VariableHeader.SubscribeProperties.ToStream(buffer)
	if err != nil {
		return 0, err
	}

	// payload
	for x := 0; x < len(s.Payload.SubPairs); x++ {
		_, err = s.Payload.SubPairs[x].TopicFilter.ToStream(buffer)
		if err != nil {
			return 0, err
		}
		_, err = s.Payload.SubPairs[x].SubscriptionOptions.ToStream(buffer)
		if err != nil {
			return 0, err
		}
	}
	// to stream
	return output.Write(buffer.Bytes())
}
