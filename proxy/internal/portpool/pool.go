package portpool

import (
	"fmt"
	"sync"
)

// PoolStats 端口池统计信息
type PoolStats struct {
	Total     int `json:"total"`
	Used      int `json:"used"`
	Available int `json:"available"`
}

// Pool 端口池，管理可用端口的分配和释放
type Pool struct {
	mu         sync.Mutex
	rangeStart int
	rangeEnd   int
	used       map[int]string // port -> envID
}

// NewPool 创建端口池
func NewPool(rangeStart, rangeEnd int) *Pool {
	return &Pool{
		rangeStart: rangeStart,
		rangeEnd:   rangeEnd,
		used:       make(map[int]string),
	}
}

// Allocate 分配一个可用端口，返回端口号
func (p *Pool) Allocate(envID string) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for port := p.rangeStart; port <= p.rangeEnd; port++ {
		if _, ok := p.used[port]; !ok {
			p.used[port] = envID
			return port, nil
		}
	}
	return 0, fmt.Errorf("端口池已满，无可用端口")
}

// AllocateSpecific 分配指定端口
func (p *Pool) AllocateSpecific(port int, envID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if port < p.rangeStart || port > p.rangeEnd {
		return fmt.Errorf("端口 %d 不在范围 [%d, %d] 内", port, p.rangeStart, p.rangeEnd)
	}
	if _, ok := p.used[port]; ok {
		return fmt.Errorf("端口 %d 已被占用", port)
	}
	p.used[port] = envID
	return nil
}

// Release 释放指定端口
func (p *Pool) Release(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.used, port)
}

// ReleaseByEnvID 释放指定环境的所有端口
func (p *Pool) ReleaseByEnvID(envID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for port, id := range p.used {
		if id == envID {
			delete(p.used, port)
		}
	}
}

// Stats 返回端口池统计信息
func (p *Pool) Stats() PoolStats {
	p.mu.Lock()
	defer p.mu.Unlock()

	total := p.rangeEnd - p.rangeStart + 1
	used := len(p.used)
	return PoolStats{
		Total:     total,
		Used:      used,
		Available: total - used,
	}
}
