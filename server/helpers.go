package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func getPort(port string) (string, error) {
	_port, err := strconv.Atoi(port)
	if err != nil {
		return "", err
	}

	if _port <= 0 || _port > 65535 {
		return "", errors.New("invalid port number: " + port)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", _port))
	for err != nil {
		log.Printf("Port %d is unavailable, trying next one.\n", _port)
		_port += 1

		if _port > 65535 {
			return "", errors.New("maximum port number reached - none available")
		}

		ls, err = net.Listen("tcp", fmt.Sprintf(":%d", _port))
	}
	ls.Close()

	return fmt.Sprintf("%d", _port), nil
}

func getIPs() (list []string, loopback string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return list, "", err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if ip == "127.0.0.1" {
					loopback = ip
				} else {
					list = append(list, ip)
				}
			}
		}
	}

	return
}

func validateDir(dir string) (string, error) {
	var err error

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	fi, err := os.Stat(dir)
	if err != nil {
		return "", err
	}

	if !fi.IsDir() {
		msg := fmt.Sprintf("Not a directory: %s", dir)
		return "", errors.New(msg)
	}

	return dir, nil
}

func generateCerts(path string, force bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, 0700); err != nil {
			return err
		}
	}

	match, err := filepath.Glob(path + "/server.*")
	if err != nil {
		log.Fatalln(err)
	}

	if !force && len(match) == 2 {
		log.Println("Using existing certificates from", path)
		return nil
	}

	openssl, err := exec.LookPath("openssl")
	if err != nil {
		return err
	}

	keyParams := []string{"genrsa", "-out", path + "/server.key", "2048"}
	certParams := []string{"req", "-batch", "-new", "-x509", "-sha256", "-key", path + "/server.key", "-out", path + "/server.crt", "-days", "365"}

	log.Println("Generating server.key")
	o, err := exec.Command(openssl, keyParams...).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(o))

	log.Println("Generating server.crt")
	_, err = exec.Command(openssl, certParams...).CombinedOutput()
	if err != nil {
		return err
	}

	log.Println("Certificate generation OK")

	return nil
}
