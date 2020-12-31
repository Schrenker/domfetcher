package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/schrenker/domfetcher/internal/fetchHTTPServer"
	"golang.org/x/crypto/ssh"
)

func getVhosts(host string, config *ssh.ClientConfig) {
	HTTPServer, err := fetchHTTPServer.FetchHTTPServer(host)
	if err != nil {
		return
	}
	fmt.Println(HTTPServer)

	// client, err := ssh.Dial("tcp", hosts[0], config)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer client.Close()

	// ss, err := client.NewSession()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer ss.Close()

	// var stdoutBuf bytes.Buffer
	// ss.Stdout = &stdoutBuf
	// ss.Run("uname -a")
	// fmt.Println(stdoutBuf.String())

}

// remember to add comments option
func loadInputFile(inputPath string) []string {
	file, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Split(string(file), "\n")
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
