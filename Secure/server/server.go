package server

import (
	"bufio"
	"crypto/ecdh"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pion/webrtc/v3"
)

type Server struct {
	Api    *webrtc.API
	Config *webrtc.Configuration
	Route  *mux.Route
}

type Room struct {
	Curve         *ecdh.Curve
	ListPublicKey []*ecdh.PublicKey
	Conns         *webrtc.PeerConnection
	DataChannel   *webrtc.DataChannel
}

func NewServer(api *webrtc.API, config *webrtc.Configuration) *Server {

	return &Server{
		Api:    api,
		Config: config,
		Route:  r,
	}
}

func (server *Server) NewRoom(curve *ecdh.Curve, config *webrtc.Configuration) (*Room, error) {

	peerConn, err := server.Api.NewPeerConnection(*server.Config)
	if err != nil {
		return nil, err
	}
	dataChannel, err := peerConn.CreateDataChannel("chat", nil)
	if err != nil {
		return nil, err
	}
	return &Room{
		Curve:         curve,
		ListPublicKey: []*ecdh.PublicKey{},
		Conns:         peerConn,
		DataChannel:   dataChannel,
	}, nil
}

func (room *Room) JoinRoom(remote *ecdh.PublicKey, conn net.Conn) {
	room.ListPublicKey = append(room.ListPublicKey, remote)
}

func (room *Room) GetRemotePublicKey(publicKey *ecdh.PublicKey) (*ecdh.PublicKey, error) {
	for _, value := range room.ListPublicKey {
		if !publicKey.Equal(value) {
			return value, nil
		}
	}
	return publicKey, nil
	// return nil, errors.New("not found remote client")
}

func (room *Room) HandleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for {
		if scanner.Scan() {
			msg := scanner.Text()
			fmt.Println("Server received message: ", msg)
			//for _, con := range room.Conns {
			fmt.Fprintf(conn, msg+"\n")
			//}
		} else {
			break
		}
	}
}
