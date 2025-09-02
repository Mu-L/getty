package getty

import (
	"testing"
)

func TestCallback(t *testing.T) {
	// 测试空列表
	cb := &callbacks{}
	if cb.Count() != 0 {
		t.Errorf("空列表期望数量为 0，实际为 %d", cb.Count())
	}
	
	// 测试添加回调函数
	var count, expected, remove, totalCount int
	totalCount = 10
	remove = 5
	
	// 添加多个回调函数
	for i := 1; i < totalCount; i++ {
		expected = expected + i
		func(ii int) { 
			cb.Add(ii, ii, func() { count = count + ii }) 
		}(i)
	}
	
	// 验证添加后的数量
	expectedCallbacks := totalCount - 1
	if cb.Count() != expectedCallbacks {
		t.Errorf("期望回调数量为 %d，实际为 %d", expectedCallbacks, cb.Count())
	}
	
	// 测试添加 nil 回调
	cb.Add(remove, remove, nil)
	if cb.Count() != expectedCallbacks {
		t.Errorf("添加 nil 回调后期望数量为 %d，实际为 %d", expectedCallbacks, cb.Count())
	}
	
	// 移除指定的回调
	cb.Remove(remove, remove)
	
	// 尝试移除不存在的回调
	cb.Remove(remove+1, remove+2)
	
	// 执行所有回调
	cb.Invoke()
	
	// 验证执行结果
	expectedCount := expected - remove
	if count != expectedCount {
		t.Errorf("期望执行结果为 %d，实际为 %d", expectedCount, count)
	}
	
	// 测试字符串类型的 handler 和 key
	cb2 := &callbacks{}
	
	// 添加回调
	cb2.Add("handler1", "key1", func() {})
	cb2.Add("handler2", "key2", func() {})
	cb2.Add("handler3", "key3", func() {})
	
	if cb2.Count() != 3 {
		t.Errorf("期望回调数量为 3，实际为 %d", cb2.Count())
	}
	
	// 移除中间的回调
	cb2.Remove("handler2", "key2")
	if cb2.Count() != 2 {
		t.Errorf("移除中间回调后期望数量为 2，实际为 %d", cb2.Count())
	}
	
	// 移除第一个回调
	cb2.Remove("handler1", "key1")
	if cb2.Count() != 1 {
		t.Errorf("移除第一个回调后期望数量为 1，实际为 %d", cb2.Count())
	}
	
	// 移除最后一个回调
	cb2.Remove("handler3", "key3")
	if cb2.Count() != 0 {
		t.Errorf("移除最后一个回调后期望数量为 0，实际为 %d", cb2.Count())
	}
	
	// 测试移除不存在的回调
	cb2.Add("handler1", "key1", func() {})
	cb2.Remove("handler2", "key2") // 尝试移除不存在的回调
	
	// 应该仍然有1个回调
	if cb2.Count() != 1 {
		t.Errorf("期望回调数量为 1，实际为 %d", cb2.Count())
	}
}
