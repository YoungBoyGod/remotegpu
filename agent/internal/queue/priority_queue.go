package queue

import (
	"github.com/YoungBoyGod/remotegpu-agent/internal/models"
)

// taskHeap 实现 heap.Interface
type taskHeap []*models.Task

func (h taskHeap) Len() int { return len(h) }

func (h taskHeap) Less(i, j int) bool {
	// 优先级数字越小优先级越高
	if h[i].Priority != h[j].Priority {
		return h[i].Priority < h[j].Priority
	}
	// 同优先级按创建时间排序
	return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h taskHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *taskHeap) Push(x any) {
	*h = append(*h, x.(*models.Task))
}

func (h *taskHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
