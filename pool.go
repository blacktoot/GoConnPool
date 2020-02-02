package pool
import (
    "time"
    "net"
    "sync"
    "sync/atomic"
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
    running int32 //already apply num
    mutex sync.Mutex
}

type connection struct {
    //conn type
    Conn interface{}
    //start time
    expiry time.Time
}

//singleton get pool instance
var (
    once sync.Once
    singleton *Pool
    err error
)
func GetPoolInst(opt ...Options) (*Pool, error) {
	once.Do(func() {
		singleton, err = newPool(opt...)
	})
	return singleton, err
}

//new pool
func newPool(opt ...Options) (*Pool, error) {
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
			    Conn : c,
				expiry : time.Now(),
			}
            p.conn = append(p.conn, conn)
            atomic.AddInt32(&p.running, 1)
        }
        go p.periodCleanExpiryConn()
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
    runningNum := int(atomic.LoadInt32(&p.running))
    if runningNum >= p.maxNum {
        return nil, poolOverload
    }

    //new one conn
    c, err := p.create(p.addr)
	if err != nil {
        return nil, createConnFail
    }
	conn := &connection{
	    Conn : c,
		expiry : time.Now(),
	}
	atomic.AddInt32(&p.running, 1)
    return conn, err
}

func(p *Pool) Put(c *connection) error {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    if len(p.conn) >= p.maxNum {
        return poolOverload
    }
    
    p.conn = append(p.conn, c)
    return nil
}

//period clean expiry conn
func (p *Pool) periodCleanExpiryConn() {
    timer := time.NewTicker(p.expiryTime)
    defer timer.Stop()
    for range timer.C {
	    //get expiry conn
	    t := time.Now().Add(-p.expiryTime)
	    p.mutex.Lock()
	    defer p.mutex.Unlock()
	    index := p.binSearch(t)
	    if index < 0 {
			return
	    }
	    expiryConn := p.conn[:index+1]
	    p.conn = p.conn[index+1:]
	    //deal expiry conn
	    for _, v := range expiryConn {
	        p.close(v)
	    }
	    expiryNum := int32(index+1)
	    atomic.AddInt32(&p.running, -expiryNum)
	}
}
