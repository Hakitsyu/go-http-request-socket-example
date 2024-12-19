package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strings"
)

const URL string = "https://github.com/"

func main() {
	client, err := newHttpsClient(URL)
	panicErrorIfIsNotNull(err)

	host := getUrlHost(client.url)

	request := "GET / HTTP/1.1\r\n" +
		fmt.Sprintf("Host: %s\r\n", host) +
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7\r\n" +
		"Accept-Encoding: gzip, deflate, br\r\n" +
		"Accept-Language: pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7\r\n" +
		"Connection: closed\r\n" +
		"\r\n"

	buffer, err := client.Request(request)
	panicErrorIfIsNotNull(err)

	for {
		line, err := buffer.ReadString('\n')

		if err != nil {
			break
		}

		fmt.Print(line)
	}
}

type DefaultHttpClient interface {
	Request(content string) (*bufio.Reader, error)
}

func NewHttpClient(u string) (DefaultHttpClient, error) {
	pu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if pu.Scheme == "https" {
		client, err := newHttpsClient(u)
		if err != nil {
			return nil, err
		}

		return client, nil
	} else {
		client, err := newHttpClient(u)
		if err != nil {
			return nil, err
		}

		return client, nil
	}
}

type HttpClient struct {
	url *url.URL
}

func newHttpClient(u string) (*HttpClient, error) {
	pu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return &HttpClient{pu}, nil
}

func (c *HttpClient) Request(content string) (*bufio.Reader, error) {
	host := getUrlHost(c.url)
	const port int = 80

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(conn), nil
}

type HttpsClient struct {
	url *url.URL
}

func newHttpsClient(u string) (*HttpsClient, error) {
	pu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return &HttpsClient{pu}, nil
}

func (c *HttpsClient) Request(content string) (*bufio.Reader, error) {
	host := getUrlHost(c.url)
	const port int = 443

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(conn), nil
}

func getUrlHost(u *url.URL) string {
	return strings.TrimPrefix(u.Hostname(), "www.")
}

func getUrlDefaultPort(url *url.URL) int {
	if url.Scheme == "https" {
		return 443
	}

	return 80
}

func panicErrorIfIsNotNull(err error) {
	if err != nil {
		panic(err)
	}
}
