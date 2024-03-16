package main
import (
 "fmt"
 "bufio"
 "net"
 "sync"
 "io/ioutil"
 "log"
 "os"
 "strconv"
 "gopkg.in/yaml.v2"
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
type Config struct {
 Database struct {
 Host     string `yaml:"host"`
 Port     int    `yaml:"port"`
 Port2    int    `yaml:"port2"`
 } `yaml:"database"`
}
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
	argCount := len(os.Args) - 1
	if argCount == 0 {
		fmt.Println("没有传递任何参数。")
		fmt.Println("acade -h")
		fmt.Println("acade --help")
		os.Exit(1) // 使用非零状态码退出，表示错误  
	}

	// 如果参数数量不符合预期（例如，你需要恰好2个参数），提前返回  
	if argCount != 1 {
		fmt.Printf("需要1个参数，但传递了%d个参数。\n", argCount)
		os.Exit(1) // 使用非零状态码退出，表示错误  
	}

	yamlfilename:=os.Args[1]
	if yamlfilename=="-h" || yamlfilename=="--help"{
            fmt.Println("acade [path to settings.yaml]")
	    os.Exit(0)
	}
         // 读取YAML文件内容  
         yamlFile, err := ioutil.ReadFile(yamlfilename)
        if err != nil {
            log.Fatalf("无法读取YAML文件: %v", err)
         }
         // 解析YAML内容  
         var config Config
         err = yaml.Unmarshal(yamlFile, &config)
         if err != nil {
             log.Fatalf("无法解析YAML内容: %v", err)
         }
         // 输出解析后的内容  
         fmt.Printf("Host: %s\n", config.Database.Host)
         fmt.Printf("Port: %d\n", config.Database.Port)
	 fmt.Printf("Port2: %d\n", config.Database.Port2)
	 Host:=config.Database.Host
	 Port:=strconv.Itoa(config.Database.Port)
	 Port2:=strconv.Itoa(config.Database.Port2)
	 listener, err := net.Listen("tcp",Host+":"+Port)
	 if err != nil {
                fmt.Println(err)
                return
        }
	listener2, err2 := net.Listen("tcp",Host+":"+Port2)
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
	// 程序正常执行完毕，使用零状态码退出  
	os.Exit(0)
}
