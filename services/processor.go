func processContract(contract string, priority int) {
    // Check if this is a hot contract (recently analyzed)
    if result, found := hotCache.Get(contract); found {
        // Fast path - return cached result
        metrics.IncCounter("cache.hit")
        return result
    }

    // Check if already being processed
    if inProgress.Contains(contract) {
        // Add to waiting list with priority
        waitList.Add(contract, priority)
        metrics.IncCounter("contract.waiting")
        return
    }

    // Add to processing set
    inProgress.Add(contract)
    
    // Enqueue for processing with appropriate priority
    queue.Push(&AnalysisRequest{
        Contract: contract,
        Priority: priority,
        RequestedAt: time.Now(),
    })
    
    metrics.IncCounter("contract.queued")
}