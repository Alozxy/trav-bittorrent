package main

import (
	"log"
	"os/exec"
	"strconv"

	"github.com/alozxy/udp-forward"
)

var forwarder *forward.Forwarder = nil

func clear_rule_v4() {

	if out, err := exec.Command("netsh", "interface", "portproxy", "reset").CombinedOutput(); err != nil {
		log.Fatalln("netsh return a non-zero value while clearing ipv4 rules:", string(out))
	}

	if forwarder != nil {
		forwarder.Close()
	}
}

func set_rule_v4() {
}

func modify_rule_v4(external_port uint16) {
	local_port := get_conf("local_port").(uint16)
	src_ip := local_ip("1.1.1.1:443")

	var err error

	if out, err := exec.Command("netsh", "interface", "portproxy", "add", "v4tov4", "listenport="+strconv.FormatUint(uint64(local_port), 10), "listenaddress="+src_ip.String(), "connectport="+strconv.FormatUint(uint64(external_port), 10), "connectaddress=127.0.0.1").CombinedOutput(); err != nil {
		log.Fatalln("netsh return a non-zero value while setting ipv4 rules:", string(out))
	}

	forwarder, err = forward.Forward("0.0.0.0:"+strconv.FormatUint(uint64(local_port), 10), "127.0.0.1:"+strconv.FormatUint(uint64(external_port), 10), forward.DefaultTimeout)
	if err != nil {
		panic(err)
	}
}
