package main

import (
	"log"
	"os/exec"
	"strconv"
)

func clear_rule_v4() {

	if out, err := exec.Command("bash", "-c", "iptables-save --counters | grep -v TRAVERSAL | iptables-restore --counters --table nat").CombinedOutput(); err != nil {
		log.Fatalln("iptables-save return a non-zero value while clearing ipv4 rules:", string(out))
	}
}

func set_rule_v4() {

	if out, err := exec.Command("bash", "-c", `iptables-restore --noflush <<-EOF
		*nat
		:TRAVERSAL -
		-A PREROUTING -m addrtype --dst-type LOCAL -j TRAVERSAL	
		COMMIT
		EOF`).CombinedOutput(); err != nil {
		log.Fatalln("iptablesi-restore return a non-zero value while setting ipv4 rules:", string(out))
	}
}

func modify_rule_v4(external_port uint16) {

	local_port := get_conf("local_port").(uint16)

	if out, err := exec.Command("bash", "-c", `iptables-restore --noflush <<-EOF
		*nat
		-F TRAVERSAL
		-I TRAVERSAL -p tcp -m tcp --dport `+strconv.FormatUint(uint64(local_port), 10)+` -j REDIRECT --to-ports `+strconv.FormatUint(uint64(external_port), 10)+`
		-I TRAVERSAL -p udp -m udp --dport `+strconv.FormatUint(uint64(local_port), 10)+` -j REDIRECT --to-ports `+strconv.FormatUint(uint64(external_port), 10)+`
		COMMIT
		EOF`).CombinedOutput(); err != nil {
		log.Fatalln("iptablesi-restore return a non-zero value while modifying ipv4 rules:", string(out))
	}
}
