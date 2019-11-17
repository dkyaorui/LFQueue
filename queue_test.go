package LFQueue

import "testing"

func TestNewQue(t *testing.T) {
    capacity := uint64(102400)
    que := NewQue(capacity)
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

}

func TestLFQueue_Pop(t *testing.T) {

}

func TestLFQueue_PushMore(t *testing.T) {

}

func TestLFQueue_PopMore(t *testing.T) {

}
