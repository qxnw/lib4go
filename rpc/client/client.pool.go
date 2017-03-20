package client

import "github.com/qxnw/lib4go/concurrent/cmap"

type poolOption struct {
	concurrent int
	minConn    int
}

//PoolOption 连接池配置选项
type PoolOption func(*poolOption)

//ClientPool 客户端连接池
type ClientPool struct {
	address string
	opt     *poolOption
	clients chan *Client
	cache   cmap.ConcurrentMap
}

//NewClientPool 初始化客户端连接池
func NewClientPool(address string, opts ...PoolOption) (p *ClientPool) {
	p = &ClientPool{opt: &poolOption{minConn: 1, concurrent: 10000}, address: address}
	p.cache = cmap.New()
	for _, opt := range opts {
		opt(p.opt)
	}
	p.clients = make(chan *Client, p.opt.minConn)
	for i := 0; i < p.opt.minConn; i++ {
		client, _ := p.create()
		p.clients <- client
	}
	return
}
