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
var udpDatos *net.UDPConn

// se asigna en build.sh
var GIT string

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println(GIT)
		return
	}
	// carga el config.json
	gloApp.Load()

	go udpDatosLoop()

	udpPingLoop()

}

var datosAddr *net.UDPAddr

func udpDatosLoop() {

	host := "0.0.0.0"

	var err error

	udpDatos, err = net.ListenUDP("udp",
		&net.UDPAddr{
			IP:   net.ParseIP(host),
			Port: gloApp.DatosPort})

	if err != nil {
		fmt.Println("ERR ListenUDP:", err)
		return
	}

	fmt.Println("UDP escuchando Datos en ", udpDatos.LocalAddr())

	buf := make([]byte, 1024)

	for {
		_, remote, err := udpDatos.ReadFrom(buf)
		if err != nil {
			fmt.Println("ERR udpListenLoop:", err)
			return

		}

		var ok bool
		datosAddr, ok = remote.(*net.UDPAddr)
		if !ok {
			fmt.Println("ERR No UDPAddr in read packet. (Windows?)")
			return
		}

		fmt.Printf("Datos desde %+v \n", datosAddr)
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

		fmt.Printf("PING desde %+v \n", udpaddr)

		respuestaPing(udpaddr, "PONG 1 \n")
		respuestaPing(udpaddr, "PONG 2 \n")
		respuestaPing(udpaddr, "PONG 3 \n")

		respuestaDatos(datosAddr, "Datos 1 \n")
		respuestaDatos(udpaddr, "Datos 2 \n")
		respuestaDatos(datosAddr, "Datos 3 \n")

		delay := 50
		fmt.Printf("Delay %v segundos", delay)
		time.Sleep(time.Second * time.Duration(delay))

		respuestaDatos(datosAddr, "FIN........ \n")
		respuestaPing(udpaddr, "FIN........ \n")

	}
}

func respuestaPing(udpaddr *net.UDPAddr, texto string) {

	time.Sleep(1 * time.Second)
	log.Printf("Enviando PING %v a %v", texto, udpaddr)

	buf := []byte(texto)
	udpPing.WriteTo(buf, udpaddr)

}
func respuestaDatos(remoteAddr *net.UDPAddr, texto string) {

	time.Sleep(1 * time.Second)
	log.Printf("Enviando Datos %v a %v", texto, remoteAddr)

	buf := []byte(texto)

	udpDatos.WriteTo(buf, remoteAddr)

}
