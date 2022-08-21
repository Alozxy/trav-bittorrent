package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func transmission_port(external_port uint16) {
	client := &http.Client{}

	status_code := rpc_request(external_port, client)
	if status_code == 409 {
		rpc_request(external_port, client)
	}
}

func rpc_request(external_port uint16, client *http.Client) int {

	username := conf.get_conf("username").(string)
	password := conf.get_conf("password").(string)
	rpc_url := "http://" + conf.get_conf("address").(string) + "/transmission/rpc"

	data := rpc_body{
		Method: "session-set",
		Arguments: arg_obj{
			Peer_port: external_port,
		},
		Tag: "",
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return 1
	}

	req, err := http.NewRequest("POST", rpc_url, strings.NewReader(string(b)))
	if err != nil {
		log.Println(err)
		return 1
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("X-Transmission-Session-Id", get_conf("X-Transmission-Session-Id").(string))

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 1
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return 1
	}

	if res.StatusCode == 409 {
		set_conf("X-Transmission-Session-Id", res.Header.Get("X-Transmission-Session-Id"))
		return 409
	}

	if res.StatusCode != 200 {
		log.Println("set config error:", string(body))
		return 1
	}

	log.Println("set config successfully")
	return 0
}

type rpc_body struct {
	Method    string  `json:"method"`
	Arguments arg_obj `json:"arguments"`
	Tag       string  `json:"tag"`
}

type arg_obj struct {
	Peer_port uint16 `json:"peer-port"`
}
