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
		data := make([]byte, 1024)
		_, errConn := conn.Read(data)
		if errConn != nil {
		    fmt.Println(errConn)
			continue
		}
		fmt.Println(string(data))
	}
}
