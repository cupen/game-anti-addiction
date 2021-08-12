# 防沉迷系统测试套件
包含官方要求的 8 个测试用例，覆盖三个接口: 实名认证、实名认证结果、玩家在线时长上报

# 使用
```bash
git clone https://github.com/cupen/game-anti-addiction --depth 1
cd _examples/testsuite/
go run main.go -appId xxx --bizId xxx --secretKey xxx --testCase testcase01 --testCode ABCDEF
...
# 如果失败，就加个 --debug 参数然后重试并查看日志。
```
# 使用（https 代理）
```bash
go run main.go \
    --appId xxx \
    --bizId xxx \
    --secretKey xxx \
    --testCase testcase02 \
    --testCode ABCDEF \
    --proxy 127.0.0.1:8080 \
    --cacert /a/path/to/ca/certificate.pem
```


# FAQ
* 连接被拒  
检查 IP 白名单配置是否正确, 或等几秒再重试（根据实际体验盲猜的，原因未知）。
