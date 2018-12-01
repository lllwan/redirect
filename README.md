### URl 转发系统

> 一个简单的url转发系统demo，基于国内运营商封锁80、443端口和url隐性转发要求备案或收费的现状开发。目前放在香港VPS上，给自己和朋友用，得益与golang的特性，在512内存的VPS上运行稳定。并发无压力。
> 目前只能使用API的方式管理。API使用JWT方式认证，还比较简陋，只实现了添加、删除转发规则。修改密码的接口，其实不如直接修改数据库方便^_^。几个接口只是为了预留添加前端界面的坑。

### 一、默认配置

+ 默认用户名密码都是admin
+ 默认配置文件config.yaml：
```yaml
HTTP_BIND: ":8899"                  # 监听端口
DATABASE: "redirect.db"             # sqlite数据库文件名称
SECRET: "6RyC2VpehJERy78Q"          # token加密秘钥。
Duration: 24                        # token有效期， 默认24小时
```
+ 使用sqlite数据库，默认数据库名：redirect.db, 首次运行会创建db并进行初始化。

### 二、认证
+ 使用用户名密码获取token，使用token做管理操作。

### 二、基本功能

+ 隐性转发
    + 支持自定义百度统计或者51la统计。
    + 支持自定义title、description、keyword、favicon
    
+ 301 重定向
    + 支持直接重定向和全站重定向
    
