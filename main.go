package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				clear_rule_v4()
				log.Fatalln(s)
			}
		}
	}()

	var local_port_64 uint64
	var interval int
	var input string
	var client string
	var address string
	var username string
	var password string
	var print_version bool
	flag.Uint64Var(&local_port_64, "l", 12345, "local port")
	flag.IntVar(&interval, "i", 10, "interval between two stun request in second")
	flag.StringVar(&input, "w", "./external.port", "file path of external port")
	flag.StringVar(&client, "c", "", "bittorrent client name")
	flag.StringVar(&address, "s", "127.0.0.1:8080", "bittorrent client address <ip:port>")
	flag.StringVar(&username, "u", "admin", "bittorrent client username")
	flag.StringVar(&password, "p", "123456", "bittorrent client password")
	flag.BoolVar(&print_version, "v", false, "show current version")
	flag.Parse()
	if print_version {
		println("trav-bittorrent", version)
		os.Exit(0)
	}
	var local_port uint16 = uint16(local_port_64)

	set_conf("local_port", local_port)
	set_conf("interval", interval)
	set_conf("input", input)
	set_conf("client", client)
	set_conf("address", address)
	set_conf("username", username)
	set_conf("password", password)

	set_conf("X-Transmission-Session-Id", "75xsfCvW52BecBVFhpWy6M1jY5oJBYdCz53rbeF6S5FzYVUx")

	var external_port uint16 = 0
	for ; ; time.Sleep(time.Duration(get_conf("interval").(int)) * time.Second) {

		bytes, err := os.ReadFile(get_conf("input").(string))
		if err != nil {
			log.Println("read file failed:", err)
			continue
		}

		port, err := strconv.ParseUint(string(bytes), 10, 16)
		if err != nil {
			log.Println("invalid file content:", err)
			continue
		}

		if port < 1024 || port > 65535 {
			log.Println("invalid port:", port)
			continue
		}

		if uint16(port) == external_port {
			log.Println("no change")
			continue
		}

		log.Println("modifying listening port...")
		switch client {
		case "qbittorrent":
			qbittorrent_port(uint16(port))
		case "transmission":
			transmission_port(uint16(port))
		default:
			log.Fatalln("unsupported client:", client)
		}

		log.Println("modifying ipv4 iptables redirect rule...")
		set_rule_v4(uint16(port))

		external_port = uint16(port)
	}

}

func local_ip(server_addr string) net.IP {

	conn, err := net.Dial("udp4", server_addr)
	if err != nil {
		log.Println("failed to get local ip")
		log.Fatalln(err)
	}
	defer conn.Close()

	local_addr := conn.LocalAddr().(*net.UDPAddr)
	return local_addr.IP
}
