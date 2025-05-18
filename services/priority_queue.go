type AnalysisRequest struct {
    Contract    string
    Priority    int // 0=low, 1=normal, 2=critical
    RequestedAt time.Time
}

type PriorityQueue struct {
    items []*AnalysisRequest
    mu    sync.Mutex
    cond  *sync.Cond
}

func NewPriorityQueue() *PriorityQueue {
    pq := &PriorityQueue{
        items: make([]*AnalysisRequest, 0),
    }
    pq.cond = sync.NewCond(&pq.mu)
    return pq
}

func (pq *PriorityQueue) Push(req *AnalysisRequest) {
    pq.mu.Lock()
    defer pq.mu.Unlock()
    
    pq.items = append(pq.items, req)
    heap.Fix(pq, len(pq.items)-1)
    pq.cond.Signal() // Wake worker
}

func (pq *PriorityQueue) Pop() *AnalysisRequest {
    pq.mu.Lock()
    defer pq.mu.Unlock()
    
    for len(pq.items) == 0 {
        pq.cond.Wait() // Sleep until signaled
    }
    
    req := pq.items[0]
    pq.items[0] = pq.items[len(pq.items)-1]
    pq.items = pq.items[:len(pq.items)-1]
    heap.Fix(pq, 0)
    
    return req
}

func (pq *PriorityQueue) worker() {
    for {
        req := pq.Pop()
        timeout := getTimeoutForPriority(req.Priority)
        
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        analyzeContract(ctx, req.Contract)
        cancel()
    }
}