# Daenerys

叫Daenerys是因为她有三条龙

------

## 项目概述
* go语言实现的web服务器，具有注册登录逻辑
* 前端模仿京东，数据库数据从京东商品爬取


### 爬虫
* go爬虫抓取[京东商品分类页](https://www.jd.com/allSort.aspx),存储到mysql表中

### 数据库
* 采用MySQL，利用InnoDB引擎提供更高粒度的锁和更好的并发性能
* 表的结构
* 建表语句

### 服务端
* golang，暂不借助其他框架

### 注册登录功能
- [x] session/cookie实现登录验证
- [ ] redis实现登录验证


### 前端页面
* 常规html+bootstrap+js实现
* 注册、登录、登陆失败、登录成功四个页面
