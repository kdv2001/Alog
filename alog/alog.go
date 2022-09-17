package alog

import (
	"fmt"
	"io"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"time"
)

type Alog struct {
	msgChan      chan string
	writer       io.Writer
	loggerIsOpen bool
	wg           *sync.WaitGroup
	prefix       string
	mu           *sync.Mutex
}

func NewAlog(w io.Writer, p string, bufferSize int) *Alog {
	if p != "" {
		p += " "
	}
	return &Alog{msgChan: make(chan string, bufferSize), wg: &sync.WaitGroup{}, writer: w, prefix: p, mu: &sync.Mutex{}}
}

func (l *Alog) StartLogging() {
	l.wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		for msg := range l.msgChan {
			l.writer.Write([]byte(msg))
		}
	}(l.wg)
	l.loggerIsOpen = true
}

func (l *Alog) formatHeader(line any) string {
	now := time.Now()
	year, month, day := now.Date()
	date := strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + strconv.Itoa(day) + " "
	hour, min, sec := now.Clock()
	t := strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec) + " "
	return fmt.Sprintf("%s%s%v", date+t, l.prefix, line)
}

func (l *Alog) Println(in ...any) {
	if !l.loggerIsOpen {
		return
	}
	str := fmt.Sprintln(in...)
	l.msgChan <- l.formatHeader(str)
}

func (l *Alog) Printf(format string, v ...any) {
	if !l.loggerIsOpen {
		return
	}
	str := fmt.Sprintf(format, v...)
	l.msgChan <- l.formatHeader(str)
}

func (l *Alog) Fatalln(in ...any) {
	if !l.loggerIsOpen {
		return
	}
	str := fmt.Sprintln(in...)
	l.writer.Write([]byte(l.formatHeader(str)))
	l.StopLogging()
	os.Exit(1)
}

func (l *Alog) Fatalf(format string, v ...any) {
	if !l.loggerIsOpen {
		return
	}
	str := fmt.Sprintf(format, v...)
	l.writer.Write([]byte(l.formatHeader(str)))
	l.StopLogging()
	os.Exit(1)
}

func (l *Alog) Panicln(v ...any) {
	if !l.loggerIsOpen {
		return
	}
	str := fmt.Sprintln(v...)
	l.writer.Write([]byte(l.formatHeader(str)))
	l.StopLogging()
	panic(str)
}

func (l *Alog) Panicf(format string, v ...any) {
	if !l.loggerIsOpen {
		return
	}
	str := fmt.Sprintf(format, v...)
	l.writer.Write([]byte(l.formatHeader(str)))
	l.StopLogging()
	panic(str)
}

func (l *Alog) SetPrefix(prefix string) {
	if !l.loggerIsOpen {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix + " "
}

func (l *Alog) StopLogging() {
	close(l.msgChan)
	l.loggerIsOpen = false
	l.wg.Wait()
}
