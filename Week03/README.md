学习笔记

1. CSP模型(communicating sequential processes) 顺序通信过程

    多个goroutine可以通过管道(channel)传输消息。
    
    在Java中，通过共享(资源)内存的形式进行线程间通信，在GO中提供了另外一种通信方式，就是CSP

2. 并发模型

    GO调度器不是抢占式调度器，而是协作式调度器。

    [GO并发模型](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html)
    
    ![image.png](https://i.loli.net/2020/12/06/tQ3dFNZzMsrJbqK.png)
    
    - 逻辑处理器(P)，每个虚拟内核提供一个逻辑处理器
    - OS Thread(M)，每一个P被分配一个M，OS将M放置在核心中处理逻辑
    - goroutines(G)，OS线程和内核的关系 类似于 goroutine和OS线程
    - GRQ(全局运行队列)
    - LRQ(本地运行队列)，每一个P分配一个LRQ，管理那些分配给这个P执行的G，M会对这些G依次进行上下文切换
    
    M与G是一个N:M的模型
    
    1. 异步网络系统调用
    
        ![image.png](https://i.loli.net/2020/12/06/seLaPxmJKNBtGod.png)
        
        当M上绑定的G1需要进行异步网络系统调度时，该G1将会移动到网络轮询器，然后LRQ上的G2会进行上下文切换，绑定到M上处理逻辑.
        
        防止G1进行网络系统调度时会阻塞M
        
        ![image.png](https://i.loli.net/2020/12/06/i8lN3f5OMaCKPU1.png)
        
        当G1完成系统调度时，G1会重新放入LRQ队尾
    
    2. 同步系统调用
       
        ![image.png](https://i.loli.net/2020/12/06/tdH8oTjN1VwSs2p.png)
        
        G在执行同步系统调用时，会阻塞M，此时G1将M1阻塞，在M1上绑定的P将会重新绑定到新的M2，然后由M2进行调度LRQ中的G
        G1执行完之后会重新放入LRQ队尾。
    
    3. 工作窃取(working-stealing)
       
        ```go
        runtime.schedule() {
           // only 1/61 of the time, check the global runnable queue for a G.
           // if not found, check the local queue.
           // if not found,
           //     try to steal from other Ps.
           //     if not, check the global runnable queue.
           //     if not found, poll network.
        }
        ```

3. 方法设计时，要注意函数执行的时间是否会过长；是否需要使用异步的形式获取结果；如果异步过程中符合条件的结果已经出现，如何终止方法执行
    
    ```go
    // 全量获取目录列表，如果目录过多，耗费时间会很长
    func ListDirectory(dir string) ([]string, error)
    
    // 通过channel异步获取目录，并放入channel中
    func ListDirectory(dir string) chan string
   
    // 可以通过一个方法回调判断是否满足跳出的条件
    func ListDirectory(dir string, fn func(string) bool) chan string
    ```
    
    filepath.WalkDir 也是类似的模型，如果函数启动 goroutine，则必须向调用方提供显式停止该goroutine 的方法。通常，将异步执行函数的决定权交给该函数的调用方通常更容易。
    
4. 使用goroutine时我们必须考虑两个问题
    - When will it terminate? 什么时候终止
    - What could make it terminate? 什么能让它终止
    
    即控制goroutine的整个生命周期。
    
    1. 使用channel控制goroutine的关闭
        ```go
        // 使用stop、done两个channel控制goroutine的关闭
        func (g *Group) Run() error {
           if len(g.fns) == 0 {
               return nil
           }
           stop := make(chan struct{})
           done := make(chan error, len(g.fns))
           for _, fn := range g.fns {
               go func(fn Run) {
                   // 在fn中使用 <-stop来控制资源的释放和关闭
                   done <- fn(stop)
               }(fn)
           }
           var err error
           for i := 0; i < cap(done); i++ {
               if i == 0 {
                   // 返回第一个error
                   err = <-done
                   close(stop)
               } else {
                   <-done
               }
           }
           close(done)
           return err
        }
        ```
    2. 使用WaitGroup和channel对goroutine进行控制
        WaitGroup可以控制一组goroutine执行完毕之后，在处理剩余逻辑。
        ```go
        func PoolExecute(int n) {
           ch := make(chan bool, n)
           var g sync.WaitGroup
           for i := 0; i < n; i++ {
               g.Add(1)
               ch <- true
               // 启动协程处理业务
               go func() {
                   // 在最后使WaitGroup减1
                   defer g.Done()
                   // DoSomething
                   //使用channel控制最多只能有n个协程
                   <-ch
               }
           }
           g.Wait()
        }
        ```
    3. 使用超时对goroutine的执行时间进行控制

5. 内存模型
    
    1. Happen-Before
    
        可见性，a操作的结果对b操作是可见的，说明a Happen-Before b
        
    2. Memory-Recording
    
        **内存重排序**
        
        为了提高内存的读写效率，减少程序指令数，最大化的提高CPU利用率，CPU会对指令重排序。编译器也会进行重排序
        
        ```go
        // goroutine1
        func run1() {
           a = 1 // (1)
           fmt.Println(b) // (2)
        }
        // goroutine2
        func run2() {
           b = 1
           fmt.Println(a)
        }
        // 这样的结果有可能会出现0 0，这里主要是因为可能会出现CPU重排序。
        ```
        
        对于(1)、(2)操作，在协程中，(2)的执行是不需要依赖(1)的，所以两个操作完全可以并行。
        
        ![内存模型](https://ss.csdn.net/p?https://mmbiz.qpic.cn/mmbiz_png/ASQrEXvmx62Cvw3EzBCJ5VBpV3E1jgC0g5gyqtznicpHMKP06LMQRufpTicjAiazJp7dxDC3cSs1icpibQAtwTEEd8A/640?wx_fmt=png)
        
        当操作(1)执行时，CPU会将A的值存在CPU中的*StoreBuffer*中，此时CPU可以继续运行(2)，之后*StoreBuffer*中的数据会被逐级写入下级缓存中，
        在L3 Cache中就可以被其他线程看到。*StoreBuffer*隐藏了(1)写数据的耗时。
        
        在多线程的操作中
        
        ![出现00的情况](https://ss.csdn.net/p?https://mmbiz.qpic.cn/mmbiz_png/ASQrEXvmx62Cvw3EzBCJ5VBpV3E1jgC0IWZ7zqlmoq3PXibxLyibLFaW63HFUJLy8B7jKYpRQacR5pl1PpcELbPw/640?wx_fmt=png)
        
        当两个线程的写指令还未从*StoreBuffer*写到内存中，两个线程优先执行了读操作，就会出现00的情况。
        这种情况就是**指令的重排序**。
        
    3. memory barrier
    
        对于多线程的程序，所有的 CPU 都会提供"锁"支持，称之为barrier，或者fence。
        
        内存屏障，barrier指令开始之后，所有对内存的操作都会在*刷新缓存*之后发生。
        
        `刷新缓存：将CPU的StoreBuffer中的数据刷新到内存中`
        
        ![缓存结构](https://imgconvert.csdnimg.cn/aHR0cHM6Ly9pbWFnZXMuY25ibG9ncy5jb20vY25ibG9nc19jb20vbGlsb2tlLzIwMTExMS8yMDExMTEyMDA0MzEzNDQ4ODMucG5n?x-oss-process=image/format,png)
        
        [缓存结构](https://blog.csdn.net/qq_21125183/article/details/80590934)
        
        1. cache coherence 缓存一致性
            
            
            
6. sync包

    1. Once 双重检查锁，执行且仅执行一次f函数
        
7. context包

    **Request-Scope**
    
    - 在服务器请求的生命周期中的function链应该传递Context，例如当前用户的信息、授权令牌。
    - 当一个请求被取消或者超时，goroutine内部启动的所有资源都应该被回收
    - 内部启动的新的goroutine应该快速退出。
    
    **Context种类**
    
        可以选择性的使用WithCancel、WithDeadline、WithTimeout、WithValue等包装上下文
    
    - WithCancel
        
        当context被cancel之后，会将cancel的信号（close done chan）层层递归传递给所有派生的context，从而让整个调用链中的goroutine退出。
        
    - WithTimeout
    
        在发出cancel信号之前加一个超时的逻辑。WithTimeout(parent, timeout)等价于WithDeadline(parent, time.Now().Add(timeout))。
        
        ```go
        t := time.Second * 10
        ctx, cancelFunc := context.WithTimeout(context.Background(), t)
        defer cancelFunc()
        select {
        case <-time.After(1 * time.Second):
            fmt.Println("执行结束")
        case <-ctx.Done():
            fmt.Println("超时")
        }
        ```
        
    - WithDeadline
    
        可以应用于资源的回收，如还剩下多少时间时不足以执行任务，直接回收资源并结束等。
        
    - WithValue
        
        存储Key-Value键值对，存储时，要避免存储业务相关的信息。go允许在多个goroutine中传递context，所以Value方法获取值是线程安全的，应该尽量保证它里面的值是不可变的。
        
        调用function (c *valueCtx) Value(key interface{}) interface{}获取时，go会依次递归获取value值，从child到parent，Background和TODO会返回nil。

        在context中存储map，在多个goroutine中修改、读取map中的值有可能会出现data race。通常读多写少，我们会使用COW的技巧，在新建Context的之前，先使用COW创建map的副本。各个Context中的数据不会被其他Context污染。
    
    **Context接口**
        
        ```go
        type Context interface {
           // 返回可以被取消的上下文-任务被取消时时间 如果可以被取消，返回true
           Deadline() (deadline time.Time, ok bool)
           // 上下文被主动取消、超时时，会关闭channel
           //   如果context不能被取消，则返回nil
           Done() <-chan struct{}
           // 如果Done没有被close，返回nil
           // 如果Context被取消之后，返回Canceled
           // 如果Context超时，返回DeadlineExceeded（这里的error是一个结构体，而不是指针）
           Err() error
           // 存储请求层面的key-value
           Value(key interface{}) interface{}
        }
        ```
        
        

## 未完待续


1. goroutine生命周期要清楚，避免goroutine泄露
    通过context控制超时
    通过channel发送消息
    通过close channel
    
2. 启动goroutine交给调用方

3. https://golang.org/ref/mem
