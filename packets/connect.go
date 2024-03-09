package packets

type CONNECT_PACKET struct {
	FixHeader CONNECT_FIX_HEADER
}

type CONNECT_FIX_HEADER struct {
	PacketType      byte
	RemainingLength remainingLength
}

type CONNECT_VARIABLE_HEADER struct {
}
