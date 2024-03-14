package types

import "io"

type MQTT_TYPE interface {
	Length() int
	FromStream(io.Reader) (int, error)
	ToStream(io.Writer) (int, error)
}

// all types defined in MQTT standard
type MQTT_BYTE struct {
	data byte
}
type MQTT_U16 struct {
	data uint16
}
type MQTT_U32 struct {
	data uint32
}
type MQTT_UTF8 struct {
	len  MQTT_U16
	data []byte
}
type MQTT_VAR_INT struct {
	data []byte
}
type MQTT_BIN struct {
	len  MQTT_U16
	data []byte
}
type MQTT_UTF8_PAIR struct {
	key   MQTT_UTF8
	value MQTT_UTF8
}
