package packets

import (
	"io"

	t "github.com/vitsumoc/vmq/types"
)

type SUBACK_PACKET struct {
	FixHeader      SUBACK_FIX_HEADER
	VariableHeader SUBACK_VARIABLE_HEADER
	Payload        SUBACK_PAYLOAD
}

type SUBACK_FIX_HEADER struct {
	PacketType      t.MQTT_BYTE
	RemainingLength t.MQTT_VAR_INT
}

type SUBACK_VARIABLE_HEADER struct {
	PacketId         t.MQTT_U16
	SubackProperties PROPERTIES
}

type SUBACK_PAYLOAD struct {
	ReasonCodes []t.MQTT_TYPE
}

func NewSubackPacket(packetType *t.MQTT_BYTE, remainingLength *t.MQTT_VAR_INT) *SUBACK_PACKET {
	return &SUBACK_PACKET{
		FixHeader: SUBACK_FIX_HEADER{
			PacketType:      *packetType,
			RemainingLength: *remainingLength,
		},
		VariableHeader: SUBACK_VARIABLE_HEADER{
			SubackProperties: *NewProperties(),
		},
		Payload: SUBACK_PAYLOAD{
			ReasonCodes: make([]t.MQTT_TYPE, 0),
		},
	}
}

func (sp *SUBACK_PACKET) FromStream(input io.Reader) (int, error) {
	length := 0
	n, err := sp.VariableHeader.PacketId.FromStream(input)
	if err != nil {
		return 0, err
	}
	length += n
	n, err = sp.VariableHeader.SubackProperties.FromStream(input)
	if err != nil {
		return 0, err
	}
	length += n
	for length < sp.FixHeader.RemainingLength.ToValue() {
		n, err = t.NewByte().FromStream(input)
		if err != nil {
			return 0, err
		}
		length += n
	}
	return length, nil
}
