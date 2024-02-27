package main
import (
	"bufio"
	"fmt"
	"net"
	"sync"
)
type client struct {
	conn     net.Conn
	connid string
	name     string
	messages chan  string
}
type Message struct {
    conn net.Conn
    connid string
    SenderName string
    Content    string
}
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan Message) // 所有接收的消息  
)

func broadcast() {
    clients := make(map[string]client) // 键是客户端的名称，值是客户端结构体实例  

    for {
        select {
        case msg := <-messages:
            // 将消息广播给除了发送者以外的所有客户端  
            for _,cli := range clients {
                if cli.name != msg.SenderName && cli.connid==msg.connid {
		   cli.messages<- msg.Content
                }
            }
        case newClient := <-entering:
            clients[newClient.name] = newClient
        case leavingClient := <-leaving:
	delete(clients, leavingClient.name)
           close(leavingClient.messages)
        }
    }
}
func handleConn(conn net.Conn,connid string) {
	input := bufio.NewScanner(conn)
	var wg sync.WaitGroup
	// 创建一个新的客户端  
	ch := make(chan string,1)  // 客户端的消息通道  
	go clientWriter(conn, ch)
	cli := client{conn: conn, connid:connid,name: conn.RemoteAddr().String(), messages: ch}
	entering <- cli
        fmt.Println(cli.name)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for input.Scan() {
			messages <-Message{SenderName: cli.name, Content: input.Text(),conn:conn,connid:connid} // cli.name + ": " + input.Text()  
			_, err := conn.Write([]byte("Received your message\n"))
			 if err != nil {
				fmt.Println("Error writing to connection:", err.Error())
				break
			 }
		}
	}()
	wg.Wait()
	leaving <- cli
	conn.Close()
}
func clientWriter(conn net.Conn, ch <-chan string) {

	for msg := range ch {
		fmt.Fprintln(conn, msg) // 注意：网络写入应该有错误处理  
	}
}
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
  listener2, err2 := net.Listen("tcp", "localhost:8001")
        if err2 != nil {
                fmt.Println(err2)
                return
        }
	go broadcast()
	go func(){
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConn(conn,"1")
	}
}()
go func(){
	for{
                conn2, err2 := listener2.Accept()
                if err2 != nil {
                        fmt.Println(err2)
                        continue
                }
                go handleConn(conn2,"2")
        }
}()
for{}
}
