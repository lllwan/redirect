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

### 三、基本功能

+ 隐性转发
    + 支持自定义百度统计或者51la统计。
    + 支持自定义title、description、keyword、favicon
    
+ 301 重定向
    + 支持直接重定向和全站重定向
    
### 四、 使用

+ 运行本系统，直接监听80端口或者配置nginx的proxy_pass。
+ 使用api或者直接修改数据库的acl表:

|参数|类型|是否必须|描述|
|---|---|---|---|
|domain|string|是|域名|
|url|string|是|转发的目标地址|
|method|string|是|hide：隐性转发，301：重定向，301all: 全站重定向|
|title|string|否|页面的title（只有在隐性转发的情况下有效）|
|keywords|string|否|关键字（只有在隐性转发的情况下有效）|
|description|string|否|描述（只有在隐性转发的情况下有效）|
|favicon|string|否|站点图标的url只有在隐性转发的情况下有效）|
|count|string|否|统计功能值只支持：baidu或者51la（只有在隐性转发的情况下有效）|
|countid|string|否|统计方案的ID（只有在隐性转发的情况下有效）|
|username|string|否|暂时可忽略，为多用户准备。|

### 五、操作实例
+ 暂缺

    
