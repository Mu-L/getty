/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package getty

import (
	"testing"
)

func TestCallback(t *testing.T) {
	// Test empty list
	cb := &callbacks{}
	if cb.Len() != 0 {
		t.Errorf("Expected count for empty list is 0, but got %d", cb.Len())
	}

	// Ensure invoking on an empty registry is a no-op (no panic).
	cb.Invoke()

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
	if cb.Len() != expectedCallbacks {
		t.Errorf("Expected callback count is %d, but got %d", expectedCallbacks, cb.Len())
	}

	// Test adding nil callback
	cb.Add(remove, remove, nil)
	if cb.Len() != expectedCallbacks {
		t.Errorf("Expected count after adding nil callback is %d, but got %d", expectedCallbacks, cb.Len())
	}

	// Replace an existing callback with a non-nil one; count should remain unchanged.
	cb.Add(remove, remove, func() { count += remove })
	if cb.Len() != expectedCallbacks {
		t.Errorf("Expected count after replacing existing callback is %d, but got %d", expectedCallbacks, cb.Len())
	}

	// Remove specified callback
	cb.Remove(remove, remove)

	// Try to remove non-existent callback
	cb.Remove(remove+1, remove+2)

	// Execute all callbacks
	cb.Invoke()

	// Verify execution result
	expectedSum := expected - remove
	if count != expectedSum {
		t.Errorf("Expected execution result is %d, but got %d", expectedSum, count)
	}

	// Test string type handler and key
	cb2 := &callbacks{}

	// Add callbacks
	cb2.Add("handler1", "key1", func() {})
	cb2.Add("handler2", "key2", func() {})
	cb2.Add("handler3", "key3", func() {})

	if cb2.Len() != 3 {
		t.Errorf("Expected callback count is 3, but got %d", cb2.Len())
	}

	// Remove middle callback
	cb2.Remove("handler2", "key2")
	if cb2.Len() != 2 {
		t.Errorf("Expected count after removing middle callback is 2, but got %d", cb2.Len())
	}

	// Remove first callback
	cb2.Remove("handler1", "key1")
	if cb2.Len() != 1 {
		t.Errorf("Expected count after removing first callback is 1, but got %d", cb2.Len())
	}

	// Remove last callback
	cb2.Remove("handler3", "key3")
	if cb2.Len() != 0 {
		t.Errorf("Expected count after removing last callback is 0, but got %d", cb2.Len())
	}

	// Test removing non-existent callback
	cb2.Add("handler1", "key1", func() {})
	cb2.Remove("handler2", "key2") // Try to remove non-existent callback

	// Should still have 1 callback
	if cb2.Len() != 1 {
		t.Errorf("Expected callback count is 1, but got %d", cb2.Len())
	}
}

func TestCallbackInvokePanicPropagation(t *testing.T) {
	cb := &callbacks{}
	cb.Add("h", "k1", func() { panic("boom") })

	// Test that panic is propagated (not swallowed by Invoke)
	defer func() {
		if r := recover(); r != nil {
			if r != "boom" {
				t.Errorf("Expected panic 'boom', got %v", r)
			}
		} else {
			t.Errorf("Expected panic to be propagated, but it was swallowed")
		}
	}()

	// This should panic and be caught by the defer above
	cb.Invoke()
}

func TestCallbackNonComparableTypes(t *testing.T) {
	cb := &callbacks{}

	// Test with non-comparable types (slice, map, function)
	nonComparableTypes := []struct {
		name     string
		handler  any
		key      any
		expected bool // whether the callback should be added
	}{
		{"slice_handler", []int{1, 2, 3}, "key", false},
		{"map_handler", map[string]int{"a": 1}, "key", false},
		{"func_handler", func() {}, "key", false},
		{"slice_key", "handler", []int{1, 2, 3}, false},
		{"map_key", "handler", map[string]int{"a": 1}, false},
		{"func_key", "handler", func() {}, false},
		{"both_non_comparable", []int{1}, map[string]int{"a": 1}, false},
		{"comparable_types", "handler", "key", true},
		{"nil_values", nil, nil, true},
		{"mixed_comparable", "handler", 123, true},
	}

	for _, tt := range nonComparableTypes {
		t.Run(tt.name, func(t *testing.T) {
			initialCount := cb.Len()

			// Try to add callback
			cb.Add(tt.handler, tt.key, func() {})

			// Check if callback was added
			finalCount := cb.Len()
			if tt.expected {
				if finalCount != initialCount+1 {
					t.Errorf("Expected callback to be added, but count remained %d", initialCount)
				}
				// Clean up for next test
				cb.Remove(tt.handler, tt.key)
			} else {
				if finalCount != initialCount {
					t.Errorf("Expected callback to be ignored, but count changed from %d to %d", initialCount, finalCount)
				}
			}
		})
	}

	// Test Remove with non-comparable types
	t.Run("RemoveNonComparable", func(t *testing.T) {
		initialCount := cb.Len()

		// Try to remove with non-comparable types
		cb.Remove([]int{1, 2, 3}, map[string]int{"a": 1})

		// Count should remain unchanged
		if cb.Len() != initialCount {
			t.Errorf("Expected count to remain %d after removing non-comparable types, but got %d", initialCount, cb.Len())
		}
	})
}
