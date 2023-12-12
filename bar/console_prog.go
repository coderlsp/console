package bar

import (
	"console/utils"
	"fmt"
	"time"
)

type ConsoleProg struct {
	name        string
	total       float64
	cur         float64
	prog        []string
	doneS       string // 已完成的进度字符
	blankS      string // 未完成的进度字符
	fps         uint8  // 进度条的刷新频率
	dataChan    chan float64
	dataStopCB  utils.CloseCB
	printStopCB utils.CloseCB
	begin       time.Time
}

func CreateConsoleProg(name string, total float64, opts ...ConsoleProgOption) *ConsoleProg {
	prog := &ConsoleProg{
		name:        name,
		total:       total,
		cur:         0,
		dataChan:    make(chan float64, 1000),
		prog:        make([]string, 100),
		doneS:       "#",
		blankS:      "_",
		fps:         20,
		dataStopCB:  utils.NewStopCB(),
		printStopCB: utils.NewStopCB(),
		begin:       time.Now(),
	}
	for _, opt := range opts {
		opt(prog)
	}
	prog.run()
	return prog
}

func (prog *ConsoleProg) Close(duration time.Duration) {
	close(prog.dataChan)
	prog.dataStopCB.Close()
	prog.dataStopCB.WaitClosed(duration)
	prog.print()
	prog.printStopCB.Close()
	prog.printStopCB.WaitClosed(duration)
}

func (prog *ConsoleProg) runPrint() {
	go func() {
		t := time.NewTicker(time.Duration(float64(time.Second) / float64(prog.fps)))
		defer func() {
			t.Stop()
		}()
		for {
			select {
			case <-t.C:
				prog.print()
			case <-prog.printStopCB.Closing():
				return
			}
		}
	}()
}

func (prog *ConsoleProg) print() {
	percent := prog.cur / prog.total
	ids := int(percent * float64(len(prog.prog)))
	for i := 0; i < len(prog.prog); i++ {
		if ids > 0 && i <= ids {
			prog.prog[i] = prog.doneS
		} else {
			prog.prog[i] = prog.blankS
		}
	}

	fmt.Printf("\r[%s]%.2f/%.2f [cost]%s", prog.prog, prog.cur, prog.total, time.Since(prog.begin).String())
}

func (prog *ConsoleProg) run() {
	prog.handle()
	prog.runPrint()
}

func (prog *ConsoleProg) Add(f float64) error {
	select {
	case prog.dataChan <- f:
	default:
		return fmt.Errorf("channel full, size:%d", cap(prog.dataChan))
	}
	return nil
}

func (prog *ConsoleProg) handle() {
	go func() {
		defer func() {
			for data := range prog.dataChan {
				prog.cur += data
			}
		}()
		for {
			select {
			case data, ok := <-prog.dataChan:
				if !ok {
					return
				}
				prog.cur += data
			case <-prog.dataStopCB.Closing():
				return
			}
		}
	}()
}
