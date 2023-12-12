package utils

import "time"

type CloseType int8

const (
	SuccessClose = iota
	TimeoutClose
)

type CloseCB struct {
	closing, closed chan struct{}
}

func NewStopCB() CloseCB {
	return CloseCB{
		closing: make(chan struct{}, 1),
		closed:  make(chan struct{}, 1),
	}
}

func (cb CloseCB) Close() {
	select {
	case cb.closing <- struct{}{}:
	default:
	}
}

func (cb CloseCB) Closing() <-chan struct{} {
	return cb.closing
}

func (cb CloseCB) WaitClosed(duration time.Duration) CloseType {
	t := time.After(duration)
	select {
	case <-t:
		return TimeoutClose
	case <-cb.closed:
		return SuccessClose
	}
}
