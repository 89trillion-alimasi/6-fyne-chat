package Protobuf

type ClientMessage struct {
	Type    string
	Msg     string
	Usrname string
}

func GetExitMessage() *Communication {
	m := &Communication{
		Type: "2",
		Msg:  "",
	}
	return m
}

func GettalkMessage(message string) *Communication {
	m := &Communication{
		Type: "1",
		Msg:  message,
	}
	return m
}

func Getuserlist() *Communication {
	m := &Communication{
		Type: "3",
		Msg:  "",
	}
	return m
}
