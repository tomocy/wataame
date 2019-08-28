package http

import (
	"errors"
	"fmt"
	"strings"
)

type Request struct {
	URI string
}

func compensateAddress(uri string) (string, error) {
	host, port, err := splitHostPort(uri)
	if err != nil {
		return "", err
	}
	if port == "" {
		port = "80"
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}

func splitHostPort(uri string) (string, string, error) {
	splited := strings.Split(uri, ":")
	if len(splited) > 2 {
		return "", "", errors.New("invalid format of uri: the format should be host[:port]")
	}
	var host, port string
	host = splited[0]
	if len(splited) == 2 {
		port = splited[1]
	}

	return host, port, nil
}

type Response []byte
