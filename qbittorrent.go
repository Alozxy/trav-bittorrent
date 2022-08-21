package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

func qbittorrent_port(external_port uint16) {

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	login(client)
	set_preferences(external_port, client)
}

func login(client *http.Client) {

	addr := conf.get_conf("address").(string)
	username := conf.get_conf("username").(string)
	password := conf.get_conf("password").(string)
	schema := "http://" + addr

	res, err := client.PostForm(schema+"/api/v2/auth/login", url.Values{"username": {username}, "password": {password}})
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if res.StatusCode != 200 {
		log.Println("login error:", string(body))
		return
	}

	log.Println("login successfully")
}

func set_preferences(external_port uint16, client *http.Client) {

	addr := conf.get_conf("address").(string)
	schema := "http://" + addr

	json := "{ \"listen_port\": " + strconv.FormatUint(uint64(external_port), 10) + " }"
	res, err := client.PostForm(schema+"/api/v2/app/setPreferences", url.Values{"json": {json}})
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if res.StatusCode != 200 {
		log.Println("set preferences error:", string(body))
		return
	}

	log.Println("set preferences successfully")

}
