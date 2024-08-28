package main

import (
	"fmt" // 导入fmt包，用于格式化输出
	"net" // 导入net包，用于网络编程
	"os"
	"strconv" // 导入strconv包，用于字符串和其他类型的转换
	"strings" // 导入strings包，用于字符串操作
	"time"    // 导入time包，用于处理时间
)

func main() {
	// 定义服务端口为字符串，使用":1200"表示监听所有网络接口上的1200端口。
	service := ":8081"
	// 解析TCP端点的地址，使用IPv4。
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	// 检查解析TCP地址是否出错。
	checkError(err)
	// 在解析的地址上创建TCP监听器。
	listener, err := net.ListenTCP("tcp", tcpAddr)
	// 检查创建TCP监听器是否出错。
	checkError(err)
	// 进入无限循环，接受传入的连接。
	for {
		// 从监听器接受新的连接。
		conn, err := listener.Accept()
		// 如果接受连接时出错，跳过本次循环的剩余部分。
		if err != nil {
			continue
		}
		// 为了并发处理，使用goroutine来处理连接。
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	// 为从连接读取数据设置截止时间（从现在开始2分钟后）。
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	// 创建一个字节切片来保存传入的请求，设置最大大小以防止洪水攻击。
	request := make([]byte, 128)
	// 确保在函数返回时关闭连接。
	defer conn.Close()
	// 进入循环，读取和处理请求。
	for {
		// 从连接读取数据到请求字节切片中。
		read_len, err := conn.Read(request)

		// 如果读取时出错，打印错误信息并退出循环。
		if err != nil {
			fmt.Println(err)
			break
		}

		// 如果没有读取到数据（客户端关闭了连接），退出循环。
		if read_len == 0 {
			break
		} else if strings.TrimSpace(string(request[:read_len])) == "timestamp" {
			// 如果客户端发送的是"timestamp"，则回复当前的Unix时间戳。
			daytime := strconv.FormatInt(time.Now().Unix(), 10)
			conn.Write([]byte(daytime))
		} else {
			// 对于其他任何请求，回复当前时间的字符串形式。
			daytime := time.Now().String()
			conn.Write([]byte(daytime))
		}

		// 清除请求切片，为下一次读取操作做准备。
		request = make([]byte, 128)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatel error: %s", err.Error())
		os.Exit(1)
	}
}
