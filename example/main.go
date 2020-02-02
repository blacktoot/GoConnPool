package main

import (
    "fmt"
	"github.com/blacktoot/GoConnPool/pool"
)

func main() {
    p, err := pool.GetPoolInst(
	    pool.ConfigInitSize(9),
	)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(p)
}
