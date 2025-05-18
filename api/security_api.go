func SetupRouter(pipeline *ml.Pipeline) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery())
    
    // Use custom logger that doesn't block
    r.Use(AsyncLogger())
    
    // Configure connection pool
    r.Use(func(c *gin.Context) {
        c.Set("dbConn", dbPool.Get())
        defer dbPool.Put(c.MustGet("dbConn"))
        c.Next()
    })
    
    r.POST("/analyze/transaction", func(c *gin.Context) {
        var tx Transaction
        if err := c.ShouldBindJSON(&tx); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        
        score, anomalies := pipeline.Process(&tx)
        c.JSON(200, gin.H{
            "score": score,
            "anomalies": anomalies,
        })
    })
    
    return r
}