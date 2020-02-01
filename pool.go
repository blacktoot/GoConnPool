package pool
import (
    "time"
    "net"
    "sync"
)

//pool strcut 
type Pool struct {
    initialNum int
    maxNum int
    expiryTime time.Duration
    create func(addr string) (interface{}, error)
    close func(interface{}) error
    conn []*connection
    addr string
    running int //already apply num
    mutex sync.Mutex
}

type connection struct {
    //conn type
    conn interface{}
    //start time
    expiry time.Time
}
//new pool
func NewPool(opt ...Options) (*Pool, error){
    //default config
    configParam := config{
        initialSize : defaultInitialSize,
        maxSize : defaultMaxSize,
        expiryTime : defaultExpiryTime,
        addr : defaultAddr,
    }

    for _, o := range opt {
        o(&configParam)
    }
    
    p := &Pool{
        initialNum : configParam.initialSize,
        maxNum : configParam.maxSize,
        expiryTime : configParam.expiryTime,
        create : build,
        close : close,
        conn : make([]*connection, 0, configParam.maxSize),
        addr : configParam.addr,
        running : 0,
    }
    if p.initialNum <= 0 || p.maxNum <= 0 || p.initialNum > p.maxNum {
        return nil, initialConfigErr
    } else if p.expiryTime <= 0 {
        return nil, expiryConfigErr
    } else if p.create == nil {
        return nil, createConfigErr
    } else if p.close == nil {
        return nil, closeConfigErr
    } else {
        //initial conn
        for i := 0; i < p.initialNum; i++ {
            c, err := p.create(p.addr)
			if err != nil {
                continue
            }
			conn := &connection{
			    conn : c,
				expiry : time.Now(),
			}
            p.conn = append(p.conn, conn)
            p.running += 1
        }
        return p, nil
    }
}

//conn create 
func build(addr string) (interface{}, error) {
    return net.Dial("tcp", addr)
}

//conn close 
func close(conn interface{}) error {
    return conn.(net.Conn).Close()
}

// pool operate interface
type operate interface {
    Get() (*connection, error)
    Put(c interface{}) error
}

//get one conn from pool
func(p *Pool) Get() (*connection, error) {
    //lock
    p.mutex.Lock()
    defer p.mutex.Unlock()
    length := len(p.conn)
    if length > 0 {
        return p.conn[length-1], nil
    }

    //conn is empty
    if p.running >= p.maxNum {
        return nil, poolOverload
    }

    //new one conn
    c, err := p.create(p.addr)
	if err != nil {
        return nil, createConnFail
    }
	conn := &connection{
	    conn : c,
		expiry : time.Now(),
	}
    p.running += 1
    return conn, err
}

func(p *Pool) Put(c interface{}) error {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    if len(p.conn) >= p.maxNum {
        return poolOverload
    }
    
	conn := &connection{
	    conn : c,
	}
    p.conn = append(p.conn, conn)
    return nil
}