package packets

import (
	"io"

	t "github.com/vitsumoc/vmq/types"
)

const CONNACK_FLAG_SESSIONPRESENT byte = 0x01

type CONNACK_PACKET struct {
	FixHeader      CONNACK_FIX_HEADER
	VariableHeader CONNACK_VARIABLE_HEADER
}

type CONNACK_FIX_HEADER struct {
	PacketType      t.MQTT_BYTE
	RemainingLength t.MQTT_VAR_INT
}

type CONNACK_VARIABLE_HEADER struct {
	ConnectAcknowledgeFlags t.MQTT_BYTE
	ConnectReasonCode       t.MQTT_BYTE
	Properties              PROPERTIES
}

func NewConnackPacket(packetType *t.MQTT_BYTE, remainingLength *t.MQTT_VAR_INT) *CONNACK_PACKET {
	return &CONNACK_PACKET{
		FixHeader: CONNACK_FIX_HEADER{
			PacketType:      *packetType,
			RemainingLength: *remainingLength,
		},
		VariableHeader: CONNACK_VARIABLE_HEADER{
			Properties: *NewProperties(),
		},
	}
}

func (c *CONNACK_PACKET) FromStream(input io.Reader) (int, error) {
	length := 0
	n, err := c.VariableHeader.ConnectAcknowledgeFlags.FromStream(input)
	if err != nil {
		return 0, err
	}
	length += n
	n, err = c.VariableHeader.ConnectReasonCode.FromStream(input)
	if err != nil {
		return 0, err
	}
	length += n
	n, err = c.VariableHeader.Properties.FromStream(input)
	if err != nil {
		return 0, err
	}
	length += n
	return length, nil
}
