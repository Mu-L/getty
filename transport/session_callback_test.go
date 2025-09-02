package getty

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSessionCallback(t *testing.T) {
	// 测试基本的添加、移除和执行回调功能
	t.Run("BasicCallback", func(t *testing.T) {
		s := &session{
			once: &sync.Once{},
			done: make(chan struct{}),
			closeCallback: callbacks{},
		}
		
		var callbackExecuted bool
		var callbackMutex sync.Mutex
		
		callback := func() {
			callbackMutex.Lock()
			callbackExecuted = true
			callbackMutex.Unlock()
		}
		
		// 添加回调
		s.AddCloseCallback("testHandler", "testKey", callback)
		if s.closeCallback.Count() != 1 {
			t.Errorf("期望回调数量为 1，实际为 %d", s.closeCallback.Count())
		}
		
		// 测试移除回调
		s.RemoveCloseCallback("testHandler", "testKey")
		if s.closeCallback.Count() != 0 {
			t.Errorf("期望回调数量为 0，实际为 %d", s.closeCallback.Count())
		}
		
		// 重新添加回调
		s.AddCloseCallback("testHandler", "testKey", callback)
		
		// 测试关闭时回调执行
		go func() {
			time.Sleep(10 * time.Millisecond)
			s.stop()
		}()
		
		// 等待回调执行
		time.Sleep(50 * time.Millisecond)
		
		callbackMutex.Lock()
		if !callbackExecuted {
			t.Error("回调函数未被执行")
		}
		callbackMutex.Unlock()
	})
	
	// 测试多个回调的添加、移除和执行
	t.Run("MultipleCallbacks", func(t *testing.T) {
		s := &session{
			once: &sync.Once{},
			done: make(chan struct{}),
			closeCallback: callbacks{},
		}
		
		var callbackCount int
		var callbackMutex sync.Mutex
		
		// 添加多个回调
		totalCallbacks := 3
		for i := 0; i < totalCallbacks; i++ {
			index := i // 捕获循环变量
			callback := func() {
				callbackMutex.Lock()
				callbackCount++
				callbackMutex.Unlock()
			}
			s.AddCloseCallback(fmt.Sprintf("handler%d", index), fmt.Sprintf("key%d", index), callback)
		}
		
		if s.closeCallback.Count() != totalCallbacks {
			t.Errorf("期望回调数量为 %d，实际为 %d", totalCallbacks, s.closeCallback.Count())
		}
		
		// 移除一个回调
		s.RemoveCloseCallback("handler0", "key0")
		expectedAfterRemove := totalCallbacks - 1
		if s.closeCallback.Count() != expectedAfterRemove {
			t.Errorf("期望回调数量为 %d，实际为 %d", expectedAfterRemove, s.closeCallback.Count())
		}
		
		// 测试关闭时剩余回调执行
		go func() {
			time.Sleep(10 * time.Millisecond)
			s.stop()
		}()
		
		time.Sleep(50 * time.Millisecond)
		
		callbackMutex.Lock()
		if callbackCount != expectedAfterRemove {
			t.Errorf("期望执行的回调数量为 %d，实际为 %d", expectedAfterRemove, callbackCount)
		}
		callbackMutex.Unlock()
	})
	
	// 测试 invokeCloseCallbacks 功能
	t.Run("InvokeCloseCallbacks", func(t *testing.T) {
		s := &session{
			once: &sync.Once{},
			done: make(chan struct{}),
			closeCallback: callbacks{},
		}
		
		var callbackResults []string
		var callbackMutex sync.Mutex
		
		// 添加多个不同类型的关闭回调
		callbacks := []struct {
			handler string
			key     string
			action  string
		}{
			{"cleanup", "resources", "清理资源"},
			{"cleanup", "connections", "关闭连接"},
			{"logging", "audit", "记录审计日志"},
			{"metrics", "stats", "更新统计信息"},
		}
		
		// 注册所有回调
		for _, cb := range callbacks {
			cbCopy := cb // 捕获循环变量
			callback := func() {
				callbackMutex.Lock()
				callbackResults = append(callbackResults, cbCopy.action)
				callbackMutex.Unlock()
			}
			s.AddCloseCallback(cbCopy.handler, cbCopy.key, callback)
		}
		
		// 验证回调数量
		expectedCount := len(callbacks)
		if s.closeCallback.Count() != expectedCount {
			t.Errorf("期望回调数量为 %d，实际为 %d", expectedCount, s.closeCallback.Count())
		}
		
		// 手动调用关闭回调（模拟 invokeCloseCallbacks）
		callbackMutex.Lock()
		callbackResults = nil // 清空之前的结果
		callbackMutex.Unlock()
		
		// 执行所有关闭回调
		s.closeCallback.Invoke()
		
		// 等待回调执行完成
		time.Sleep(10 * time.Millisecond)
		
		// 验证所有回调都被执行
		callbackMutex.Lock()
		if len(callbackResults) != expectedCount {
			t.Errorf("期望执行 %d 个回调，实际执行了 %d 个", expectedCount, len(callbackResults))
		}
		
		// 验证回调执行顺序（应该按照添加顺序执行）
		expectedActions := []string{"清理资源", "关闭连接", "记录审计日志", "更新统计信息"}
		for i, result := range callbackResults {
			if i < len(expectedActions) && result != expectedActions[i] {
				t.Errorf("位置 %d: 期望执行 '%s'，实际执行了 '%s'", i, expectedActions[i], result)
			}
		}
		callbackMutex.Unlock()
		
		// 测试移除回调后再次执行
		s.RemoveCloseCallback("cleanup", "resources")
		
		callbackMutex.Lock()
		callbackResults = nil
		callbackMutex.Unlock()
		
		// 再次执行回调
		s.closeCallback.Invoke()
		time.Sleep(10 * time.Millisecond)
		
		// 验证移除后的执行结果
		callbackMutex.Lock()
		expectedAfterRemove := expectedCount - 1
		if len(callbackResults) != expectedAfterRemove {
			t.Errorf("移除一个回调后期望执行 %d 个回调，实际执行了 %d 个", expectedAfterRemove, len(callbackResults))
		}
		callbackMutex.Unlock()
	})
	
	// 测试边界情况
	t.Run("EdgeCases", func(t *testing.T) {
		// 测试空回调列表的情况
		s := &session{
			once: &sync.Once{},
			done: make(chan struct{}),
			closeCallback: callbacks{},
		}
		
		// 验证空列表
		if s.closeCallback.Count() != 0 {
			t.Errorf("空列表期望数量为 0，实际为 %d", s.closeCallback.Count())
		}
		
		// 执行空回调列表（不应该panic）
		s.closeCallback.Invoke()
		
		// 添加一个回调然后移除，再次执行
		s.AddCloseCallback("test", "key", func() {})
		s.RemoveCloseCallback("test", "key")
		
		// 移除后执行空列表（不应该panic）
		s.closeCallback.Invoke()
		
		if s.closeCallback.Count() != 0 {
			t.Errorf("移除后期望数量为 0，实际为 %d", s.closeCallback.Count())
		}
	})
}
