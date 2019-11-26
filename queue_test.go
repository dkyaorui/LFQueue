package LFQueue

import (
    "fmt"
    "sync"
    "testing"
)

var inCapacity = 1024

func TestNewQue(t *testing.T) {
    que := NewQue(inCapacity)
    capacity := uint64(getCapacity(inCapacity))
    if que.capacity != capacity {
        t.Error("que's capacity is wrong")
    }
    if que.endIndex != capacity-1 {
        t.Error("que's endIndex is wrong")
    }
    if que.writeCursor != 0 {
        t.Error("que's writeCursor is wrong")
    }
    if que.readCursor != 0 {
        t.Error("que's readCursor is wrong")
    }
    if uint64(len(que.ringBuffer)) != capacity {
        t.Error("que's len of ringBuffer is wrong")
    }
    for _, value := range que.availableBuffer {
        if value != -1 {
            t.Error("que's value of available is not -1")
        }
    }
}

func TestLFQueue_Push(t *testing.T) {
    que := NewQue(inCapacity)
    capacity := uint64(getCapacity(inCapacity))
    var wg sync.WaitGroup
    for item := uint64(1); item <= capacity+2; item++ {
        wg.Add(1)
        go func(i uint64) {
            defer wg.Done()
            var err *QueError
            err = que.Push(i)
            if err != nil {
                if err.StatusCode == QueueIsFull {
                    fmt.Printf("err: %s, value: %d\n", err, i)
                } else {
                    t.Errorf("que's push method run with wrong:%+v", err)
                }
            }
        }(item)
    }
    wg.Wait()
    checkMap := make(map[uint64]int)
    for index, value := range que.ringBuffer {
        if value.value == nil {
            continue
        }
        if checkMap[value.value.(uint64)] == 0 {
            checkMap[value.value.(uint64)] = index
        } else {
            t.Errorf("data error occurs!index:%d, value:%d", index, value.value.(uint64))
        }
    }
}

func TestLFQueue_Pop(t *testing.T) {
    que := NewQue(inCapacity)
    var wg sync.WaitGroup
    item := 1
    for ; item <= inCapacity; item++ {
        wg.Add(1)
        go func(item int) {
            defer wg.Done()
            err := que.Push(item)
            if err != nil {
                t.Errorf("que's push method run with wrong:%s", err)
            }
        }(item)
    }
    wg.Wait()
    var checkMap sync.Map
    index := 0
    flag := true
    for flag {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            value, err := que.Pop()
            if err != nil {
                if err.StatusCode == QueueIsEmpty {
                    flag = false
                    return
                }
            }
            _, ok := checkMap.Load(value.(int))
            if ok {
                t.Errorf("data error occurs!index:%d, value:%d", i, value.(int))
            } else {
                checkMap.Store(value.(int), 1)
            }
        }(index)
        index++
    }
    wg.Wait()
}

func BenchmarkLFQueue_Push(b *testing.B) {
    var wg sync.WaitGroup
    que := NewQue(b.N)
    capacity := uint64(b.N)
    for item := uint64(1); item <= capacity; item++ {
        wg.Add(1)
        go func(i uint64) {
            defer wg.Done()
            var err *QueError
            err = que.Push(i)
            if err != nil {
                if err.StatusCode == QueueIsFull {
                    fmt.Printf("err: %s, value: %d\n", err, i)
                } else {
                    b.Errorf("que's push method run with wrong:%+v", err)
                }
            }
        }(item)
    }
    wg.Wait()
}

func BenchmarkLFQueue_Pop(b *testing.B) {
    b.StopTimer()
    que := NewQue(b.N)
    var wg sync.WaitGroup
    item := 1
    for ; item <= b.N; item++ {
        wg.Add(1)
        go func(item int) {
            defer wg.Done()
            err := que.Push(item)
            if err != nil {
                b.Errorf("que's push method run with wrong:%s", err)
            }
        }(item)
    }
    wg.Wait()
    b.StartTimer()
    for index := 0; index < b.N; index++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            _, err := que.Pop()
            if err != nil {
                if err.StatusCode != QueueIsEmpty {
                    b.Errorf("error: %+v", err)
                }
            }
        }(index)
    }
    wg.Wait()
}
