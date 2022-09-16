package alog

import (
	"fmt"
	"io"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Alog struct {
	msgChan chan string
	writer  io.Writer
	wg      *sync.WaitGroup
	prefix  string
	mu      *sync.Mutex
}

func NewAlog(w io.Writer, p string) *Alog {
	if p != "" {
		p += " "
	}
	return &Alog{msgChan: make(chan string), wg: &sync.WaitGroup{}, writer: w, prefix: p, mu: &sync.Mutex{}}
}

func (l *Alog) StartLogging() {
	l.wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		for msg := range l.msgChan {
			l.writer.Write(l.formatHeader(msg))
		}
	}(l.wg)

}

func (l *Alog) formatHeader(line any) []byte {
	now := time.Now()
	year, month, day := now.Date()
	date := strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + strconv.Itoa(day) + " "
	hour, min, sec := now.Clock()
	t := strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec) + " "
	return []byte(fmt.Sprintf("%s%s%v", date+t, l.prefix, line))
}

func (l *Alog) Println(in ...any) {
	str := fmt.Sprintln(in...)
	l.msgChan <- str
}

func (l *Alog) Printf(format string, v ...any) {
	str := fmt.Sprintf(format, v...)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	l.msgChan <- str
}

func (l *Alog) Fatalln(in ...any) {
	str := fmt.Sprintln(in...)
	l.writer.Write([]byte(str))
	l.StopLogging()
	os.Exit(1)
}

func (l *Alog) Fatalf(format string, v ...any) {
	str := fmt.Sprintf(format, v...)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	l.writer.Write([]byte(str))
	l.StopLogging()
	os.Exit(1)
}

func (l *Alog) Panicln(v ...any) {
	str := fmt.Sprintln(v...)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	l.writer.Write([]byte(str))
	l.StopLogging()
	panic(str)
}

func (l *Alog) Panicf(format string, v ...any) {
	str := fmt.Sprintf(format, v...)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	l.writer.Write([]byte(str))
	l.StopLogging()
	panic(str)
}

func (l *Alog) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix + " "
}

func (l *Alog) StopLogging() {
	close(l.msgChan)
	l.wg.Wait()
}
