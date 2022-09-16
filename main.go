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
		alog2.Printf("I'm " + strconv.Itoa(goNum) + " async. " + "My msg num is: " + strconv.Itoa(j))
	}
}

func testSync(goNum int, messageNum int, w *sync.WaitGroup) {
	defer w.Done()
	for j := 0; j < messageNum; j++ {
		log.Println("I'm " + strconv.Itoa(goNum) + " sync. " + "My msg num is: " + strconv.Itoa(j))
	}
}

func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Println("args != 2")
		return
	}
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
	wg2 := sync.WaitGroup{}

	//sync
	start := time.Now()
	for i := 0; i < threadsNum; i++ {
		wg2.Add(1)
		go testSync(i, messageNum, &wg2)
	}
	wg2.Wait()
	workTimeSync := time.Since(start)

	//async
	s := alog.NewAlog(os.Stdout, "Debug")
	wg := sync.WaitGroup{}

	start = time.Now()
	s.StartLogging()
	for i := 0; i < threadsNum; i++ {
		wg.Add(1)
		go testAsync(i, messageNum, &wg, s)
	}
	wg.Wait()
	s.StopLogging()
	workTimeAsync := time.Since(start)

	fmt.Println("Work async time: " + workTimeAsync.String())
	fmt.Println("Work sync time: " + workTimeSync.String())
	fmt.Scanln()
}
