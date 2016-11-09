# lib4go

对常用功能进行封装提供友好的公共库代码，供其它golang程序调用

基本库
* 线程安全的map
* 数据库操作（oralce,sqlite）
* 编码处理（gbk,utf-8字符串转换,base64,hex）
* 加解密（md5,des,rsa,sha1等）
* 图片库（图片，文字绘制）
* 系统资源（获取CPU，内存，硬盘使用情况）

其它三方组件封装

* 日志（写入文件）
* influxdb操作（存，取）
* MQ操作（发布，订阅等）
* scheduler（基于cron的任务处理）
* LUA脚本引擎
* HTTP请求(支持证书,cookie等)
* HTTP Server
