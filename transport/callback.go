package getty

// callbackCommon 表示回调链表中的一个节点
// 每个节点包含处理器标识、键值、回调函数和指向下一个节点的指针
type callbackCommon struct {
	handler interface{}        // 处理器标识，用于标识回调的来源或类型
	key     interface{}        // 回调的唯一标识键，与 handler 组合使用
	call    func()             // 实际要执行的回调函数
	next    *callbackCommon    // 指向下一个节点的指针，形成链表结构
}

// callbacks 是一个单向链表结构，用于管理多个回调函数
// 支持动态添加、移除和执行回调
type callbacks struct {
	first *callbackCommon    // 指向链表第一个节点的指针
	last  *callbackCommon    // 指向链表最后一个节点的指针，用于快速添加新节点
}

// Add 向回调链表中添加一个新的回调函数
// 参数说明:
//   - handler: 处理器标识，可以是任意类型
//   - key: 回调的唯一标识键，与 handler 组合使用
//   - callback: 要执行的回调函数，如果为 nil 则忽略
func (t *callbacks) Add(handler, key interface{}, callback func()) {
	// 防止添加空回调函数
	if callback == nil {
		return
	}
	
	// 创建新的回调节点
	newItem := &callbackCommon{handler, key, callback, nil}
	
	if t.first == nil {
		// 如果链表为空，新节点成为第一个节点
		t.first = newItem
	} else {
		// 否则将新节点添加到链表末尾
		t.last.next = newItem
	}
	// 更新最后一个节点的指针
	t.last = newItem
}

// Remove 从回调链表中移除指定的回调函数
// 参数说明:
//   - handler: 要移除的回调的处理器标识
//   - key: 要移除的回调的唯一标识键
// 注意: 如果找不到匹配的回调，此方法不会产生任何效果
func (t *callbacks) Remove(handler, key interface{}) {
	var prev *callbackCommon
	
	// 遍历链表查找要移除的节点
	for callback := t.first; callback != nil; prev, callback = callback, callback.next {
		// 找到匹配的节点
		if callback.handler == handler && callback.key == key {
			if t.first == callback {
				// 如果是第一个节点，更新 first 指针
				t.first = callback.next
			} else if prev != nil {
				// 如果是中间节点，更新前一个节点的 next 指针
				prev.next = callback.next
			}
			
			if t.last == callback {
				// 如果是最后一个节点，更新 last 指针
				t.last = prev
			}
			
			// 找到并移除后立即返回
			return
		}
	}
}

// Invoke 执行链表中所有注册的回调函数
// 按照添加的顺序依次执行每个回调
// 注意: 如果某个回调函数为 nil，会被跳过
func (t *callbacks) Invoke() {
	// 从头节点开始遍历整个链表
	for callback := t.first; callback != nil; callback = callback.next {
		// 确保回调函数不为 nil 再执行
		if callback.call != nil {
			callback.call()
		}
	}
}

// Count 返回链表中回调函数的数量
// 返回值: 当前注册的回调函数总数
func (t *callbacks) Count() int {
	var count int
	
	// 遍历链表计数
	for callback := t.first; callback != nil; callback = callback.next {
		count++
	}
	
	return count
}
