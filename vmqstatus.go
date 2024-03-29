package vmq

import "errors"

// vmq state machine
type VMQ_STATUS int

const (
	STATUS_IDLE         VMQ_STATUS = iota // no net conn
	STATUS_CONNECTING                     // no conn ack
	STATUS_CONNECTED                      // recive connack
	STATUS_DISCONNECTED                   // disconnected
)

func (v *vmq) setStatus(s VMQ_STATUS) error {
	// idle => connecting
	if v.status == STATUS_IDLE && s == STATUS_CONNECTING {
		v.status = STATUS_CONNECTING
		return nil
	}
	// connecting => connected
	if v.status == STATUS_CONNECTING && s == STATUS_CONNECTED {
		v.status = STATUS_CONNECTED
		v.onConnect(v)
		return nil
	}
	// => STATUS_DISCONNECTED
	if s == STATUS_DISCONNECTED {
		v.status = STATUS_DISCONNECTED
		v.onDisConnect(v)
		return nil
	}
	// => idle (may cause by error)
	if s == STATUS_IDLE {
		// TODO clean some data here
		v.status = STATUS_IDLE
		v.onError(v)
		return nil
	}
	return errors.New("vmq status error")
}
