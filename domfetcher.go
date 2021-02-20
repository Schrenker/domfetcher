package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/schrenker/domfetcher/internal/fetchhttpserver"
	"github.com/schrenker/domfetcher/internal/fetchvhosts"
	"golang.org/x/crypto/ssh"
)

func getVhosts(host string, config *ssh.ClientConfig) ([]string, error) {
	HTTPServer, err := fetchhttpserver.FetchHTTPServer(host)
	if err != nil {
		return nil, err
	}
	return fetchvhosts.FetchVhosts(host, HTTPServer, config), nil
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

func fetchFromFile(user, keyPath, inputPath string) map[string][]string {
	hosts := loadInputFile(inputPath)
	config := parseSSHAuth(user, keyPath)
	hostMap := make(map[string][]string)
	for i := range hosts {
		vhosts, err := getVhosts(hosts[i], config)
		if err != nil {
			log.Fatalf("Couldn't fetch vhosts for %v. Err: %v\n", hosts[i], err)
		}
		hostMap[hosts[i]] = vhosts
	}

	// for k, v := range hostMap {
	// 	fmt.Println(k)
	// 	for j := range v {
	// 		fmt.Println(v[j])
	// 	}
	// }

	return hostMap
}

func main() {
	user := flag.String("u", "", "SSH login user")
	identity := flag.String("i", "", "passwordless identity file")
	hostsFile := flag.String("f", "", "hosts file")

	flag.Parse()

	if *user == "" || *hostsFile == "" || *identity == "" {
		log.Fatalln("User, hosts file and identity file required for this to work")
	}
	hostMap := fetchFromFile(*user, *identity, *hostsFile)
	// jsonString, err := json.Marshal(hostMap)
	jsonString, err := json.MarshalIndent(hostMap, "", "")
	if err != nil {
		log.Fatalln(err)
	}

	_ = ioutil.WriteFile("output.json", jsonString, 0644)
}
