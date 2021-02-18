package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/schrenker/domfetcher/internal/fetchhttpserver"
	"github.com/schrenker/domfetcher/internal/fetchvhosts"
	"golang.org/x/crypto/ssh"
)

func getVhosts(host string, config *ssh.ClientConfig) {
	HTTPServer, err := fetchhttpserver.FetchHTTPServer(host)
	if err != nil {
		return
	}
	vhosts := fetchvhosts.FetchVhosts(host, HTTPServer, config)
	for i := range vhosts {
		fmt.Println(vhosts[i])
	}
	//add vhost fetching here
}

func parseSSHAuth(user, keyPath string) *ssh.ClientConfig {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalln(err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalln(err)
	}

	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func loadInputFile(inputPath string) []string {
	file, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Fatalln(err)
	}
	list := strings.Split(string(file), "\n")
	result := make([]string, 0)
	for i := range list {
		if list[i][0] != '#' {
			result = append(result, list[i])
		}
	}
	return result
}

func fetchFromFile(user, keyPath, inputPath string) {
	hosts := loadInputFile(inputPath)
	config := parseSSHAuth(user, keyPath)
	for i := range hosts {
		getVhosts(hosts[i], config)
	}
}

func main() {
	fetchFromFile("kylos", "private/kbkey", "private/input.txt")
}
