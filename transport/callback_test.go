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
	if cb.Count() != 0 {
		t.Errorf("Expected count for empty list is 0, but got %d", cb.Count())
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
	if cb.Count() != expectedCallbacks {
		t.Errorf("Expected callback count is %d, but got %d", expectedCallbacks, cb.Count())
	}

	// Test adding nil callback
	cb.Add(remove, remove, nil)
	if cb.Count() != expectedCallbacks {
		t.Errorf("Expected count after adding nil callback is %d, but got %d", expectedCallbacks, cb.Count())
	}

	// Replace an existing callback with a non-nil one; count should remain unchanged.
	cb.Add(remove, remove, func() { count += remove })
	if cb.Count() != expectedCallbacks {
		t.Errorf("Expected count after replacing existing callback is %d, but got %d", expectedCallbacks, cb.Count())
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

func TestCallbackInvokePanicSafe(t *testing.T) {
	cb := &callbacks{}
	var ran bool
	cb.Add("h", "k1", func() { panic("boom") })
	cb.Add("h", "k2", func() { ran = true })
	// Expect: Invoke swallows panics and continues executing remaining callbacks.
	cb.Invoke()
	if !ran {
		t.Errorf("Expected subsequent callbacks to run even if one panics")
	}
}
