type CompilerAnalyzer struct {
    versionCache *ttlcache.Cache[string, semver.Version]
}

func (c *CompilerAnalyzer) ValidateCompatibility(callerVer, targetVer string) bool {
    // Check cache first
    cacheKey := callerVer + ":" + targetVer
    if result, found := c.versionCache.Get(cacheKey); found {
        return result
    }
    
    // Normalize versions
    normalizedCaller := normalizeCompilerVersion(callerVer)
    normalizedTarget := normalizeCompilerVersion(targetVer)
    
    // Reject unstable versions
    if isUnstableVersion(normalizedCaller) || isUnstableVersion(normalizedTarget) {
        c.versionCache.Set(cacheKey, false, 24*time.Hour)
        return false
    }

    // Compare major.minor versions
    result := semver.MajorMinor(normalizedCaller) == semver.MajorMinor(normalizedTarget)
    c.versionCache.Set(cacheKey, result, 24*time.Hour)
    return result
}