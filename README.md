# Alog

Alog это пакет для асинхронной записи логов в вашем приложение. В чем её преимущество, спросите вы.
Её преимущество в том, что во время записи лога, не происходит прерывания работы основной программы.
Сообщение, которое необходимо записать, передается в отдельную горотину и запись происходит там

# Тестовый запуск

Для тестового запуска выполните:

```
go build -o testAlog .
./testAlog [gorutineNum] [logNum]
```

где gorutineNum - количество горутин, пишущих логи, logNum - количество логов для записи

# Пример

```go
func testAsync(goNum int, messageNum int, w *sync.WaitGroup, alog *alog.Alog) {
    defer w.Done()
    for j := 0; j < messageNum; j++ {
        //вызов метода записи
        alog.Println("I'm " + strconv.Itoa(goNum) + " async. " + "Msg num is: " + strconv.Itoa(j))
    }
}


func main() {
    threadsNum := 3 //количество пишущих горутин
    messageNum := 50 //количество записанных логов
    bufferSize := 50 //размер буфера канала
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

# Методы

Данный пакет поддерживает следующие методы:

* `NewAlog(w io.Writer, p string, bufferSize int) *Alog` - создает объект логера. Здесь можно задать интерфейс вывода,
  префикс сообщения, размер буфера канала.
* `StartLogging()` - запускает горутину логирования
* `Println(in ...any)` - вывод лога в io.writer
* `Printf(format string, v ...any)` - вывод лога в io.writer
* `Fatalln(in ...any)` - вывод лога в io.writer, c завершением программы(os.exit(1))
* `Fatalf(format string, v ...any)` - вывод лога в io.writer, c завершением программы(os.exit(1))
* `Panicln(v ...any)` - вывод лога в io.writer, c выбросом паники(данный метод не завершает горутину логирования, это
  нужно сделать самостоятельно)
* `Panicf(format string, v ...any)` - вывод лога в io.writer, c выбросом паники(данный метод не завершает горутину
  логирования, это нужно сделать самостоятельно)
* `SetPrefix(prefix string)` - смена префикса
* `StopLogging()` - остановка логера

# Формат лога

```
 2022/9/17 11:20:33 [prefix] [message]
```