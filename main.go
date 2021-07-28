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
		wg.Add(1)
		go ping(ip, &wg)
	}

	wg.Wait()
	// time.Sleep(1 * time.Hour)
}

func getIPs() []string {
	cmd := exec.Command("arp", "-a")
	out, err := cmd.Output()
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
	RunCommand(wg, "ping", ip, "-i 0.2", "-s 300", "-c 100")
}

func RunCommand(wg *sync.WaitGroup, name string, arg ...string) error {
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
		wg.Done()
		return err
	}

	wg.Done()
	return nil
}
