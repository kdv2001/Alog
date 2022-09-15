package main

import (
	"alog/alog"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func testAsync(goNum int, messageNum int, w *sync.WaitGroup, alog2 *alog.Alog) {
	defer w.Done()
	for j := 0; j < messageNum; j++ {
		alog2.Println("i'm " + strconv.Itoa(goNum) + " go. " + "My msg num is: " + strconv.Itoa(j))
	}
}

func testSync(goNum int, messageNum int, w *sync.WaitGroup) {
	defer w.Done()
	for j := 0; j < messageNum; j++ {
		log.Println("i'm " + strconv.Itoa(goNum) + " go. " + "My msg num is: " + strconv.Itoa(j))
	}
}

func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Println("args != 2")
		return
	}
	s := alog.NewAlog(os.Stdout, "DEBUG")
	s.StartLogging()

	wg := sync.WaitGroup{}

	threadsNum, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("second argument is not int")
		return
	}
	messageNum, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("third argument is not int")
		return
	}
	start := time.Now()
	for i := 0; i < threadsNum; i++ {
		wg.Add(1)
		go testAsync(i, messageNum, &wg, s)
	}
	wg.Wait()
	s.StopLogging()

	workTimeAsync := time.Since(start)
	start = time.Now()
	for i := 0; i < threadsNum; i++ {
		wg.Add(1)
		go testSync(i, messageNum, &wg)
	}
	wg.Wait()
	workTimeSync := time.Since(start)
	fmt.Println("work async time " + workTimeAsync.String())
	fmt.Println("work sync time " + workTimeSync.String())
	fmt.Scanln()
}
