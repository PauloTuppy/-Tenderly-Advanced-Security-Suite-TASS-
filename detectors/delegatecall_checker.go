func DetectDelegatecallMismatch(caller, target *Contract) []Issue {
    var issues []Issue
    // Compare storage layouts
    storageLayoutMismatches := compareStorageLayouts(caller.StorageLayout, target.StorageLayout)
    // Check for delegatecall-safe annotations
    if !target.HasAnnotation("delegatecall-safe") && len(storageLayoutMismatches) > 0 {
        issues = append(issues, NewIssue("Storage layout mismatch in delegatecall target", High))
    }
    return issues
}