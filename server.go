package gows

import (
	"bytes"
	"fmt"
	"net"
)

type WsConn struct {
	conn      *net.TCPConn
	ouputChan chan MessagePacket
}

type MessagePacket struct {
	message []byte
	err     error
}

type AcceptPacket struct {
	conn *WsConn
	err  error
}

// Configuration of the WebSocket Server
type ServerConfig struct {

	// Addr used by the Server
	Addr string
}

type WebSocketServer struct {
	config ServerConfig

	acceptChan chan AcceptPacket
}

// Create a New Server with given configuration and return the reference
func NewWebSocketServer(cfg ServerConfig) *WebSocketServer {

	return &WebSocketServer{
		config:     cfg,
		acceptChan: make(chan AcceptPacket, 1),
	}

}

// Starts the Server and Handles the Incoming Connections
func (S *WebSocketServer) Start() error {

	_listener, err := net.Listen("tcp", S.config.Addr)

	if err != nil {
		return err
	}

	fmt.Printf("Server listening at localhost%s\n", S.config.Addr)
	go accpetConn(S, _listener)
	return nil

}

func (S *WebSocketServer) Accept() (*WsConn, error) {
	acpakcet := <-S.acceptChan
	return acpakcet.conn, acpakcet.err
}

func (ws *WsConn) Read() ([]byte, error) {
	msgpacket := <-ws.ouputChan
	return msgpacket.message, msgpacket.err
}

func (ws *WsConn) Write(data []byte) error {

	message, err := prepareMessage(data)
	if err != nil {
		return fmt.Errorf("preparingMessage: %s", err)
	}

	_, err = ws.conn.Write(message)
	if err != nil {
		return fmt.Errorf("writing to socket: %s", err)
	}
	return nil

}

func accpetConn(S *WebSocketServer, _listener net.Listener) {
	for {
		rawConn, err := _listener.Accept()
		if err != nil {
			fmt.Println("Error in getting the connection:", err)
			S.acceptChan <- AcceptPacket{
				conn: nil,
				err:  err,
			}
			continue
		}

		tcpConn := rawConn.(*net.TCPConn)

		fmt.Printf("Got a %s connection with addr:%s\n", tcpConn.RemoteAddr().Network(), tcpConn.RemoteAddr().String())

		err = handleHandShake(tcpConn)
		if err != nil {
			fmt.Printf("Error in handleHandShake: %s\n", err)
			fmt.Printf("Closing the %s Connection at %s\n", tcpConn.RemoteAddr().Network(), tcpConn.RemoteAddr().String())
			S.acceptChan <- AcceptPacket{
				conn: nil,
				err:  err,
			}
			tcpConn.Close()
		}

		wsconn := WsConn{
			conn:      tcpConn,
			ouputChan: make(chan MessagePacket, 1),
		}

		S.acceptChan <- AcceptPacket{
			conn: &wsconn,
			err:  nil,
		}

		go handleConn(&wsconn)
	}
}

// Handles the incoming Connections
func handleConn(wsconn *WsConn) {

	tcpConn := wsconn.conn

	buffer := &bytes.Buffer{}
	buff_size := 2

	for {
		buff := make([]byte, buff_size)
		n, err := tcpConn.Read(buff)
		if err != nil {
			fmt.Printf("Error in Reading Message: %s\n", err)
			return
		}
		buffer.Write(buff[:n])

		if n < buff_size {
			// Got the Complete message

			message, err := handleMessage(buffer.Bytes())
			if err != nil {
				fmt.Printf("Error in handling Mesasge: %s\n", err)
			}

			wsconn.ouputChan <- MessagePacket{
				message: message,
				err:     err,
			}

			buffer.Reset()

		}

	}
}
