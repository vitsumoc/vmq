package packets

import t "github.com/vitsumoc/vmq/types"

type CONNECT_PACKET struct {
	FixHeader CONNECT_FIX_HEADER
}

type CONNECT_FIX_HEADER struct {
	PacketType      t.MQTT_BYTE
	RemainingLength t.MQTT_VAR_INT
}

type CONNECT_VARIABLE_HEADER struct {
}
