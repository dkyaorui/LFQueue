package LFQueue

import (
    "fmt"
    "reflect"
    "runtime"
    "sync/atomic"
)

type LFNode struct {
    value interface{}
}

type LFQueue struct {
    capacity        uint64   // 容量 大小取为2^n， 方便位运算
    endIndex        uint64   // 数组结束索引
    writeCursor     uint64   // 写游标
    readCursor      uint64   // 读游标
    ringBuffer      []LFNode // 数据组
    availableBuffer []int    // 标记组，默认值-1
}

// 返回值：数据，错误
func (q *LFQueue) Pop() (interface{}, error) {
    _, next, err := q.getReadNext(1)
    if err != nil {
        return nil, err
    }
    q.availableBuffer[next&(q.endIndex)] = -1
    return q.ringBuffer[next&(q.endIndex)].value, nil

}

// 返回值：结束游标，错误
func (q *LFQueue) Push(value interface{}) error {
    next, err := q.getWriteNext(1)
    if err != nil {
        return err
    }
    q.availableBuffer[next&(q.endIndex)] = 1
    q.ringBuffer[next&(q.endIndex)] = LFNode{value: value}
    return nil
}

// 返回值： 数据，错误
func (q *LFQueue) PopMore(n uint64) ([]interface{}, error) {
    current, next, err := q.getReadNext(n)
    if err != nil {
        return nil, err
    }
    values := make([]interface{}, int(next-current))
    current++ // 游标后移一位开始读数据
    for index, _ := range values {
        i := current + uint64(index)
        values[index] = q.ringBuffer[i&(q.endIndex)]
        q.availableBuffer[i&(q.endIndex)] = -1
    }
    return values, nil
}

// 返回值： 结束游标，错误
func (q *LFQueue) PushMore(in interface{}) error {
    values := reflect.ValueOf(in)
    if values.Kind() != reflect.Array && values.Kind() != reflect.Slice {
        return fmt.Errorf("not array or slice")
    }
    var n uint64
    n = uint64(values.Len())
    next, err := q.getWriteNext(n)
    if err != nil {
        return err
    }
    current := next - (n - 1)
    // 向申请到的区间写入数据
    num := values.Len()
    for i := 0; i < num; i++ {
        q.availableBuffer[current&(q.endIndex)] = 1
        q.ringBuffer[current&(q.endIndex)] = LFNode{value: values.Index(i).Interface()}
        current++
    }
    return nil
}

/*
获取下一个写入范围结束端点

n: 申请可写范围

next: 结束点
err: 错误
*/
func (q *LFQueue) getWriteNext(n uint64) (next uint64, err error) {
    if n < 1 {
        n = 1
    }
    var current uint64

    for {
        current = q.writeCursor
        next = current + n
        // 如果申请的空间已被写入或者队列当前游标和申请的开始不同则等待
        if q.checkAvailableCapacity(current, n) && atomic.CompareAndSwapUint64(&q.writeCursor, current, next&q.endIndex) {
            break
        }
        runtime.Gosched()
    }
    return next, nil
}

// 检查当前游标开始n空间内是否都可写
func (q *LFQueue) checkAvailableCapacity(current uint64, n uint64) bool {
    // 申请的空间都未被标记时才可写入，未被标记时值为默认值-1
    end := current + n
    current++
    for current <= end {
        if q.availableBuffer[current&(q.endIndex)] != -1 {
            return false
        }
        current++
    }
    return true
}

/*
获取当前游标开始n内可读空间

n: 申请可读范围

next: 结束点
num: 实际可读空间
err: 错误
*/
func (q *LFQueue) getReadNext(n uint64) (start uint64, next uint64, err error) {
    if n < 1 {
        n = 1
    }
    var current uint64

    for {
        current = q.readCursor
        if q.availableBuffer[(current+1)&(q.endIndex)] == -1 {
            return 0, 0, fmt.Errorf("there is no data can read")
        }
        next = q.checkAvailableRead(current, n)
        if atomic.CompareAndSwapUint64(&q.readCursor, current, next&q.endIndex) {
            break
        }
        runtime.Gosched()
    }
    return current, next, nil
}

// 返回n以内最长可读空间
func (q *LFQueue) checkAvailableRead(current uint64, n uint64) uint64 {
    end := current + n
    current++

    for current <= end {
        index := current & (q.endIndex)
        if q.availableBuffer[index] == -1 {
            return current
        }
        current++
    }
    return end
}

// 获取最近的2的指数
func getCapacity(in int) int {
    in--
    in |= in >> 1
    in |= in >> 2
    in |= in >> 4
    in |= in >> 8
    in |= in >> 16
    in++
    return in
}

func NewQue(inCapacity int) *LFQueue {
    capacity := uint64(getCapacity(inCapacity))
    que := LFQueue{
        capacity:        capacity,
        endIndex:        capacity - 1,
        writeCursor:     0,
        readCursor:      0,
        ringBuffer:      make([]LFNode, capacity),
        availableBuffer: make([]int, capacity),
    }
    for index, _ := range que.availableBuffer {
        que.availableBuffer[index] = -1
    }
    return &que
}
