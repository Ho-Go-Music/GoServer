# HoGo Muisc Server
## 1. Introduction
base net/http, gin, gorm, mysql, redis, jwt, viper, logrus, go mod, docker, docker-compose, k8s, helm
## 2. Quick Start
## 3. third-party
#### [gob](https://pkg.go.dev/encoding/gob): 软件包gob管理gobs的流-在编码器（发射器）和解码器（接收器）之间交换的二进制值。
一个典型的用途是传输远程过程调用（RPC）的参数和结果，例如net/rpc提供的参数和结果。 该实现为流中的每种数据类型编译一个自定义编解码器，当使用单个编码器传输值流时，效率最高，从而摊销编译成本。
1. `gob.Register(time.Time{}) `这行代码的作用是在 gob 包的默认编码器和解码器中注册 time.Time 类型。这样，gob 包就能正确地序列化和反序列化 time.Time 类型的数据。  在你的代码中，你将 time.Time 类型的数据存储到了 Redis 中。Redis 本身并不知道如何处理 Go 语言的 time.Time 类型，所以你需要先将 time.Time 类型的数据转换（序列化）为可以被 Redis 存储的格式（比如字符串或者字节流），然后再从 Redis 中取出数据时，将其转换（反序列化）回 time.Time 类型。  这就是为什么你需要使用 gob 包，并且需要在 gob 包中注册 time.Time 类型。因为 gob 包可以帮助你完成 time.Time 类型数据的序列化和反序列化操作。

#### [go-redis](https://redis.io/docs/latest/develop/connect/clients/go/): go-redis是一个用于Go语言的Redis客户端。
#### [gorm](https://gorm.io/zh_CN/docs/index.html): GORM是一个适用于Go的ORM库，它的主要特性是支持多种数据库，例如MySQL、PostgreSQL、SQLite、SQL Server等。
#### [gin](https://gin-gonic.com/zh-cn/docs/): Gin是一个用Go（Golang）编写的Web框架。它具有类似于Martini的API，但性能更好。如果您需要性能和良好的生产质量，您会发现Gin非常有用。
#### [go-toml](https://github.com/pelletier/go-toml):解析toml配置文件