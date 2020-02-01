package pool
import(
    "time" 
)
//default config
const (
    defaultInitialSize = 5
    defaultMaxSize = 10
    defaultExpiryTime = 5*time.Second
    defaultAddr string = "127.0.0.1:8080"
)
//pool config
type config struct {
    initialSize int
    maxSize int
    expiryTime time.Duration
    addr string
}

//func option
type Options func(*config)

func ConfigInitSize(size int) Options{
    return func(c *config) {
        c.initialSize = size
    }
}

func ConfigMaxSize(size int) Options {
    return func(c *config) {
        c.maxSize = size
    }
}

func ConfigExpiryTime(t time.Duration) Options {
    return func(c *config) {
        c.expiryTime = t
    }
}

func ConfigAddr(addr string) Options {
    return func(c *config) {
        c.addr = addr
    }
}