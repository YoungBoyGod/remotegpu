package forwarder

import (
	"io"
	"log/slog"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
)

// TCPForwarder TCP 端口转发器
type TCPForwarder struct {
	listenPort int
	targetAddr string
	listener   net.Listener
	stopCh     chan struct{}
	wg         sync.WaitGroup
	connCount  atomic.Int64
}

// NewTCPForwarder 创建 TCP 转发器
func NewTCPForwarder(listenPort int, targetHost string, targetPort int) *TCPForwarder {
	return &TCPForwarder{
		listenPort: listenPort,
		targetAddr: net.JoinHostPort(targetHost, strconv.Itoa(targetPort)),
		stopCh:     make(chan struct{}),
	}
}

// Start 启动转发器，监听端口并接受连接
func (f *TCPForwarder) Start() error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(f.listenPort))
	if err != nil {
		return err
	}
	f.listener = ln

	f.wg.Add(1)
	go f.acceptLoop()

	slog.Info("TCP 转发器已启动", "listen", f.listenPort, "target", f.targetAddr)
	return nil
}

// Stop 停止转发器
func (f *TCPForwarder) Stop() {
	close(f.stopCh)
	if f.listener != nil {
		f.listener.Close()
	}
	f.wg.Wait()
	slog.Info("TCP 转发器已停止", "listen", f.listenPort)
}

// ConnCount 返回当前活跃连接数
func (f *TCPForwarder) ConnCount() int64 {
	return f.connCount.Load()
}

// acceptLoop 接受连接循环
func (f *TCPForwarder) acceptLoop() {
	defer f.wg.Done()

	for {
		conn, err := f.listener.Accept()
		if err != nil {
			select {
			case <-f.stopCh:
				return
			default:
				slog.Error("接受连接失败", "port", f.listenPort, "error", err)
				continue
			}
		}

		f.wg.Add(1)
		go f.handleConn(conn)
	}
}

// handleConn 处理单个连接，双向转发数据
func (f *TCPForwarder) handleConn(src net.Conn) {
	defer f.wg.Done()
	defer src.Close()

	f.connCount.Add(1)
	defer f.connCount.Add(-1)

	dst, err := net.Dial("tcp", f.targetAddr)
	if err != nil {
		slog.Error("连接目标失败", "target", f.targetAddr, "error", err)
		return
	}
	defer dst.Close()

	// 双向拷贝
	done := make(chan struct{})
	go func() {
		io.Copy(dst, src)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(src, dst)
		done <- struct{}{}
	}()

	// 等待任意一个方向结束
	<-done
}
