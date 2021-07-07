package server

import (
	"6-fyne-chat/controller"
	"fyne.io/fyne/widget"
	"log"
	"sort"
	"strings"
	"time"
)

func Updatechatroom(Message, UserList *widget.Label) {
	for {
		select {
		case m := <-controller.ClientMessage:
			//fmt.Println("rec msg:", m)
			var str strings.Builder
			t := time.Now().Format("2006-01-02 15:04:05")
			str.WriteString(string(t))
			str.WriteString("  ")
			str.WriteString("\n")
			str.WriteString(m.Msg)
			str.WriteString("\n")
			//fmt.Printf("%p",&Message)
			Message.Text += str.String()
			Message.Refresh()
			log.Println("message:###", m)
		case m := <-controller.Userlist:
			h := strings.Split(m.Msg, ",")
			sort.Strings(h)
			var s string
			for _, v := range h {
				s = s + v + ","
			}
			UserList.Text = s[1:]
			UserList.Refresh()
		default:
		}
	}
}
