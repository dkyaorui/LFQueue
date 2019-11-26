# LFQueue

无锁队列，使用原子操作避免并发冲突。

思路：[高性能队列——Disruptor](https://tech.meituan.com/2016/11/18/disruptor.html)

version: 0.1.1

数据正常：已测试

性能:

`Pop()`

> goos: darwin
>
> goarch: amd64
>
> pkg: LFQueue
>
> BenchmarkLFQueue_Pop-8   	 4192549	       266 ns/op
> 
> PASS 

`Push()`
>goos: darwin
>  
>goarch: amd64
>
>pkg: LFQueue
>
>BenchmarkLFQueue_Push-8   	 4164613	       271 ns/op
>
>PASS
