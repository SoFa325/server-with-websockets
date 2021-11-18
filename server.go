package main

import (
	"flag"
	"fmt"
	"strings"
	"os/exec"
	"bufio"
	"net"
	"github.com/tatsushid/go-fastping"
    "time"
)

func trace (conn *net.TCPConn, h string) {
    cmd := exec.Command("tracert", h)
    fmt.Println(h)
    dateOut, _ := cmd.Output()
    fmt.Println(dateOut)
    
    start := 0
    for i := 0; i < len(dateOut); i++ {
        if dateOut[i] == byte(10) {
            //fmt.Println(string(dateOut[start:i]))
            conn.Write([]byte(string(dateOut[start:i])+"\n"))
            start = i + 1
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func Pinger(conn *net.TCPConn, h string) {
    p := fastping.NewPinger()
    ra, _ := net.ResolveIPAddr("ip4:icmp", h)
    p.AddIPAddr(ra)
    p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
        out := fmt.Sprintf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
        conn.Write([]byte(out+"\n"))
    }
    p.OnIdle = func() {
        out := fmt.Sprintf("finish")
        conn.Write([]byte(out+"\n"))
    }
    p.RunLoop()
    time.Sleep(time.Millisecond*3000)
    p.Stop()
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
func serve(conn *net.TCPConn) {
	defer conn.Close()
    z:= bufio.NewReader(conn)
	for {
		message, _ := z.ReadString('\n')
    	// Распечатываем полученое сообщение
    	fmt.Print("Message Received:", string(message))
    	// Процесс выборки для полученной строки
		if (strings.HasPrefix(message, "ping")) {
			go Pinger(conn, message[4:len(message)-1])
		}
		if (strings.HasPrefix(message, "trace")) {
			go trace(conn, message[5:len(message)-1])
		}
	}
}


func main() {
    // Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:9000", "specify ip address and port")
	flag.Parse()

    // Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		fmt.Println("address resolution failed", "address", addrStr)
	} else {
		fmt.Println("resolved TCP address", "address", addr.String())

        // Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			fmt.Println("listening failed", "reason", err)
		} else {
            // Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					fmt.Println("cannot accept connection", "reason", err)
				} else {
					fmt.Println("accepted connection", "address", conn.RemoteAddr().String())

                    // Запуск go-программы для обслуживания клиентов.
					go serve(conn)
				}
			}
		}
	}
}
