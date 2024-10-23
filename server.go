package gows

import (
	"fmt"
	"net"
)

// Configuration of the WebSocket Server
type ServerConfig struct {

	// Addr used by the Server
	Addr string
}

type WebSocketServer struct {
	config ServerConfig
}

// Create a New Server with given configuration and return the reference
func NewWebSocketServer(cfg ServerConfig) *WebSocketServer {

	return &WebSocketServer{
		config: cfg,
	}

}

// Starts the Server and Handles the Incoming Connections
func (S *WebSocketServer) Start() error {

	_listener, err := net.Listen("tcp", S.config.Addr)

	if err != nil {
		return err
	}

	fmt.Printf("Server listening at localhost%s\n", S.config.Addr)

	for {
		conn, err := _listener.Accept()
		if err != nil {
			fmt.Println("Error in getting the connection:", err)
			continue
		}

		fmt.Printf("Got a %s connection with addr:%s\n", conn.RemoteAddr().Network(), conn.RemoteAddr().String())

		go S.handleConn(conn.(*net.TCPConn))
	}

}

// Handles the incoming Connections
func (S *WebSocketServer) handleConn(conn *net.TCPConn) {
	err := handleHandShake(conn)
	if err != nil {
		fmt.Println("Error in handleHandShake:", err)
		conn.Write(invalidHandshake())
		conn.Close()
		return
	}
}
