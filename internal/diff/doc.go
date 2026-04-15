// Package diff provides functionality for comparing two versions of a
// HashiCorp Vault secret and rendering the differences.
//
// Usage:
//
//	old := map[string]interface{}{"db_pass": "old123"}
//	new := map[string]interface{}{"db_pass": "new456", "api_key": "abc"}
//
//	result := diff.Compare("secret/myapp", old, new)
//	diff.Render(os.Stdout, result, diff.FormatOptions{
//		ColorEnabled: true,
//		MaskValues:   false,
//	})
//
// ChangeType values: Added, Removed, Modified, Unchanged.
// Use Result.HasChanges() to detect if any meaningful differences exist.
package diff
