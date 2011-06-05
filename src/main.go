package main

import (
	"http"
	"log"
	"fmt"
	"socketio"
	"json"
	"bytes"
	"redis"
	"strings"
)

const (
	HttpListen = "0.0.0.0:8080"
	RedisAddr = "127.0.0.1:6379"
)

var SocketIO socketio.SocketIO
var Sockets map[string] *socketio.Conn
var Users map[string] *socketio.Conn
var RUsers map[*socketio.Conn] string

func initSocketIO() {
    log.Println("init socketio server")
    config := socketio.DefaultConfig
	config.Origins = []string{"127.0.0.1:3000"}
	SocketIO := socketio.NewSocketIO(&config)

    go func() {
		if err := SocketIO.ListenAndServeFlashPolicy(":843"); err != nil {
			log.Println(err)
		}
	}()

	SocketIO.OnConnect(func(c *socketio.Conn) {
        Sockets[c.String()] = c
        c.Send(Message{[]string{"Server", "hello"}})
    })

    SocketIO.OnDisconnect(func(c *socketio.Conn) {
        //log.Println(c.String(), " disconnetced")
        hash := c.String()
        Sockets[hash] = nil, false
        uid, exists := RUsers[c]
        if exists {
            RUsers[c] = "", false
            Users[uid] = nil, false
        }
    })

    SocketIO.OnMessage(func(c *socketio.Conn, msg socketio.Message) {
        var info = Info{};
        v, _ := msg.JSON()
        b := bytes.NewBufferString(v)
        err := json.Unmarshal(b.Bytes(), &info)
        if err != nil {
            log.Println("Client: ", c.String(), " send a invalid message: ", msg.Data())
            return
        }
        if info.Hash != "" {
            log.Println("Link Socket: ", c.String(), " to: ", info.Hash)
            Users[info.Hash] = c
            RUsers[c] = info.Hash
        }
    })

    mux := SocketIO.ServeMux()
	err := http.ListenAndServe(HttpListen, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}
}

func initPubSub() {
    log.Println("init redis client")
    var client redis.Client
    client.Addr = RedisAddr
    sub := make(chan string, 1)
    sub <- "user:*:*"
    messages := make(chan redis.Message, 0)
    go client.Subscribe(nil, nil, sub, nil, messages)
    go func(){
        for {
            msg := <-messages
            subs := strings.Split(msg.Channel, ":", 3)
            uid := subs[1]
            s, exists := Users[uid]
            if exists {
                var pack = Message{};
                err := json.Unmarshal(msg.Message, &pack)
                if err != nil {
                    log.Println("Unmarshal msg err: ", string(msg.Message))
                    continue
                }
                pack.Message = append(pack.Message, subs[2])
                s.Send(pack)
                log.Println("Send to: ", uid, " Data: ", pack)
            }

        }
    }()
}

func startConsole(cmds chan<- string) {
	var input string
	for {
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Println("Input Error: ", err.String())
			continue
		}

		if input == "exit" {
			close(cmds)
			log.Println("Server exit.")
			break
		} else {
		}
		cmds <- input
	}
}

func processCMD(cmds <-chan string) {
	for {
		if v, closed := <-cmds; closed != false {
			log.Println("input: ", v)
		} else {
			break
		}
	}
}


func main() {
     Sockets = make(map[string] *socketio.Conn)
    Users = make(map[string] *socketio.Conn)
    RUsers = make(map[*socketio.Conn] string)
    cmds := make(chan string)
	go processCMD(cmds)
	go initSocketIO()
    go initPubSub()
	startConsole(cmds)

}

