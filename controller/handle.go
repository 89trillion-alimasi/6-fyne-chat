package controller

import (
	Protobuf "6-fyne-chat/model"
	"6-fyne-chat/view"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

var (
	Conn *websocket.Conn

	receivermessage = make(chan *Protobuf.Communication)
	sendmessage     = make(chan *Protobuf.Communication)
	Userlist        = make(chan *Protobuf.ClientMessage, 1)
	ClientMessage   = make(chan *Protobuf.ClientMessage, 1)
)

func GetConnection(url, username string) error {
	header := http.Header{}
	header["login-username"] = []string{username}
	c, _, err := websocket.DefaultDialer.Dial(url, header)
	Conn = c

	go receivermessageclient()
	go serverclient()

	return err
}

//发送断开连接消息
func SendExit() error {
	//fmt.Println("已经退出去了-----")
	sendmessage <- Protobuf.GetExitMessage()
	err := Conn.Close()
	if err != nil {
		return err
	}
	Conn = nil
	return nil
}

func Sendmessage(message string) {
	bytes := Protobuf.GettalkMessage(message)
	//fmt.Println("send message :", message)

	sendmessage <- bytes
	bytes = Protobuf.Getuserlist()
	sendmessage <- bytes
	log.Printf("SendMessage success, message is [%s]\n", message)
}

//读消息
func receivermessageclient() {
	for {
		_, message, err := Conn.ReadMessage()
		if err != nil {
			log.Println("receivermessageMessage error:", err)
			return
		}
		var m Protobuf.Communication
		if err := proto.Unmarshal(message, &m); err != nil {
			log.Println("Unmarshal error:", err)
			return
		}
		receivermessage <- &m
		//log.Printf("receivermessageMessage: %s", &m)
	}
}

func serverclient() {
	for {
		select {
		// send
		case m := <-sendmessage:
			//fmt.Println("sendmessager:", m)
			bytes, err := proto.Marshal(m)
			if err != nil {
				log.Println("Marshal error:", err)
				return
			}
			if err := Conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
				log.Println("sendmessageMessage error:", err)
				return
			}
		// rec
		case m := <-receivermessage:

			switch m.Type {
			case "1":
				msg := &Protobuf.ClientMessage{
					Usrname: m.Username,
					Msg:     m.Msg,
				}
				ClientMessage <- msg
			case "3":
				msg := &Protobuf.ClientMessage{
					Usrname: m.Username,
					Msg:     m.Msg,
				}
				Userlist <- msg
			default:
			}
		}
	}
}

func Updatechatroom() {
	for {
		select {
		case m := <-ClientMessage:
			//fmt.Println("rec msg:", m)
			var str strings.Builder
			t := time.Now().Format("2006-01-02 15:04:05")
			str.WriteString(string(t))
			str.WriteString("  ")
			str.WriteString("\n")
			str.WriteString(m.Msg)
			str.WriteString("\n")
			view.Message.Text += str.String()
			view.Message.Refresh()
			log.Println("message:###", m)
		case m := <-Userlist:
			h := strings.Split(m.Msg, ",")
			sort.Strings(h)
			var s string
			for _, v := range h {
				s = s + v + ","
			}
			view.UserList.Text = s[1:]
			view.UserList.Refresh()
		default:
		}
	}
}
