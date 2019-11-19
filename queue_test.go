package LFQueue

import (
    "fmt"
    "strings"
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
    item := uint64(1)
    for ; item <= capacity * 2; item++ {
        go func() {
            var err error
            err = que.Push(item)
            if err != nil {
                t.Errorf("que's push method run with wrong:%+v", err)
            }
        }()
    }
    checkMap:= make(map[uint64]int)
    for index, value := range que.ringBuffer {
        if checkMap[value.value.(uint64)] == 0 {
            checkMap[value.value.(uint64)] = index
            fmt.Println(index, value)
        }else{
            t.Errorf("data error occurs!index:%d, value:%d", index, value.value.(uint64))
        }
    }
}

func TestLFQueue_Pop(t *testing.T) {
    que := NewQue(inCapacity)
    item := uint64(1)
    for ; item <= 1024; item++ {
            var err error
            err = que.Push(item)
            if err != nil {
                t.Errorf("que's push method run with wrong:%+v", err)
            }
    }
    checkMap:= make(map[uint64]int)
    index := 0
    for {
        value, err := que.Pop()
        if err != nil {
            if strings.Contains(err.Error(), "no data") {
                break
            }
        }
        if checkMap[value.(uint64)] == 0 {
            checkMap[value.(uint64)] = index
            fmt.Println(index, value)
        }else{
            t.Errorf("data error occurs!index:%d, value:%d", index, value.(uint64))
        }
        index++
    }
}
