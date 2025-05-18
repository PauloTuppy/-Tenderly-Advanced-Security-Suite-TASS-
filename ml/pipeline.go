type FeatureExtractor interface {
    Extract(tx *Transaction) []float32
}

type AnomalyDetector interface {
    Detect(features []float32) float32
}

type Pipeline struct {
    extractors []FeatureExtractor
    detectors  []AnomalyDetector
    threshold  float32
}

func (p *Pipeline) Process(tx *Transaction) (score float32, anomalies []string) {
    features := make([][]float32, len(p.extractors))
    // Extract features concurrently
    var wg sync.WaitGroup
    for i, extractor := range p.extractors {
        wg.Add(1)
        go func(i int, e FeatureExtractor) {
            features[i] = e.Extract(tx)
            wg.Done()
        }(i, extractor)
    }
    wg.Wait()
    
    // Run detectors and aggregate scores
    // ...
}