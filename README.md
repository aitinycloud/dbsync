# dbsync
database sync tools.

## config
config.json file.

1. SrcDB source database config.
2. DesDB destcation database config.
3. DataSync data sync handle config.


## build
make

## run
./bin/start.sh

## 说明
dbsync是数据库同步工具,使用go编写，主要特性有
1. 支持数据库种类多(oracal,postgresql,mysql等,其他添加中)
2. 灵活支持自定义SQL,支持多表和视图
3. 支持字段映射,可以配置字段的映射关系
4. 性能优化,批量操作

## bench
a. 本地虚拟机 centos6 2核2G 
b. postgres (PostgreSQL) 10.6
c. mysql  Ver 14.14 Distrib 5.1.73

1. insert 117275 条记录 14S TPS 8376/s
2. update 117275 条记录 44S TPS 2665/s (update操作中计算md5会花一点时间)

## 问题
如果有任何疑问或错误，欢迎在 issues 进行提问或给予修正意见。 如果喜欢或对你有所帮助，欢迎 Star。

我是小铭哥 专注做后台开发和系统架构
我在重庆两江新区。

<p align="center"><img src="doc/wx/2074030723.jpg" alt="wx Logo"></p>



## License

所有文章采用[知识共享署名-非商业性使用-相同方式共享 3.0 中国大陆许可协议](https://creativecommons.org/licenses/by-nc-sa/3.0/cn/)进行许可。
