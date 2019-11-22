package LFQueue

import (
    "fmt"
    "strings"
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
            var err error
            err = que.Push(i)
            if err != nil {
                if strings.Contains(err.Error(), "the queue is full"){
                    fmt.Printf("err: %s, value: %d\n", err, i)
                }else {
                    t.Errorf("que's push method run with wrong:%+v", err)
                }
            }
        }(item)
    }
    wg.Wait()
    fmt.Println(que.ringBuffer)
    checkMap:= make(map[uint64]int)
    for index, value := range que.ringBuffer {
        if value.value == nil {
            continue
        }
        if checkMap[value.value.(uint64)] == 0 {
            checkMap[value.value.(uint64)] = index
        }else{
            t.Errorf("data error occurs!index:%d, value:%d", index, value.value.(uint64))
        }
    }
}

func TestLFQueue_Pop(t *testing.T) {
    que := NewQue(inCapacity)
    var wg sync.WaitGroup
    item := uint64(1)
    for ; item <= 1024; item++ {
        wg.Add(1)
        go func(item uint64) {
            defer wg.Done()
            var err error
            err = que.Push(item)
            if err != nil {
                t.Errorf("que's push method run with wrong:%+v", err)
            }
        }(item)
    }
    wg.Wait()
    var checkMap sync.Map
    index := 0
    flag := true
    for flag{
        wg.Add(1)
       go func(i int) {
           defer wg.Done()
           value, err := que.Pop()
           if err != nil {
               if strings.Contains(err.Error(), "no data") {
                   flag = false
                   return
               }
           }
           _, ok := checkMap.Load(value.(uint64))
           if ok {
               t.Errorf("data error occurs!index:%d, value:%d", i, value.(uint64))
           }else {
               checkMap.Store(value.(uint64), 1)
           }
       }(index)
       index++
    }
    wg.Wait()
}
