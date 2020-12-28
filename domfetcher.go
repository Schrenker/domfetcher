package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/ssh"
)

func checkHTTPServer(host string) ([]string, error) {
	url := strings.Split(host, ":")[0]

	res, err := http.Head("http://" + url)
	if err != nil {
		return nil, err
	}

	return res.Header["Server"], nil
}

func getVhosts(host string, config *ssh.ClientConfig) {
	HTTPServer, err := checkHTTPServer(host)
	if err != nil {
		fmt.Println(err)
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
	getVhosts(hosts[0], config)
}

func main() {
	fetchFromFile("centos", "private/id_rsa", "private/input.txt")
}
