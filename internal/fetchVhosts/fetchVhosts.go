package fetchVhosts

import (
	"bytes"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

var serverMap = map[string]func(*ssh.Session) []string{
	"nginx": getNginxVhosts,
}

func FetchVhosts(host, httpServer string, config *ssh.ClientConfig) []string {
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalln(err)
	}
	defer session.Close()

	vhosts := make([]string, 0)
	if _, ok := serverMap[httpServer]; ok {
		vhosts = serverMap[httpServer](session)
	}

	return vhosts
}

func getNginxVhosts(session *ssh.Session) []string {
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run("sudo -n nginx -T | grep server_name")
	return strings.Split(stdoutBuf.String(), "\n")
}
