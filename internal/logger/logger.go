package logger

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type output struct {
	writer io.Writer
	closer io.Closer
}

var (
	mu      sync.Mutex
	outputs = map[string]output{}
)

func AddOutput(name string, w io.Writer) {
	mu.Lock()
	defer mu.Unlock()

	out := output{writer: w}
	if closer, ok := w.(io.Closer); ok {
		out.closer = closer
	}
	outputs[name] = out
}

func RemoveOutput(name string) {
	mu.Lock()
	defer mu.Unlock()

	out, ok := outputs[name]
	if !ok {
		return
	}
	if out.closer != nil {
		out.closer.Close()
	}
	delete(outputs, name)
}

func Println(v ...any) {
	write(fmt.Sprintln(v...))
}

func Printf(format string, v ...any) {
	write(fmt.Sprintf(format, v...))
}

func write(message string) {
	mu.Lock()
	defer mu.Unlock()

	line := fmt.Sprintf("%s %s", time.Now().Format("2006/01/02 15:04:05"), message)
	for name, out := range outputs {
		if _, err := io.WriteString(out.writer, line); err != nil {
			if out.closer != nil {
				out.closer.Close()
			}
			delete(outputs, name)
		}
	}
}
