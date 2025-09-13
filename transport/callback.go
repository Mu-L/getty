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

// callbackCommon represents a node in the callback linked list
// Each node contains handler identifier, key, callback function and pointer to next node
type callbackCommon struct {
	handler interface{}        // Handler identifier, used to identify the source or type of callback
	key     interface{}        // Unique identifier key for callback, used in combination with handler
	call    func()             // Actual callback function to be executed
	next    *callbackCommon    // Pointer to next node, forming linked list structure
}

// callbacks is a singly linked list structure for managing multiple callback functions
// Supports dynamic addition, removal and execution of callbacks
type callbacks struct {
	first *callbackCommon    // Pointer to the first node of the linked list
	last  *callbackCommon    // Pointer to the last node of the linked list, used for quick addition of new nodes
}

// Add adds a new callback function to the callback linked list
// Parameters:
//   - handler: Handler identifier, can be any type
//   - key: Unique identifier key for callback, used in combination with handler
//   - callback: Callback function to be executed, ignored if nil
func (t *callbacks) Add(handler, key interface{}, callback func()) {
	// Prevent adding empty callback function
	if callback == nil {
		return
	}
	
	// Create new callback node
	newItem := &callbackCommon{handler, key, callback, nil}
	
	if t.first == nil {
		// If linked list is empty, new node becomes the first node
		t.first = newItem
	} else {
		// Otherwise add new node to the end of linked list
		t.last.next = newItem
	}
	// Update pointer to last node
	t.last = newItem
}

// Remove removes the specified callback function from the callback linked list
// Parameters:
//   - handler: Handler identifier of the callback to be removed
//   - key: Unique identifier key of the callback to be removed
// Note: If no matching callback is found, this method has no effect
func (t *callbacks) Remove(handler, key interface{}) {
	var prev *callbackCommon
	
	// Traverse linked list to find the node to be removed
	for callback := t.first; callback != nil; prev, callback = callback, callback.next {
		// Found matching node
		if callback.handler == handler && callback.key == key {
			if t.first == callback {
				// If it's the first node, update first pointer
				t.first = callback.next
			} else if prev != nil {
				// If it's a middle node, update the next pointer of the previous node
				prev.next = callback.next
			}
			
			if t.last == callback {
				// If it's the last node, update last pointer
				t.last = prev
			}
			
			// Return immediately after finding and removing
			return
		}
	}
}

// Invoke executes all registered callback functions in the linked list
// Executes each callback in the order they were added
// Note: If a callback function is nil, it will be skipped
func (t *callbacks) Invoke() {
	// Traverse the entire linked list starting from the head node
	for callback := t.first; callback != nil; callback = callback.next {
		// Ensure callback function is not nil before executing
		if callback.call != nil {
			callback.call()
		}
	}
}

// Count returns the number of callback functions in the linked list
// Return value: Total number of currently registered callback functions
func (t *callbacks) Count() int {
	var count int
	
	// Traverse linked list to count
	for callback := t.first; callback != nil; callback = callback.next {
		count++
	}
	
	return count
}
