学习笔记

1. [项目布局](https://github.com/golang-standards/project-layout/blob/master/README_zh.md)
    
    1. 标准Go项目布局
        
        - /cmd
        
            项目的主干，一般存放可执行文件。一般使用项目名作为二级目录名，如/cmd/myapp。负责应用的启动、关闭、配置初始化等。
        
        - /internal
        
            存放私有程序和库代码（不希望他人在其他程序或库代码中导入），由Go编译器本身执行不可被导入。
            
        - /pkg
        
            存放外部应用可以使用的库代码。/internal/pkg一般用于项目中的跨多个应用的公共共享代码，但其作用域仅在当前项目工程中。
            
        - /api
        
            API协议定义目录，如xxxapi.proto的protobuf文件，以及生成的go文件。
        
        - /configs
        
            配置文件等。
        
        - /test
        
            测试文件。
        
    2. Kit Go项目布局
        
        
        
2. 微服务服务类型
    
    - interface
    
        对外的BFF服务，对前端暴露HTTP/gRPC接口。
    
    - service
    
        对内的微服务。
    
    - admin
    
        区别于service，数据权限更高的service。
        
    - job
    
        流式任务处理的服务，上游一般依赖 message broker。
    
    - task
    
        定时任务，类似 cronjob，部署到 task 托管平台中。
        
        
3. 应用的Lifecycle

    应用程序的初始化、执行、销毁。



####给自己的作业

1. 写一个小工具用于生成自己的文件目录