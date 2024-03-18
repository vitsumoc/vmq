package packets

import (
	"bytes"
	"io"

	t "github.com/vitsumoc/vmq/types"
)

type DISCONNECT_PACKET struct {
	FixHeader      DISCONNECT_FIX_HEADER
	VariableHeader DISCONNECT_VARIABLE_HEADER
}

type DISCONNECT_FIX_HEADER struct {
	PacketType      t.MQTT_BYTE
	RemainingLength t.MQTT_VAR_INT
}

type DISCONNECT_VARIABLE_HEADER struct {
	DisconnectReasonCode t.MQTT_BYTE
	DisconnectProperties PROPERTIES
}

func NewDisconnectPacketFromConf(dc *DisconnectConf) *DISCONNECT_PACKET {
	packet := DISCONNECT_PACKET{}
	// fix header
	packet.FixHeader.PacketType.FromValue(byte(PACKET_TYPE_DISCONNECT))
	remainingLen := 0

	// variable header
	// Byte 1 in the Variable Header is the Disconnect Reason Code.
	// If the Remaining Length is less than 1 the value of 0x00 (Normal disconnection) is used.
	if dc.disconnectReasonCode != RC_NORMAL || dc.disconnectProperties.Length() > 1 {
		packet.VariableHeader.DisconnectReasonCode.FromValue(byte(dc.disconnectReasonCode))
		remainingLen += packet.VariableHeader.DisconnectReasonCode.Length()
		packet.VariableHeader.DisconnectProperties = dc.disconnectProperties
		remainingLen += packet.VariableHeader.DisconnectProperties.Length()
	}

	// set remainingLen
	packet.FixHeader.RemainingLength.FromValue(remainingLen)

	return &packet
}

func (d *DISCONNECT_PACKET) ToStream(output io.Writer) (int, error) {
	var err error
	var buffer = bytes.NewBuffer(nil)
	// fix header
	_, err = d.FixHeader.PacketType.ToStream(buffer)
	if err != nil {
		return 0, err
	}
	_, err = d.FixHeader.RemainingLength.ToStream(buffer)
	if err != nil {
		return 0, err
	}

	// variable header
	if d.FixHeader.RemainingLength.ToValue() > 0 {
		_, err = d.VariableHeader.DisconnectReasonCode.ToStream(buffer)
		if err != nil {
			return 0, err
		}
		_, err = d.VariableHeader.DisconnectProperties.ToStream(buffer)
		if err != nil {
			return 0, err
		}
	}

	// to stream
	return output.Write(buffer.Bytes())
}

func NewDisconnectPacket(packetType *t.MQTT_BYTE, remainingLength *t.MQTT_VAR_INT) *DISCONNECT_PACKET {
	return &DISCONNECT_PACKET{
		FixHeader: DISCONNECT_FIX_HEADER{
			*packetType,
			*remainingLength,
		},
		VariableHeader: DISCONNECT_VARIABLE_HEADER{
			DisconnectReasonCode: *t.NewByte().FromValue(byte(RC_NORMAL)),
			DisconnectProperties: *NewProperties(),
		},
	}
}

func (d *DISCONNECT_PACKET) FromStream(input io.Reader) (int, error) {
	if d.FixHeader.RemainingLength.ToValue() == 0 {
		return 0, nil
	}
	length := 0
	n, err := d.VariableHeader.DisconnectReasonCode.FromStream(input)
	if err != nil {
		return 0, err
	}
	length += n
	n, err = d.VariableHeader.DisconnectProperties.FromStream(input)
	if err != nil {
		return 0, err
	}
	length += n
	return length, nil
}
