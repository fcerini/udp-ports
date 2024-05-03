package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var gloApp AppConfig
var udpPing *net.UDPConn
var udpPong *net.UDPConn

// se asigna en build.sh
var GIT string

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println(GIT)
		return
	}
	// carga el config.json
	gloApp.Load()

	go udpPongLoop()

	udpPingLoop()

}

func udpPongLoop() {

	host := "0.0.0.0"

	var err error

	udpPong, err = net.ListenUDP("udp",
		&net.UDPAddr{
			IP:   net.ParseIP(host),
			Port: gloApp.PongPort})

	if err != nil {
		fmt.Println("ERR ListenUDP:", err)
		return
	}

	fmt.Println("PONG en ", udpPong.LocalAddr())

	buf := make([]byte, 1024)

	for {
		_, remote, err := udpPong.ReadFrom(buf)
		if err != nil {
			fmt.Println("ERR udpListenLoop:", err)
			return

		}

		udpaddr, ok := remote.(*net.UDPAddr)
		if !ok {
			fmt.Println("ERR No UDPAddr in read packet. (Windows?)")
			return
		}

		fmt.Printf("Nuevo PONG? desde %+v \n", udpaddr)
	}
}

// Listen for and handle UDP packets.
func udpPingLoop() {

	host := "0.0.0.0"

	var err error

	udpPing, err = net.ListenUDP("udp",
		&net.UDPAddr{
			IP:   net.ParseIP(host),
			Port: gloApp.PingPort})

	if err != nil {
		fmt.Println("ERR ListenUDP:", err)
		return
	}

	fmt.Println("UDP escuchando PING en ", udpPing.LocalAddr())

	buf := make([]byte, 1024)

	for {
		_, remote, err := udpPing.ReadFrom(buf)
		if err != nil {
			fmt.Println("ERR udpListenLoop:", err)
			return

		}

		udpaddr, ok := remote.(*net.UDPAddr)
		if !ok {
			fmt.Println("ERR No UDPAddr in read packet. (Windows?)")
			return
		}

		fmt.Printf("Nuevo PING desde %+v \n", udpaddr)

		//udpaddr.Port = 6666
		//fmt.Printf("Respondiendo a %+v \n", udpaddr)

		send(udpaddr, "PONG 1 \n")
		send(udpaddr, "PONG 2 \n")
		send(udpaddr, "PONG 3 \n")
		send(udpaddr, "FIN........ \n")

	}
}

func send(udpaddr *net.UDPAddr, resultado string) {

	time.Sleep(1 * time.Second)
	log.Printf("Enviando %v ", resultado)

	buf := []byte(resultado)
	udpPong.WriteTo(buf, udpaddr)

}
