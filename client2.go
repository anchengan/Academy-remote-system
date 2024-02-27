package main  
  
import (  
	"bufio"  
	  "fmt" 
	//  "io"
 //"os/exec"   
	"net"  
	"os"  
	"strings"  
)  

func main() {  
	server := "localhost:8000"  
	server2:="localhost:8001"
	// 连接到服务器  
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
       
	// 处理从服务器接收的消息  
	//go func(){
//		reader2 := bufio.NewReader(conn2)
  //              for {
    //                    message2, err2 := reader2.ReadString('\n')
      //                  if err2 != nil {
        //                        fmt.Print("Connection lost: ", err2)
          //                     return
           //             }
             //          input_lock=false
               //         fmt.Print("Reply from server: " + message2)
		//	input_lock=true
//                }
//	}()
	//go func() {  
		//reader := bufio.NewReader(conn)  
		//for {  
		//	message, err := reader.ReadString('\n')  
		//	if err != nil {  
		//		fmt.Print("Connection lost: ", err)  
		//		return  
		//	}  
		//	fmt.Print("Message from server: " + message)  
		//	parts := strings.SplitN(message, " ", 2)  
  
		//	if len(parts) == 2 {  
		// 取后面那段  
		//secondPart := parts[1]  
				// 分割命令和参数
                  //      args := []string{"/bin/bash", "-c", secondPart}

                        // 创建命令对象
                //                cmd := exec.Command(args[0], args[1:]...)

                           // 执行命令并获取输出
                           // output, err := cmd.Output()
                          //  if err != nil {
                          //    fmt.Println("执行命令时出错:", err)
                        //        return
                      //        }
		      
                       // 打印命令输出
                    //fmt.Println(string(output))
                   //fmt.Fprintf(conn2, string(output)+"\n")

          //          }
	//	}  
	//}()  
  
	// 读取用户输入并发送到服务器  
 
	fmt.Println("Connected to chat server. Type 'exit' to disconnect.")  
	reader2 := bufio.NewReader(conn2)
 
pauseChan := true

scanner := bufio.NewReader(os.Stdin) 
	for{ 
    	
           if pauseChan==true{  
		//scanner := bufio.NewReader(os.Stdin) 
		   text,_ := scanner.ReadString('\n') 
                //text = strings.TrimSuffix(text, "\n") 
		// 退出命令  
		if strings.ToLower(text) == "exit" {  
			break  
		}  
  
		// 发送消息到服务器  
		fmt.Fprintf(conn, text) 
		pauseChan=false
	}else{
pauseChan=true

              
                        message2, err2 := reader2.ReadString('\n')

    if message2=="EOFACA\n"{
        continue
  }
                        if err2 != nil {
				fmt.Print("Connection lost: ", err2)
                                return
                        }
                 pauseChan=false       
                        fmt.Print("Reply from server: " + message2)
                      
		}
	}	
            
  
}
