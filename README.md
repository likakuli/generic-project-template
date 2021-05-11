[![Build Status](https://travis-ci.com/likakuli/generic-project-template.svg?branch=master)](https://travis-ci.com/likakuli/generic-project-template)
## 通用Restful API项目模板

欢迎使用，这是一个用Go编写的简单通用的Restful API项目，遵循SOLID原则。

部分灵感来自于 [service-pattern-go](https://github.com/irahardianto/service-pattern-go)

### 依赖

* [Gin](https://github.com/gin-gonic/gin)
* [Gorm](https://github.com/go-gorm/gorm)
* [Testify (Test & Mock framework)](https://github.com/stretchr/testify)
* [Mockery (Mock generator)](https://github.com/vektra/mockery)
* [Hystrix-Go (Circuit Breaker)](https://github.com/afex/hystrix-go)

### 开始

* [安装](#安装)
* [介绍](#介绍)
* [目录结构](#目录结构)
* [Mocking](#mocking)
* [Tesing](#testing)
* [能力支持](#能力支持)
* [TODO](#TODO)

### 安装

**克隆项目代码**

```shell
git clone https://github.com/likakuli/generic-project-template.git
```

**启动mysql服务并初始化数据库**

```shell
docker run --name=mysql -it -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d registry.cn-beijing.aliyuncs.com/likakuli/mysql
```

**注意**：如果是在MacOS上使用Docker for Mac启动的容器，则需要安装 [docker-connector](https://github.com/wenjunxiao/mac-docker-connector) ，否则无法在本机通过容器IP访问容器，因为参考[这里](https://docs.docker.com/docker-for-mac/networking/#there-is-no-docker0-bridge-on-macos)。安装命令如下

```shell
# 安装 docker-connector
brew install wenjunxiao/brew/docker-connector  
# 把 docker 的所有 bridge 网络都添加到路由中
docker network ls --filter driver=bridge --format "{{.ID}}" | xargs docker network inspect --format "route {{range .IPAM.Config}}{{.Subnet}}{{end}}" >> /usr/local/etc/docker-connector.conf  
# 启动服务
sudo brew services start docker-connector  
# 在 docker 端运行 wenjunxiao/mac-docker-connector，需要使用 host 网络，并且允许 NET_ADMIN
docker run -it -d --restart always --net host --cap-add NET_ADMIN --name connector wenjunxiao/mac-docker-connector
```

镜像涉及到的Dockerfile与sql放置在`docker`文件夹下

**运行单元测试**

```shell
make test
```

**编译程序**

```shell
make build
```

**运行程序**

```shell
# 替换配置文件MySQL connectionString
./docker/replace_ip.sh
# 启动程序
./generic-project-template --config=./conf/config.toml --log_dir=./logs --v=1
```

**访问程序**

```shell
curl http://localhost:8080/api/v1/score/Lucy/vs/Lily
```

### 介绍

这是一个简单通用的Restful API项目，内置依赖注入、Mocking等功能，旨在方便快速的编写安全可靠的Restful API代码。不同的数据结构之间通过接口来访问，避免直接引用具体的实现，这样就可以实现依赖注入及采用Mock结构进行单元测试的效果。

举例来说：

`IPlayerServie --> IPlayerRepository`

```go
type PlayerController struct {
	service interfaces.IPlayerService
}

type playerService struct {
	repo interfaces.IPlayerRepository
}
```
### 目录结构

```text
/cmd
  /- apiserver
/conf
/pkg
  /- config
  /- controllers
  /- interfaces
  /- models
  /- repositories
  /- server
    /- middlewares
    |- container.go
    |- router.go
    |- server.go
  /- services
  /- viewmodels
```
**controllers**

控制器文件夹下包含所有`Gin Route Handler`，里面只包含处理`Request`和`Response`的逻辑，不包含任何业务逻辑和数据访问逻辑。仅依赖于`interfaces`下的`IService`接口，不依赖于具体实现。

**interafces**

接口文件夹下存放所有`IService`和`IRepository`接口定义及通过`Mockery`自动生成的用于单元测试的文件，不包含具体接口实现。

**models**

模型文件下下存放所有与数据库映射的实体模型对应的`Go Struct`，只包含数据结构，不包含数据访问逻辑。可以由 [gen](https://github.com/smallnest/gen) 根据数据库表结构自动生成，详情参考[这里](https://mp.weixin.qq.com/s/J7NO_kybMtatpWCnghu6Ag )

**repositories**

仓库文件夹下存放所有数据库访问逻辑，且实现了`interfaces`下定义的`IRepository`接口，主要用到`models`文件夹下定义的实体结构。

**services**

服务文件夹下存放所有实现了`services`下定义的`IService`接口的逻辑，共`controllers`直接使用。其中涉及到的数据库访问部分均通过调用`interfaces`下的`IRepository`接口实现，不依赖任何具体实现。

**viewmodels**

视图模型文件夹下存放所有需要与API交互的实体，主要包含从API获取到的结构和返回值得结构。与`models`的区别在于前者对应api层，后者对应数据库层。

**router**

路由文件夹下包含了所有可以对外提供服务的`Restful API`的路由注册逻辑。

**container**

容器文件下包含了所有依赖注入需要的`Provider`的逻辑，且在此选择具体使用的接口实现类型。

### Mocking

为方便进行单元测试，使用`Mockery`自动`interfaces`下接口实现，例如生成`IPlayerService的实现`，只需要进入`interfaces`文件夹下执行如下命令即可，最后会在`interfaces`下自动创建`mocks`文件夹来存放自动生成的文件。

```shell
mockery -name=IPlayerService
```

需要提前安装`mokery`工具

### Testing

有了依赖注入和`Mock`功能后，就可以针对任意接口实现编写单元测试了，示例中添加了针对`services`he`controllers`的单测，供参考。

### 能力支持

- [x] Tracing
- [x] PProf
- [x] Prometheus Metrics
- [x] Health Check
- [x] Mock
- [x] Testing
- [x] Circuit Breaker
- [x] Rate Limit
- [x] Common [go-utils](https://github.com/leopoldxx/go-utils)
- [ ] ...