package view

import (
	"6-fyne-chat/controller"
	server2 "6-fyne-chat/server"
	"errors"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"image/color"
)

var username *widget.Entry //用户名
var server *widget.Entry   //服务器名字
var UserList *widget.Label //用户列表
var Message *widget.Label  //消息
var status *widget.Label   //状态

func Chat(w fyne.Window) *fyne.Container {
	op := container.NewGridWithRows(1, loginOperator(w), statusView(w))
	top := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(2), ConnView(w), op)
	line := canvas.NewLine(color.White)
	line.Resize(fyne.NewSize(100, 1))
	h3 := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), layout.NewSpacer(), line, layout.NewSpacer())
	//list:=fyne.NewContainerWithLayout(layout.NewGridLayoutWithRows(1),alluser(w)
	//fmt.Printf("%p,%p\n",&Message,&UserList)
	UserList = widget.NewLabel("")
	Message = widget.NewLabel("")

	go server2.Updatechatroom(Message, UserList)
	chat := fyne.NewContainerWithLayout(layout.NewGridLayoutWithRows(4), top, h3, alluser(w), chatingroom(w))

	return chat
}

func loginOperator(w fyne.Window) fyne.CanvasObject {
	form := &widget.Form{
		CancelText: "quit",
		SubmitText: "link",
		OnSubmit: func() {
			server := server.Text
			name := username.Text
			if controller.Conn != nil {
				dialog.NewError(errors.New("already connected"), w)
				return
			} else if len(server) == 0 || len(name) == 0 {
				dialog.NewError(errors.New("Can not be empty"), w)
				return
			}
			err := controller.GetConnection(server, name)
			if err != nil {
				dialog.NewError(err, w)
			} else {
				status.Text = "success"
				dialog.NewInformation("success", "you can talk", w)
			}
		},
		OnCancel: func() {
			if controller.Conn == nil {
				dialog.NewError(errors.New("connect first"), w)
				return
			}
			err := controller.SendExit()
			if err != nil {
				dialog.NewError(err, w)
			} else {
				dialog.NewError(errors.New("disconnected"), w)
			}
			status.Text = "ERROR"
		},
	}
	return form
}

//状态
func statusView(_ fyne.Window) fyne.CanvasObject {
	space := widget.NewLabel("")

	text := widget.NewLabel("Status:")
	status = widget.NewLabel("ERROR")

	return container.NewVBox(space, container.NewVBox(text, status))
}

func ConnView(_ fyne.Window) fyne.CanvasObject {
	username = widget.NewEntry()
	server = widget.NewEntry()
	server.SetText("ws://localhost:8080/ws")

	form := widget.NewForm(
		&widget.FormItem{Text: "Server:", Widget: server},
		&widget.FormItem{Text: "Username:", Widget: username},
	)

	return container.NewScroll(form)
}

func chatingroom(w fyne.Window) fyne.CanvasObject {
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("please input there")
	form := &widget.Form{
		SubmitText: "send",
		OnSubmit: func() {
			if controller.Conn == nil {
				dialog.NewError(errors.New("you should conn first"), w)
			} else if len(input.Text) == 0 {
				dialog.NewError(errors.New("input someting"), w)
			} else {
				controller.Sendmessage(input.Text)
				input.Text = ""
			}
		},
	}
	form.Append("", input)
	return form

}

func alluser(w fyne.Window) fyne.CanvasObject {

	hbox := container.NewHBox(Message)
	record := container.NewHScroll(hbox)

	a := widget.NewGroup("Online User", UserList)
	b := widget.NewGroup("chating", record)
	left := container.NewVScroll(b)
	line := canvas.NewLine(color.White)
	line.Resize(fyne.NewSize(1, 200))
	h3 := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), line)
	return fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(3), a, h3, left)

}
