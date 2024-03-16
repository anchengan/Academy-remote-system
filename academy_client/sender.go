package main
import "C"
import "fmt"
import "net"
import "os"
import "bufio"
import "strings"
//export Sender  
func Sender(name *C.char,passwd *C.char,ip_info *C.char,port_info *C.char) {
	//fmt.Printf("device code: %s,password:%s\n", C.GoString(name),C.GoString(passwd))
        //fmt.Printf("ip info: %s,port info:%s\n", C.GoString(ip_info),C.GoString(port_info))
        // 调用strings.Split一次，并将结果存储在切片中
        systemlock:=true
        ips := strings.Split( C.GoString(ip_info), "-")
	ips=ips[:len(ips)-1]
        ports := strings.Split( C.GoString(port_info), "-")
	ports=ports[:len(ports)-1]
	for{
	    systemlock=true
            for id, ip := range ips {
                fmt.Println(ip)
                fmt.Println(ports[2*id])
                fmt.Println(ports[2*id+1])
		server := ip+":"+ports[2*id]
                server2 :=ip+":"+ports[2*id+1]
		go func(){
                    conn, err := net.Dial("tcp", server)
                    conn2,err2 := net.Dial("tcp",server2)
		    if err != nil {
                        fmt.Println(err)
                        return
                    }
                    if err2 != nil {
                        fmt.Println(err2)
                        return
                    }
                    defer conn.Close()
                    defer conn2.Close()
		    fmt.Println("Connected to chat server. Type 'exit' to disconnect.")
                    reader2 := bufio.NewReader(conn2)
		    pauseChan := true
		    scanner := bufio.NewReader(os.Stdin)
		    for{
                        if pauseChan==true{
                            text,_:=scanner.ReadString('\n')
			    if strings.ToLower(text) == "exit"{
                                break
	                    }
			    text=C.GoString(name)+string('-')+C.GoString(passwd)+string('-')+text
			    fmt.Fprintf(conn,text)
			    pauseChan=false
		        }else{
                            pauseChan=true
			    message2,err2 := reader2.ReadString('\n')
			    if message2=="EOFACA\n"{
				    continue
		            }
			    if err2 != nil {
                                fmt.Print("Connection lost: ",err2)
				return
		            }
			    pauseChan=false
			    fmt.Print("Reply from server: " + message2)
			}
                    }
		systemlock=false
	        }()
	    }
            for{
                if systemlock==false{
                    break
                }
            }
        }
}
func main() {} // main函数必须定义，但在这里不是执行入口 
