
### 生成rpc代码

```shell
goctl rpc protoc ./pb/user.proto --go_out=./pb/ --go-grpc_out=./pb/ --zrpc_out=. --style=go_zero
```