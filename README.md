# Daenerys


## 项目概述
* go语言实现的web服务器，具有注册登录逻辑
* 数据库数据从京东商品爬取
* 实例网页可在[Daenerys](http://154.8.143.128:18080)访问

### 项目体系    
![](root/static/struct2.png)


### 爬虫 jdSpider
* go爬虫抓取[京东商品分类页](https://www.jd.com/allSort.aspx),存储到mysql表中
* 正则表达式解析网页
* 爬取各分类下商品信息（未完成）


### MySQL数据库 sqlgo
* InnoDB引擎，更高粒度的锁，更好的并发性能
* 数据库结构
![](root/static/db2.png)
* 建表语句
    * 分类目录
    
        `create table IF NOT EXISTS classTable( /* 创建分类目录*/
         								class_id int auto_increment, /*分类id，自增*/
         								class_name varchar(20),/*分类名称*/
         								class_href varchar(50),/*分类链接*/
         								primary key(class_id))
         								engine=InnoDB default charset=utf8
         								`
        
    * 商品信息表（获取时如何跟分类进行联动）
    
        `create table IF NOT EXISTS goodsTable(/*创建商品信息表*/
								goods_id int auto_increment,/*商品id*/
								class_id int,/*商品对应分类id*/
								goods_name varchar(100),/*商品名称*/
								goods_price float,		/*商品价格*/
								goods_href text,		/*商品链接*/
								primary key(goods_id))
								engine=InnoDB default charset=utf8
								`
    * 父子关系表（如何递归地建立关系）
    
        `create table IF NOT EXISTS classRelate(/*创建分类关系表*/
        class_id int,/*分类id*/
		pid int,/*父分类id*/
		primary key(class_id))
        engine=InnoDB default charset=utf8`
	
	* 用户信息表

		`create table IF NOT EXISTS userInfo(/*创建用户信息表*/
								user_id int auto_increment,/*用户id*/
								user_name varchar(20),/*用户名*/
								user_email varchar(50),/*用户邮箱*/
								primary key(user_id))
								engine=InnoDB default charset=utf8
								`


### 服务端 web.go
* golang，暂不借助其他框架
* 实现根据分类id或分类名称进行查询

* 查询语句
    * queryById  
    
        `queryById = select a.class_id,a.class_name,a.class_href 
     					from classTable as a join classRelate as b 
     					where a.class_id = b.class_id and 
     					(b.pid = ? or b.class_id = ?);` 
    * queryByName
    
        `queryByName = select a.class_id,a.class_name,a.class_href 
      						from classTable as a join classRelate as b 
      						where a.class_id = b.class_id and 
      						a.class_name regexp '%s';`
   


### 注册登录功能 redis
- [x] redis实现登录验证


### 前端页面
* html+bootstrap+js实现
* 登录和分类管理两个页面
* 登录界面（./root/test.html）
	* 登录
		* dchest/captcha 验证码
		* js局部加载输出提示
	* 注册
		* 正则表达式进行注册邮箱校验
		* js局部加载输出提示
	* 找回密码（未完成）

* 分类管理界面（./root/test2.html）
	* 分类管理
		* 分类id或分类名称查询
		* ajxa加载模板，生成各分类信息（未完成）
		select b.class_name,left(a.goods_name,10),a.goods_price from goodsTable as a join classTable as b where b.class_name regexp '水' and b.class_id = a.class_id;

	* 商品管理（未完成）
	* 用户管理（未完成）

