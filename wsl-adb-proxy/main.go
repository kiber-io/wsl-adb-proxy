package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

const (
	WSL_BIND_HOST = "127.0.0.1"
	WSL_BIND_PORT = 5037
	WIN_ADB_PORT  = 5037
)

func getWinHostIP() (string, error) {
	cmd := exec.Command("ip", "route", "show")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get host IP address: %w", err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), "default") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return parts[2], nil
			}
		}
	}
	return "", fmt.Errorf("failed to get host IP address")
}

func forward(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			_, err2 := dst.Write(buf[:n])
			if err2 != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
}

func handleClient(client net.Conn, winHost string) {
	remote, err := net.Dial("tcp", fmt.Sprintf("%s:%d", winHost, WIN_ADB_PORT))
	if err != nil {
		client.Close()
		return
	}
	go forward(client, remote)
	go forward(remote, client)
}

func startProxy() error {
	winHost, err := getWinHostIP()
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", WSL_BIND_HOST, WSL_BIND_PORT))
	if err != nil {
		return err
	}
	fmt.Printf("ADB proxy listening on %s:%d\n", WSL_BIND_HOST, WSL_BIND_PORT)
	for {
		client, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(client, winHost)
	}
}

func main() {
	if err := startProxy(); err != nil {
		fmt.Println("Error:", err)
	}
}
