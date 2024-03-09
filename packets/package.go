// enum about packet
package packets

// MQTT Control Packet type
const (
	PACKET_TYPE_CONNECT     byte = 0x10
	PACKET_TYPE_CONNACK     byte = 0x20
	PACKET_TYPE_PUBLISH     byte = 0x30
	PACKET_TYPE_PUBACK      byte = 0x40
	PACKET_TYPE_PUBREC      byte = 0x50
	PACKET_TYPE_PUBREL      byte = 0x60
	PACKET_TYPE_PUBCOMP     byte = 0x70
	PACKET_TYPE_SUBSCRIBE   byte = 0x80
	PACKET_TYPE_SUBACK      byte = 0x90
	PACKET_TYPE_UNSUBSCRIBE byte = 0xA0
	PACKET_TYPE_UNSUBACK    byte = 0xB0
	PACKET_TYPE_PINGREQ     byte = 0xC0
	PACKET_TYPE_PINGRESP    byte = 0xD0
	PACKET_TYPE_DISCONNECT  byte = 0xE0
	PACKET_TYPE_AUTH        byte = 0xF0
)
