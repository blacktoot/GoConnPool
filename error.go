package pool
import (
    "errors"
)

var (
    initialConfigErr = errors.New("inital config error") 
    createConfigErr = errors.New("build config error")
    closeConfigErr = errors.New("close config error")
    expiryConfigErr = errors.New("expiry config error")
    poolOverload = errors.New("pool over load")
    createConnFail = errors.New("create conn fail")
)