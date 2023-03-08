package xlog

import "testing"

func TestNewProduceLogger(t *testing.T) {

	l := NewProduceLogger()
	l.Info("aaaaaaaaaaaaaaaaaa")

}
