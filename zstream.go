package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
)

var nalPrefix []byte
var startCMD []byte

var camera *net.TCPConn
var client *net.TCPConn
var listen *net.TCPListener

func nalSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	idx := bytes.Index(data, nalPrefix)
	if idx == -1 {
		return 0, nil, nil
	}

	return idx + 3, data[0:idx], nil
}

func connectCamera(addr string) (*bufio.Scanner, error) {
	if camera != nil {
		camera.Close()
		camera = nil
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	camera, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	_, err = camera.Write(startCMD)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(camera)
	scanner.Split(nalSplit)

	return scanner, nil
}

func listenForClient(addr string) error {
	if listen != nil {
		listen.Close()
		listen = nil
	}

	if client != nil {
		client.Close()
		client = nil
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	listen, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	client, err = listen.AcceptTCP()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	nalPrefix = []byte{0, 0, 1}
	startCMD = []byte{0x55, 0x55, 0xaa, 0xaa,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x50}
}

func main() {
	listenAddr := flag.String("l", ":8888", "Listen address")
	cameraAddr := flag.String("c", "", "Camera address (address:port)")
	flag.Parse()

	if *cameraAddr == "" {
		fmt.Printf("Camera address required\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	for {
		if listen != nil {
			listen.Close()
			listen = nil
		}

		if client != nil {
			client.Close()
			client = nil
		}

		if camera != nil {
			camera.Close()
			camera = nil
		}

		if listenForClient(*listenAddr) != nil {
			continue
		}

		scanner, err := connectCamera(*cameraAddr)
		if err != nil {
			continue
		}

		started := false
		for scanner.Scan() {
			nal := scanner.Bytes()
			if !started {
				if nal[0] == 0x67 {
					started = true
				} else {
					continue
				}
			}
			_, err = client.Write(nalPrefix)
			if err != nil {
				break
			}

			_, err = client.Write(nal)
			if err != nil {
				break
			}
		}
	}
}
