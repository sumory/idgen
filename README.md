## id generator

[![Build Status](https://travis-ci.org/sumory/idgen.svg?branch=master)](https://travis-ci.org/sumory/idgen) [![](http://gocover.io/_badge/github.com/sumory/idgen)](http://gocover.io/github.com/sumory/idgen) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/sumory/idgen/master/LICENSE)

#### 介绍

 - 基于`snowflake`算法实现的id生成器
 - 这是go版本，java版本可查看[IdWorker.java](https://github.com/sumory/uc/blob/master/src/com/sumory/uc/id/IdWorker.java)




#### 使用

go get github.com/sumory/idgen

使用前请先了解`snowflake`算法，并知晓其注意事项.


##### 基本使用

每个由idgen生成的id都是int64的正整数，且每个id都可以解析得到它的生成者的标识`workerId`.

```go
workerId := 1
err, idWorker := idgen.NewIdWorker(workerId)
err, nextId := idWorker.NextId()
```

##### 获取short Id

idgen使用[baseN4go](https://github.com/sumory/baseN4go)缩短id，具体参见baseN4go使用方法.

```go
workerId := 1
err, idWorker := idgen.NewIdWorker(workerId)
err, nextId := idWorker.ShortId()
```

##### 获取生成器标识workerId

```go
workerId := 1
err, idWorker := idgen.NewIdWorker(workerId)
err, nextId := idWorker.NextId()
wId := idWorker.WorkerId(newId)//wId == workerId
```

##### 其它
参见测试文件[idgen_test.go](./idgen_test.go)和源文件




#### 测试

需要[goconvey](https://github.com/smartystreets/goconvey)支持

```shell
go get github.com/smartystreets/goconvey
go test -v -cover // or $GOPATH/bin/goconvey
```

