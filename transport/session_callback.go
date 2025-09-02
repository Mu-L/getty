package getty

// AddCloseCallback adds a close callback function to the Session
// 
// Parameters:
//   - handler: handler identifier, used to identify the source or type of the callback
//   - key: unique identifier key for the callback, used in combination with handler
//   - f: callback function to be executed when the session is closed
//
// Notes:
//   - If the session is already closed, this addition will be ignored
//   - The combination of handler and key must be unique, otherwise it will override previous callbacks
//   - Callback functions will be executed in the order they were added when the session closes
func (s *session) AddCloseCallback(handler, key any, f CallBackFunc) {
	if s.IsClosed() {
		return
	}
	s.closeCallbackMutex.Lock()
	s.closeCallback.Add(handler, key, f)
	s.closeCallbackMutex.Unlock()
}

// RemoveCloseCallback removes the specified Session close callback function
// 
// Parameters:
//   - handler: handler identifier of the callback to be removed
//   - key: unique identifier key of the callback to be removed
//
// Return value: none
//
// Notes:
//   - If the session is already closed, this removal operation will be ignored
//   - If no matching callback is found, this operation will have no effect
//   - The removal operation is thread-safe
func (s *session) RemoveCloseCallback(handler, key any) {
	if s.IsClosed() {
		return
	}
	s.closeCallbackMutex.Lock()
	s.closeCallback.Remove(handler, key)
	s.closeCallbackMutex.Unlock()
}

// invokeCloseCallbacks executes all registered close callback functions
// 
// Function description:
//   - Executes all registered close callbacks in the order they were added
//   - Uses read lock to protect the callback list, ensuring concurrency safety
//   - This method is typically called automatically when the session closes
//
// Notes:
//   - This is an internal method, not recommended for external direct calls
//   - If panic occurs during callback execution, it will be caught and logged
//   - Callback functions should avoid long blocking operations, async processing is recommended for time-consuming tasks
func (s *session) invokeCloseCallbacks() {
	s.closeCallbackMutex.RLock()
	s.closeCallback.Invoke()
	s.closeCallbackMutex.RUnlock()
}

// CallBackFunc defines the callback function type when Session closes
// 
// Usage notes:
//   - Callback function accepts no parameters
//   - Callback function returns no values
//   - Callback function should handle resource cleanup, state updates, etc.
//   - It's recommended to avoid accessing closed session state in callback functions
type CallBackFunc func()
