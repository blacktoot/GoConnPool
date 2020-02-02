package main

import (
    "fmt"
    "net"
	"github.com/blacktoot/GoConnPool/pool"
)

func main() {
    p, err := pool.GetPoolInst(
	    pool.ConfigInitSize(3),
	)
    if err != nil {
        fmt.Println(err)
        return
    }
    conn, err = p.Get()
    if err != nil {
        fmt.Println(err)
        return
    }
    //deal transaction
    n, err := conn.Conn.(net.Conn).Write([]byte("This is a new conn"))
    if err != nil {
		fmt.Println("write data fail")
		return
    }
    fmt.Println("send data num:", n)
    if err := p.Put(conn); err != nil {
        fmt.Println(err)
    }
}
