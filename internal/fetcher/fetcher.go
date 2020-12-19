package fetcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

func parseSSHAuth(user, path string) *ssh.ClientConfig {
	key, err := ioutil.ReadFile(path)
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

func loadInputFile(path string) []string {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Split(string(file), "\n")
}

//FetchFromFile ...
func FetchFromFile(path string) {
	config := parseSSHAuth("centos", "private/id_rsa")
	hosts := loadInputFile(path)
	client, err := ssh.Dial("tcp", hosts[0], config)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	ss, err := client.NewSession()
	if err != nil {
		log.Fatalln(err)
	}
	defer ss.Close()

	var stdoutBuf bytes.Buffer
	ss.Stdout = &stdoutBuf
	ss.Run("uname -a")
	fmt.Println(stdoutBuf.String())
}
