package bar

import (
	"testing"
	"time"
)

func TestCreateConsoleProg(t *testing.T) {
	prog := CreateConsoleProg("test", 100, SetLength(25))
	for i := 0; i < 100; i++ {
		time.Sleep(200 * time.Millisecond)
		_ = prog.Add(1)
	}
	prog.Close(time.Second)
}
