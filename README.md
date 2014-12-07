## id generator

#### 介绍

 - 基于`snowflake`算法实现的id生成器
 - 这是go版本，java版本可查看[IdWorker.java][1]


#### 使用

go get github.com/sumory/idgen

使用前请先了解`snowflake`算法，并知晓其注意事项.


#### 基本使用

每个由idgen生成的id都是int64的正整数，且每个id都可以解析得到它的生成者的标识`workerId`.

```
workerId := 1
err, idWorker := idgen.NewIdWorker(workerId)
err, nextId := idWorker.NextId()
```

#### 获取short Id

idgen使用[baseN4go](https://github.com/sumory/baseN4go)缩短id，具体参见baseN4go使用方法.

```
workerId := 1
err, idWorker := idgen.NewIdWorker(workerId)
err, nextId := idWorker.ShortId()
```

#### 获取生成器标识workerId

```
workerId := 1
err, idWorker := idgen.NewIdWorker(workerId)
err, nextId := idWorker.NextId()
wId := idWorker.WorkerId(newId)//wId == workerId
```

#### 其它
参见测试文件[idgen_test.go](./idgen_test.go)和源文件