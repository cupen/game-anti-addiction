# Introduction
Go SDK of china national video game anti-addiction system.

# Status
alpha

# TestSuite
[testcase01~08(Chinese)](https://github.com/cupen/game-anti-addiction/tree/master/_examples/testsuite)


# Usages
* idcard check
  ```go
  c := auth.NewClient(appId, bizId, secretKey)
  req := idcard.CheckRequest{IDNum:"xx", Name:"xx", AI:"xx"}
  resp, err := req.Do(c)
  ```

* idcard query 
  ```go
  req := idcard.QueryRequest{AI:"xx"}
  resp, err := req.Do(c)
  ```

* behavior upload
  ```go
  events := []behavior.LoginOutEvent{{}, {}} 
  req := behavior.LoginOutRequest{Collection: events}
  resp, err := req.Do(c)
  ```

* message queue and producer, consumer
  ```go
  c := auth.NewClient(appId, bizId, secretKey)
  queue, err := redisstream.New(redisUrl, "behavior")

  // producer
  obj := behavior.NewLogin(...) // or NewLogout(...)
  data, _ := json.Marshal(obj)
  err = queue.Write(data)

  // consumer
  c := auth.NewClient(appId, bizId, secretKey)
  consumerFunc := behavior.ConsumerFunc(c, 128, 100)
  consumer := consumer.New(queue, consumerFunc)
  consumer.Start()

  // consumer(manually)
  msgList, err := queue.Read(1024, 1*time.Second)
  reqList, err := behavior.DecodeLoginOutRequest(msgList, 128)
  for _, req := range reqList {
      resp, err := req.Do(c)
  }
  ```

# More Usages
<details>
	<summary> out-of-box way </summary>

 * all-in-one
  ```go
  c := auth.NewClient(appId, bizId, secretKey)
  gaa, err := outofbox.New(c, redisUrl)
  gaa.Start(nil)

  queue := gaa.GetBehaviorQueue()
  queue.Write(...)
  ```
</details>

# License
MIT License


# About video game anti-addiction (Chinese)

* [网络游戏防沉迷系统](https://zh.wikipedia.org/zh-hans/%E7%BD%91%E7%BB%9C%E6%B8%B8%E6%88%8F%E9%98%B2%E6%B2%89%E8%BF%B7%E7%B3%BB%E7%BB%9F)
* [国家层面的实名验证系统已基本建成](http://youxiputao.com/articles/21386)

* [国家《未成年人保护法》第七十五条](http://www.gov.cn/xinwen/2020-10/18/content_5552113.htm)  
