package bar

type ConsoleProgOption func(prog *ConsoleProg)

// SetLength 设置进度条的长度
func SetLength(length int) ConsoleProgOption {
	return func(prog *ConsoleProg) {
		prog.prog = make([]byte, length)
	}
}

// SetDoneStr 设置进度条的完成字符
func SetDoneStr(s byte) ConsoleProgOption {
	return func(prog *ConsoleProg) {
		prog.doneS = s
	}
}

// SetBlankStr 设置进度条的背景字符
func SetBlankStr(s byte) ConsoleProgOption {
	return func(prog *ConsoleProg) {
		prog.blankS = s
	}
}

// SetFPS 设置进度条的刷新频率
func SetFPS(fps uint8) ConsoleProgOption {
	return func(prog *ConsoleProg) {
		prog.fps = fps
	}
}
