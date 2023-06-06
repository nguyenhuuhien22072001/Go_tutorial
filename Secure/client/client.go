package client

import (
	"bufio"
	"crypto/ecdh"
	"example/encrypt"
	"example/server"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	Protocol   string
	Host       string
	PrivateKey *ecdh.PrivateKey
	PublicKey  *ecdh.PublicKey
}

func NewClient(curve ecdh.Curve, key *io.Reader, protocol, host string) (*Client, error) {
	privateKey, err := curve.GenerateKey(*key)
	if err != nil {
		return nil, err
	}

	return &Client{
		Protocol:   protocol,
		Host:       host,
		PrivateKey: privateKey,
		PublicKey:  privateKey.PublicKey(),
	}, nil
}

func (client *Client) GenerateExchangeKey(remote *ecdh.PublicKey) ([]byte, error) {
	return client.PrivateKey.ECDH(remote)
}

func (client *Client) SendMsg(conn net.Conn, room *server.Room) {
	remotePublicKey, err := room.GetRemotePublicKey(client.PublicKey)
	if err != nil {
		fmt.Fprint(conn, err.Error()+"\n")
		return
	}
	keyE, err := client.GenerateExchangeKey(remotePublicKey)
	if err != nil {
		fmt.Fprint(conn, err.Error()+"\n")
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		scanner.Scan()
		msg := scanner.Text()
		msgE := encrypt.Encrypt(msg, string(keyE))
		fmt.Fprint(conn, msgE+"\n")
	}
	// messageE := Encrypt(message, string(keyE))
	// fmt.Println(messageE)
	// return nil
}

func (client *Client) ConnectToServer() (net.Conn, error) {
	conn, err := net.Dial(client.Protocol, client.Host)
	if err != nil {
		fmt.Println("Error connecting to server:", err.Error())
		return nil, err
	}
	return conn, nil
}

func (client *Client) RecvMsg(conn net.Conn, room *server.Room) {
	scanner := bufio.NewScanner(conn)
	for {
		if scanner.Scan() {
			msgE := scanner.Text()
			remotePublicKey, err := room.GetRemotePublicKey(client.PublicKey)
			if err != nil {
				fmt.Fprintf(conn, err.Error()+"\n")
				return
			}
			keyE, err := client.GenerateExchangeKey(remotePublicKey)
			if err != nil {
				fmt.Fprintf(conn, err.Error()+"\n")
				return
			}
			msg := encrypt.Encrypt(msgE, string(keyE))
			fmt.Println("Received message:", msg)
		} else {
			break
		}
	}
}
