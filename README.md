# Alog

Alog это пакет для асинхронной записи логов в вашем приложение. В чем её преимущество, спросите вы.
Её преимущество в том, что во время записи лога, не происходит прерывания работы основной программы.
Сообщение, которое необходимо записать, передается в отдельную горотину и запись происходит там

# Пример

```go
func testAsync(goNum int, messageNum int, w *sync.WaitGroup, alog2 *alog.Alog) {
	defer w.Done()
	for j := 0; j < messageNum; j++ {
		//вызов метода записи
        alog2.Println("I'm " + strconv.Itoa(goNum) + " async. " + "Msg num is: " + strconv.Itoa(j))
    }
}


func main() {
    threadsNum := 3
	messageNum := 50
	bufferSize := 50
	prefix := "Info"
    s := alog.NewAlog(os.Stdout, prefix, bufferSize) //создаем объект логера 
    wg := sync.WaitGroup{}
    s.StartLogging()
    for i := 0; i < threadsNum; i++ { //запускаем пишущие горутины 
        wg.Add(1)
        go testAsync(i, messageNum, &wg, s) 
    }
    wg.Wait()
    s.StopLogging() // останавливаем логгер
}
```
В данном примере мы запускаем 3 горотины, которые одновременно будут записывать по 50 логов каждый.

Пимер вывода:
```
    2022/9/17 11:20:33 I'm 0 async. My msg num is: 31
    2022/9/17 11:20:33 I'm 2 async. My msg num is: 43
    2022/9/17 11:20:33 I'm 1 async. My msg num is: 27
    2022/9/17 11:20:33 I'm 2 async. My msg num is: 44
    2022/9/17 11:20:33 I'm 0 async. My msg num is: 32
```
