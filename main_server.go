package main

import (
  "fmt"
  "github.com/skorobogatov/input"
  "github.com/sparrc/go-ping"
)

func attack(s string) {
	pinger, err := ping.NewPinger(s)
	fmt.Println("Your are under attack")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	go attack(s)
	pinger.Run()
}

func main() {
  var n int // количество запросов
	var h string //хост
	input.Scanf("%s", &h)
  pinger, err := ping.NewPinger(h)
	input.Scanf("%d", &n)
	if n > 0 {
		pinger.Count = n
	}
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}
	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()

  /*
  for i := 0; i < 10; i++ {
		go attack(h)
	}
  //*/
	fmt.Println("Введите любой символ для завершения")
	var c string
	input.Scanf("%s", c)
}
