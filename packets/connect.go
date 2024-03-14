package packets

import (
	"io"

	t "github.com/vitsumoc/vmq/types"
)

const CONNECT_FLAG_CLEANSTART byte = 0x02
const CONNECT_FLAG_WILLFLAG byte = 0x04
const CONNECT_FLAG_WILLQOS_0 byte = 0x00
const CONNECT_FLAG_WILLQOS_1 byte = 0x08
const CONNECT_FLAG_WILLQOS_2 byte = 0x10
const CONNECT_FLAG_WILLRETAIN byte = 0x20
const CONNECT_FLAG_PASSWORD byte = 0x40
const CONNECT_FLAG_USERNAME byte = 0x80

type CONNECT_PACKET struct {
	FixHeader      CONNECT_FIX_HEADER
	VariableHeader CONNECT_VARIABLE_HEADER
	Payload        CONNECT_PAYLOAD
}

type CONNECT_FIX_HEADER struct {
	PacketType      t.MQTT_BYTE
	RemainingLength t.MQTT_VAR_INT
}

type CONNECT_VARIABLE_HEADER struct {
	ProtocolName  t.MQTT_UTF8
	ProtocolLevel t.MQTT_BYTE
	ConnectFlags  t.MQTT_BYTE
	KeepAlive     t.MQTT_U16
	Properties    PROPERTIES
}

type CONNECT_PAYLOAD struct {
	ClientID       t.MQTT_UTF8
	WillProperties PROPERTIES
	WillTopic      t.MQTT_UTF8
	WillPayload    t.MQTT_BIN
	UserName       t.MQTT_UTF8
	Password       t.MQTT_BIN
}

func NewConnectPacket(cc *ConnectConf) *CONNECT_PACKET {
	packet := CONNECT_PACKET{}
	// fix header
	packet.FixHeader.PacketType.FromValue(byte(PACKET_TYPE_CONNECT))
	remainingLen := 0

	// variable header
	packet.VariableHeader.ProtocolName.FromValue(MQTT_PROTOCOL_NAME)
	remainingLen += packet.VariableHeader.ProtocolName.Length()
	packet.VariableHeader.ProtocolLevel.FromValue(MQTT_PROTOCOL_LEVEL)
	remainingLen += packet.VariableHeader.ProtocolLevel.Length()
	// connect flag
	var cf byte = 0x00
	if cc.cfCleanStart {
		cf = cf | CONNECT_FLAG_CLEANSTART
	}
	if cc.cfWillFlag {
		cf = cf | CONNECT_FLAG_WILLFLAG
	}
	if cc.cfWillQos == 0 {
		cf = cf | CONNECT_FLAG_WILLQOS_0
	}
	if cc.cfWillQos == 1 {
		cf = cf | CONNECT_FLAG_WILLQOS_1
	}
	if cc.cfWillQos == 2 {
		cf = cf | CONNECT_FLAG_WILLQOS_2
	}
	if cc.cfWillRetain {
		cf = cf | CONNECT_FLAG_WILLRETAIN
	}
	if cc.cfPassword {
		cf = cf | CONNECT_FLAG_PASSWORD
	}
	if cc.cfUsername {
		cf = cf | CONNECT_FLAG_USERNAME
	}
	packet.VariableHeader.ConnectFlags.FromValue(cf)
	remainingLen += packet.VariableHeader.ConnectFlags.Length()
	packet.VariableHeader.KeepAlive.FromValue(cc.keepAlive)
	remainingLen += packet.VariableHeader.KeepAlive.Length()
	packet.VariableHeader.Properties = cc.properties
	remainingLen += packet.VariableHeader.Properties.PropertyLength.Length()
	remainingLen += packet.VariableHeader.Properties.PropertyLength.ToValue()

	// payload
	packet.Payload.ClientID.FromValue(cc.clientID)
	remainingLen += packet.Payload.ClientID.Length()
	if cc.cfWillFlag {
		packet.Payload.WillProperties = cc.willProperties
		remainingLen += packet.Payload.WillProperties.PropertyLength.Length()
		remainingLen += packet.Payload.WillProperties.PropertyLength.ToValue()
		packet.Payload.WillTopic.FromValue(cc.willTopic)
		remainingLen += packet.Payload.WillTopic.Length()
		packet.Payload.WillPayload.FromValue(cc.willPayload)
		remainingLen += packet.Payload.WillPayload.Length()
	}
	if cc.cfUsername {
		packet.Payload.UserName.FromValue(cc.username)
		remainingLen += packet.Payload.UserName.Length()
	}
	if cc.cfPassword {
		packet.Payload.Password.FromValue(cc.password)
		remainingLen += packet.Payload.Password.Length()
	}

	// set remainingLen
	packet.FixHeader.RemainingLength.FromValue(remainingLen)

	return &packet
}

func (c *CONNECT_PACKET) ToStream(output io.Writer) (int, error) {
	var length, n int
	var err error
	// fix header
	n, err = c.FixHeader.PacketType.ToStream(output)
	if err != nil {
		return length, err
	}
	length += n
	n, err = c.FixHeader.RemainingLength.ToStream(output)
	if err != nil {
		return length, err
	}
	length += n

	// variable header
	n, err = c.VariableHeader.ProtocolName.ToStream(output)
	if err != nil {
		return length, err
	}
	length += n
	n, err = c.VariableHeader.ProtocolLevel.ToStream(output)
	if err != nil {
		return length, err
	}
	length += n
	n, err = c.VariableHeader.ConnectFlags.ToStream(output)
	if err != nil {
		return length, err
	}
	length += n
	n, err = c.VariableHeader.KeepAlive.ToStream(output)
	if err != nil {
		return length, err
	}
	length += n
	n, err = c.VariableHeader.Properties.ToStream(output)
	if err != nil {
		return length, err
	}
	length += n

	// payload
	n, err = c.Payload.ClientID.ToStream(output)
	if err != nil {
		return length, err
	}
	length += n
	// there is will
	if c.VariableHeader.ConnectFlags.ToValue()&CONNECT_FLAG_WILLFLAG > 0x00 {
		n, err = c.Payload.WillProperties.ToStream(output)
		if err != nil {
			return length, err
		}
		length += n
		n, err = c.Payload.WillTopic.ToStream(output)
		if err != nil {
			return length, err
		}
		length += n
		n, err = c.Payload.WillPayload.ToStream(output)
		if err != nil {
			return length, err
		}
		length += n
	}
	// there is username
	if c.VariableHeader.ConnectFlags.ToValue()&CONNECT_FLAG_USERNAME > 0x00 {
		n, err = c.Payload.UserName.ToStream(output)
		if err != nil {
			return length, err
		}
		length += n
	}
	// there is password
	if c.VariableHeader.ConnectFlags.ToValue()&CONNECT_FLAG_PASSWORD > 0x00 {
		n, err = c.Payload.Password.ToStream(output)
		if err != nil {
			return length, err
		}
		length += n
	}

	return length, nil
}
