type ABTestRunner struct {
    wasmAnalyzer        *WASMAnalyzer
    microserviceAnalyzer *SlitherAdapter
    testContracts       []Contract
    results             chan ABTestResult
}

func (r *ABTestRunner) Run(ctx context.Context) map[string]ABTestSummary {
    var wg sync.WaitGroup
    
    for _, contract := range r.testContracts {
        wg.Add(1)
        go func(c Contract) {
            defer wg.Done()
            
            // Run both analyzers in parallel
            wasmCh := make(chan AnalysisResult, 1)
            microCh := make(chan AnalysisResult, 1)
            
            go runWithTimeout(r.wasmAnalyzer, c, wasmCh, 500*time.Millisecond)
            go runWithTimeout(r.microserviceAnalyzer, c, microCh, 500*time.Millisecond)
            
            // Collect and compare results
            select {
            case <-ctx.Done():
                return
            case r.results <- compareResults(c, <-wasmCh, <-microCh):
            }
        }(contract)
    }
    
    go func() {
        wg.Wait()
        close(r.results)
    }()
    
    return aggregateResults(r.results)
}