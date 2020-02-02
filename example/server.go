package main
import "fmt"
import "net"
const (
    defaultAddr string = "127.0.0.1:8080"
)
func main() {
    l, err := net.Listen("tcp", defaultAddr)
	defer l.Close()
	if err != nil {
	    fmt.Println(err)
		return
	}
	for {
	    conn, err := l.Accept()
		if err != nil {
		    fmt.Println(err)
			continue
		}
		go dealConn(conn)
	}
}

func dealConn(c net.Conn) {
	data := make([]byte, 1024)
	_, errConn := c.Read(data)
	if errConn != nil {
		fmt.Println(errConn)
		return
	}
	fmt.Println(string(data))
}
