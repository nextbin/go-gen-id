# go-id-gen
一个生成全局唯一ID的服务，支持 HTTP、RPC 调用。

## 理论能力
- 支持最多1024个机器
- 每毫秒每个机器生成4096个ID
- 时间支持69年

## 开发环境
- Golang: 1.13.1
- OS: Mac OS 10.14
- IDE: GoLand 2018.3.5

## 部署
1. 安装 Golang 及依赖包
2. 修改配置
    1. cofing.go: MachineId（第一个服务机器可以不操作）
    2. cofing.go: RedisAddr（可选，不使用Redis检查MachineId机制可以不操作）
    3. cofing.go: MysqlDataSourceNaming（可选，不使用MySQL检查MachineId机制可以不操作）
3. MySQL建表（可选，不使用MySQL检查MachineId机制可以不操作）
4. 启动程序

## 依赖包

1. HTTP: github.com/gin-gonic/gin
2. MySQL: github.com/go-sql-driver/mysql
3. Redis: github.com/gomodule/redigo
3. 日志: github.com/sirupsen/logrus

## TODO

- [ ] 支持 RPC 调用
- [ ] 使用包管理工具
- [ ] IP白名单
- [ ] 支持更多的 MachineId 检查方式（Redis-Sentinel、Mongo）

## 参考资料

> https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake.html
> 
> https://github.com/twitter-archive/snowflake/tree/snowflake-2010
>
> https://github.com/gin-gonic/gin
> 
> https://github.com/go-sql-driver/mysql
> 
> https://github.com/gomodule/redigo
>
> https://github.com/sirupsen/logrus