package getty

import (
	"testing"
)

func TestCallback(t *testing.T) {
	// Test empty list
	cb := &callbacks{}
	if cb.Count() != 0 {
		t.Errorf("Expected count for empty list is 0, but got %d", cb.Count())
	}

	// Test adding callback functions
	var count, expected, remove, totalCount int
	totalCount = 10
	remove = 5

	// Add multiple callback functions
	for i := 1; i < totalCount; i++ {
		expected = expected + i
		func(ii int) {
			cb.Add(ii, ii, func() { count = count + ii })
		}(i)
	}

	// Verify count after adding
	expectedCallbacks := totalCount - 1
	if cb.Count() != expectedCallbacks {
		t.Errorf("Expected callback count is %d, but got %d", expectedCallbacks, cb.Count())
	}

	// Test adding nil callback
	cb.Add(remove, remove, nil)
	if cb.Count() != expectedCallbacks {
		t.Errorf("Expected count after adding nil callback is %d, but got %d", expectedCallbacks, cb.Count())
	}

	// Remove specified callback
	cb.Remove(remove, remove)

	// Try to remove non-existent callback
	cb.Remove(remove+1, remove+2)

	// Execute all callbacks
	cb.Invoke()

	// Verify execution result
	expectedCount := expected - remove
	if count != expectedCount {
		t.Errorf("Expected execution result is %d, but got %d", expectedCount, count)
	}

	// Test string type handler and key
	cb2 := &callbacks{}

	// Add callbacks
	cb2.Add("handler1", "key1", func() {})
	cb2.Add("handler2", "key2", func() {})
	cb2.Add("handler3", "key3", func() {})

	if cb2.Count() != 3 {
		t.Errorf("Expected callback count is 3, but got %d", cb2.Count())
	}

	// Remove middle callback
	cb2.Remove("handler2", "key2")
	if cb2.Count() != 2 {
		t.Errorf("Expected count after removing middle callback is 2, but got %d", cb2.Count())
	}

	// Remove first callback
	cb2.Remove("handler1", "key1")
	if cb2.Count() != 1 {
		t.Errorf("Expected count after removing first callback is 1, but got %d", cb2.Count())
	}

	// Remove last callback
	cb2.Remove("handler3", "key3")
	if cb2.Count() != 0 {
		t.Errorf("Expected count after removing last callback is 0, but got %d", cb2.Count())
	}

	// Test removing non-existent callback
	cb2.Add("handler1", "key1", func() {})
	cb2.Remove("handler2", "key2") // Try to remove non-existent callback

	// Should still have 1 callback
	if cb2.Count() != 1 {
		t.Errorf("Expected callback count is 1, but got %d", cb2.Count())
	}
}
