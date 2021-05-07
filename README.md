## GenericProjectTemplate

编译：  
`make`

运行：

`genric-project-template --config=./conf/config.toml`   
其中 `--config` 为项目配置文件参数路径，其他参数为glog自带的参数


项目代码结构：

* cmd: 命令入口，可包含多个命令
* config: 配置参数
* pkg:   
    * api: 路由注册、参数校验等
    * dao: 数据库访问，所有需要访问数据库的逻辑都在这里
    * service: 连接api和dao层，封装业务逻辑，提供开箱即用的服务
    * model: 对象实体，与数据库交互

**其中model和dao部分由gen自动生成，gen的可执行文件在gen文件夹下，目前放了linux版和mac版，
命令如下（--connstr按需替换即可）：**  
``
gen --sqltype=mysql --connstr "root:123456@tcp(172.17.8.101:3306)/modelarts?charset=utf8mb4&parseTime=True&loc=Local" --database modelarts --gorm --module "github.com/likakuli/generic-project-template/pkg" --overwrite --template=./templates --generate-dao
``  
关于gen的使用教程可以参考[这里](https://mp.weixin.qq.com/s/J7NO_kybMtatpWCnghu6Ag)

当前能力支持：

* api调用详细信息trace跟踪
* /healthz: 健康检查，成功返回200 ok
* /metrics: prometheus metrics提供，包括基础指标和api调用等相关指标
* pprof: 支持pprof
* 全局接口限流

待支持能力（添加中）：  

* 完善单元测试
* dao和service部分接口封装
* ...
  