package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"example/client"
	"example/server"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/pion/webrtc/v3"
)

func main() {
	r := mux.NewRouter()
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err.Error())
	}

	// Create clientA
	curve := ecdh.P256()
	// Create a new WebRTC API
	api := webrtc.NewAPI()
	// Create a new peer connection configuration
	config := webrtc.Configuration{}
	server := server.NewServer(api, &config)
	if len(os.Args) < 2 {
		fmt.Println("Please specify server or client")
		os.Exit(1)
	}
	if os.Args[1] == "server" {
		ln, err := net.Listen("tcp", ":8080")
		room, err := server.NewRoom(&curve, server.Config)
		if err != nil {
			fmt.Println("Error starting server:", err.Error())
			os.Exit(1)
		}
		defer ln.Close()

		for {
			conn, err := ln.Accept()
			//room.Conns = append(room.Conns, conn)

			if err != nil {
				fmt.Println("Error accepting connection:", err.Error())
				continue
			}
			go room.HandleConnection(conn)
		}
	} else if os.Args[1] == "client" {
		client, err := client.NewClient(curve, &rand.Reader, protocol, host)
		if err != nil {
			log.Fatalf("Error create client: %v", err)

		}
		conn, err := client.ConnectToServer()
		//room.JoinRoom(client.PublicKey, conn)
		if err != nil {
			log.Fatalf("Error create client: %v", err)

		}
		defer conn.Close()

		go client.SendMsg(conn, room)
		go client.RecvMsg(conn, room)

		select {}
	} else {
		fmt.Println("Invalid argument")
		os.Exit(1)
	}

}
