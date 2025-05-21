## TCP 消息格式
```
\r\n + uint32消息类型 + uint32消息体长度 + 消息体
```

### 消息类型
- 1: 心跳包
- 2: 指令包


## 消息结构

不同的消息类型，对应不同结构：

1. hearbeat
{
    "name":"csmm-1",
    "time":"2022-07-13 17:03:06.26",
    "stat": "idle"
}

2. 手机向服务器取任务
null

3. 服务器应答任务内容给手机
{
    "id": "1001",
    "type": "shell",
    "data": "ls /"
}

4. 手机向服务器发任务应答
{
    "id": "1003",
    "type": "shell",
    "stat": "running",
    "resp": {
        "type": "raw_frame",
        "width": 1080,
        "height": 2400,
        "bits_per_pixel": 24,
        "bytes":"base64"
    }
}


## mysql

部署之前,服务器时区要设置成东8

```
CREATE SCHEMA `pcdn` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

```


```
curl -v -H "X-API: user/register" -d \
'{
	"email": "demo@moli.bot",
	"password": "123456"
}' "http://127.0.0.1:10000/usercenter"
```

curl -v -H "X-API: user/register" -d \
'{
	"email": "user001@moli.bot",
	"password": "123456"
}' "http://127.0.0.1:10000/usercenter"



