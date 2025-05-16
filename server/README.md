

## mysql

部署之前,服务器时区要设置成东8

```
CREATE SCHEMA `arbitrage` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

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



