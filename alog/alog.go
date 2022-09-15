package alog

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

type Alog struct {
	msgChan chan string
	stop    chan bool
	writer  io.Writer
	mu      *sync.Mutex
	wg      *sync.WaitGroup
	prefix  string
}

func NewAlog(w io.Writer, p string) *Alog {
	return &Alog{msgChan: make(chan string), stop: make(chan bool), writer: w, mu: &sync.Mutex{}, wg: &sync.WaitGroup{}, prefix: p}
}

func (l *Alog) StartLogging() {
	l.wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		wg := sync.WaitGroup{}
	Loop:
		for {
			select {
			case msg := <-l.msgChan:
				wg.Add(1)
				go l.write(msg, &wg)
			case <-l.stop:
				wg.Wait()
				close(l.msgChan)
				break Loop
			}
		}
	}(l.wg)
}

func (l *Alog) write(str string, wg *sync.WaitGroup) {
	//str = fmt.Sprintf("%s  %s  MSG: %s", time.Now().Format(time.RFC1123), l.prefix, str)
	defer wg.Done()
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer.Write([]byte(str))
}

func (l *Alog) Printf(format string, v ...any) {
	str := fmt.Sprintf(format, v...)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	l.msgChan <- str
}

func (l *Alog) Println(str string) {
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	l.msgChan <- str
}
func (l *Alog) StopLogging() {
	l.stop <- true
	l.wg.Wait()
	close(l.stop)
}

func (l *Alog) SetPrefix(newPrefix string) {
	l.prefix = newPrefix
}

func (l *Alog) ChangeWriter(writer io.Writer) {
	l.writer = writer
}
