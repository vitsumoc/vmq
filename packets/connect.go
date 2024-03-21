package packets

import (
	"bytes"
	"io"

	t "github.com/vitsumoc/vmq/types"
)

type CONNECT_FLAG_WILLQOS byte

const CONNECT_FLAG_CLEANSTART byte = 0x02
const CONNECT_FLAG_WILLFLAG byte = 0x04
const CONNECT_FLAG_WILLQOS_0 CONNECT_FLAG_WILLQOS = 0x00
const CONNECT_FLAG_WILLQOS_1 CONNECT_FLAG_WILLQOS = 0x08
const CONNECT_FLAG_WILLQOS_2 CONNECT_FLAG_WILLQOS = 0x10
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
	ProtocolName      t.MQTT_UTF8
	ProtocolLevel     t.MQTT_BYTE
	ConnectFlags      t.MQTT_BYTE
	KeepAlive         t.MQTT_U16
	ConnectProperties PROPERTIES
}

type CONNECT_PAYLOAD struct {
	ClientID       t.MQTT_UTF8
	WillProperties PROPERTIES
	WillTopic      t.MQTT_UTF8
	WillPayload    t.MQTT_BIN
	UserName       t.MQTT_UTF8
	Password       t.MQTT_BIN
}

func NewConnectPacketFromConf(cc *ConnectConf) *CONNECT_PACKET {
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
	cf = cf | byte(cc.cfWillQos)
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
	packet.VariableHeader.ConnectProperties = cc.connectProperties
	remainingLen += packet.VariableHeader.ConnectProperties.PropertyLength.Length()
	remainingLen += packet.VariableHeader.ConnectProperties.PropertyLength.ToValue()

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
	var err error
	var buffer = bytes.NewBuffer(nil)
	// fix header
	_, err = c.FixHeader.PacketType.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	_, err = c.FixHeader.RemainingLength.ToStream(buffer)
	if err != nil {
		return 0, err
	}

	// variable header
	_, err = c.VariableHeader.ProtocolName.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	_, err = c.VariableHeader.ProtocolLevel.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	_, err = c.VariableHeader.ConnectFlags.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	_, err = c.VariableHeader.KeepAlive.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	_, err = c.VariableHeader.ConnectProperties.ToStream(buffer)
	if err != nil {
		return 0, err
	}

	// payload
	_, err = c.Payload.ClientID.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	// there is will
	if c.VariableHeader.ConnectFlags.ToValue()&CONNECT_FLAG_WILLFLAG > 0x00 {
		_, err = c.Payload.WillProperties.ToStream(buffer)
		if err != nil {
			return 0, err
		}
		_, err = c.Payload.WillTopic.ToStream(buffer)
		if err != nil {
			return 0, err
		}
		_, err = c.Payload.WillPayload.ToStream(buffer)
		if err != nil {
			return 0, err
		}
	}
	// there is username
	if c.VariableHeader.ConnectFlags.ToValue()&CONNECT_FLAG_USERNAME > 0x00 {
		_, err = c.Payload.UserName.ToStream(buffer)
		if err != nil {
			return 0, err
		}
	}
	// there is password
	if c.VariableHeader.ConnectFlags.ToValue()&CONNECT_FLAG_PASSWORD > 0x00 {
		_, err = c.Payload.Password.ToStream(buffer)
		if err != nil {
			return 0, err
		}
	}
	// to stream
	return output.Write(buffer.Bytes())
}
