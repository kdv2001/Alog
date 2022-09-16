package alog

import (
	"fmt"
	"io"
	_ "net/http/pprof"
	"strings"
	"sync"
	"time"
)

type Alog struct {
	msgChan chan string
	//stop    chan bool
	writer io.Writer
	//mu      *sync.Mutex
	wg     *sync.WaitGroup
	numgos int
	//prefix  string	stop    chan bool
}

func NewAlog(w io.Writer, p string, numgos int) *Alog {
	return &Alog{msgChan: make(chan string), numgos: numgos, wg: &sync.WaitGroup{}, writer: w /*, stop: make(chan bool), writer: w, mu: &sync.Mutex{}, wg: &sync.WaitGroup{}, prefix: p*/}
}

func (l *Alog) StartLogging() {
	l.wg.Add(l.numgos)
	for i := 0; i < l.numgos; i++ {
		go func(w *sync.WaitGroup) {
			defer w.Done()
			for msg := range l.msgChan {
				//fmt.Fprintln(l.writer, msg)
				msg = fmt.Sprintf("%s  MSG: %s\n", time.Now().String(), msg)

				l.writer.Write([]byte(msg))
			}
		}(l.wg)
	}
}

func (l *Alog) Printf(format string, v ...any) {
	str := fmt.Sprintf(format, v...)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	l.msgChan <- str
}

func (l *Alog) Println(in any) {
	str := fmt.Sprintf("%s", in)

	l.msgChan <- str
}
func (l *Alog) StopLogging() {
	close(l.msgChan)
	l.wg.Wait()
}

/*func (l *Alog) SetPrefix(newPrefix string) {
	l.prefix = newPrefix
}

func (l *Alog) ChangeWriter(writer io.Writer) {
	l.writer = writer
}*/
