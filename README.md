# websocketSample
## 介绍
这是一个即时通知的例子，所有websocket连接订阅一个rabbitmq的queue，消息可以通过rabbitmq传送；使用非阻塞的epoll管理连接，当连接有消息处理的时候触发回调函数，支持百万级别的连接。一般的websocket例子，我们总会启2个goroutine，一个go read(),一个go write(),然后在里面for select等待；100w个连接就有200w的goroutine，每个2-8K,那就4G-16G内存。
## 运行 
这个例子中使用到rabbitmq,如果不需要，屏蔽掉main.go里面的 go  Manager.Start() 这一行代码再编译  
运行serser端：
go get github.com/gmluckly/websocketSample   
go build  
./websocketSample 

运行client端  
cd clientSample  
go build   
./clientSample -n 2000    2000表示请求连接的数量