# Go-Ws

Experimental WebSocket Server for Go. Use it on your own risk.  
For now, can recieve and write messages back to websocket. 

Create a Websocket Server and start it.    
```go
server := gows.NewWebSocketServer(gows.ServerConfig{
    Addr: ":5000",
})
server.Start()
```

Accept the incoming connection  
```go
wsconn,_ := server.Accept()
```

Connections will be accepted as FIFO. It's user responsibility to handle individual connections.  

Read from the accepted connection  
```go
message,_ := wsconn.Read()
fmt.Println(string(message))
```

Write to the connection  
```go
wsconn.Write([]byte("Hello World"))
```

**Notes**  

1. For now, It dont handle opcode specific actions  
2. No Extensions  
3. Message frames should follow the rules established for the protocol.  
