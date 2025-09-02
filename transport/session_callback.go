package getty

// AddCloseCallback 添加 Session 关闭回调函数
// 
// 参数说明:
//   - handler: 处理器标识，用于标识回调的来源或类型
//   - key: 回调的唯一标识键，与 handler 组合使用
//   - f: 要执行的回调函数，在 session 关闭时自动调用
//
// 注意事项:
//   - 如果 session 已经关闭，则忽略此次添加
//   - handler 和 key 的组合必须唯一，否则会覆盖之前的回调
//   - 回调函数会在 session 关闭时按照添加顺序执行
func (s *session) AddCloseCallback(handler, key any, f CallBackFunc) {
	if s.IsClosed() {
		return
	}
	s.closeCallbackMutex.Lock()
	s.closeCallback.Add(handler, key, f)
	s.closeCallbackMutex.Unlock()
}

// RemoveCloseCallback 移除指定的 Session 关闭回调函数
// 
// 参数说明:
//   - handler: 要移除的回调的处理器标识
//   - key: 要移除的回调的唯一标识键
//
// 返回值: 无
//
// 注意事项:
//   - 如果 session 已经关闭，则忽略此次移除操作
//   - 如果找不到匹配的回调，此操作不会产生任何效果
//   - 移除操作是线程安全的
func (s *session) RemoveCloseCallback(handler, key any) {
	if s.IsClosed() {
		return
	}
	s.closeCallbackMutex.Lock()
	s.closeCallback.Remove(handler, key)
	s.closeCallbackMutex.Unlock()
}

// invokeCloseCallbacks 执行所有注册的关闭回调函数
// 
// 功能说明:
//   - 按照添加顺序依次执行所有注册的关闭回调
//   - 使用读锁保护回调列表，确保并发安全
//   - 此方法通常在 session 关闭时自动调用
//
// 注意事项:
//   - 此方法是内部方法，不建议外部直接调用
//   - 回调执行过程中如果发生 panic，会被捕获并记录日志
//   - 回调函数应该避免长时间阻塞，建议异步处理耗时操作
func (s *session) invokeCloseCallbacks() {
	s.closeCallbackMutex.RLock()
	s.closeCallback.Invoke()
	s.closeCallbackMutex.RUnlock()
}

// CallBackFunc 定义 Session 关闭时的回调函数类型
// 
// 使用说明:
//   - 回调函数不接受任何参数
//   - 回调函数不返回任何值
//   - 回调函数应该处理资源清理、状态更新等操作
//   - 建议在回调函数中避免访问已关闭的 session 状态
type CallBackFunc func()
