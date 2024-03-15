properties conf

properties to stream

properties test

connect to stream

dial tcp

发数据要做缓冲区, 这样一个字段一个字段的发太累人

处理读取时报错的问题

发送 disconnect

要做一个 onConnect 的回调 (要考虑成功和失败是否分成两个函数 onConnectFail, 或是用参数)