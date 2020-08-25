# websocketSample
这是一个即时通知的例子，所有websocket订阅一个queue,消息可以通过rabbitmq，支持百万级别的连接
一般的websocket例子，我们总会启2个goroutine，一个go read(),一个go write(),然后在里面for select等待；100w个连接就有200w的goroutine，每个2-8K,那就4G-16G内存