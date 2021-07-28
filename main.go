package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

func main() {
	ips := getIPs()

	var wg sync.WaitGroup
	for _, ip := range ips {
		go ping(ip, &wg)
	}

	wg.Wait()
}

func getIPs() []string {
	out, err := exec.Command("arp", "-a").Output()
	if err != nil {
		log.Fatal(err)
	}

	stringifiedOutput := string(out)

	var ips []string

	regexp := regexp.MustCompile(`\(([^)]+)\)`)

	results := strings.Split(stringifiedOutput, "\n")
	for _, line := range results {
		match := regexp.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}
		ips = append(ips, match[1])
	}

	return ips
}

func ping(ip string, wg *sync.WaitGroup) {
	RunCommand("ping", ip, "-i 0.2", "-s 300", "-c 100")
	wg.Add(1)
}

func RunCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		if strings.Contains(string(tmp), "packet loss") {
			fmt.Print("\n", string(tmp))
		}
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}
