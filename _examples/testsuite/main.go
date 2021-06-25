package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cupen/game-anti-addiction/auth"
	"github.com/cupen/game-anti-addiction/behavior"
	"github.com/cupen/game-anti-addiction/idcard"
)

var (
	appId     = flag.String("appId", "", "应用标识（APPID）")
	bizId     = flag.String("bizId", "", "游戏备案识别码（bizId）")
	secretKey = flag.String("secretKey", "", "用户密钥（Secret Key）")

	testCase  = flag.String("testCase", "", "测试用例: testcase01 ~ testcase08")
	testCode  = flag.String("testCode", "", "测试码")
	dataDir   = flag.String("dataDir", "./pre-data/", "预制测试数据目录")
	isDebug   = flag.Bool("debug", false, "打开调试日志")
	debugArgs = flag.String("debugArgs", "nonce=&timestamp=", "使用指定参数进行签名和加密. 方便对比官网调试工具生成的例子")
	proxy     = flag.String("proxy", "", "代理地址, 用来查看完整的 https 请求链上下文")
	cacert    = flag.String("cacert", "", "CA证书")
)

func parseCLI() {
	flag.Parse()
	if *appId == "" || *bizId == "" || *secretKey == "" {
		flag.PrintDefaults()
		return
	}
	if *testCode == "" {
		log.Fatalf("missing '--testCode'")
		return
	}
	if *testCase == "" {
		log.Fatalf("missing '--testCase'")
		return
	}
}

func main() {
	parseCLI()
	_debugArgs, err := url.ParseQuery(*debugArgs)
	if err != nil {
		fmt.Printf("invalid debugArgs: %s", *debugArgs)
		return
	}
	c := auth.NewClient(*appId, *bizId, *secretKey)
	if *proxy != "" {
		setupProxy(c, *proxy, *cacert)
	}
	if *isDebug {
		c.ApplyOptions(
			auth.WithDebug(true),
			auth.WithDebugArgs(_debugArgs),
		)
	}
	switch *testCase {
	case "testcase01":
		runTestCase01(c, *dataDir)
	case "testcase02":
		runTestCase02(c, *dataDir)
	case "testcase03":
		runTestCase03(c, *dataDir)
	case "testcase04":
		runTestCase04(c, *dataDir)
	case "testcase05":
		runTestCase05(c, *dataDir)
	case "testcase06":
		runTestCase06(c, *dataDir)
	case "testcase07":
		runTestCase07(c, *dataDir)
	case "testcase08":
		runTestCase08(c, *dataDir)
	}
}

func runTestCase01(c *auth.Client, dataDir string) {
	fpath := dataDir + "idcard-check.success.json"
	cases := []idcard.CheckRequest{}
	mustReadJson(fpath, &cases)
	for i, _case := range cases {
		resp, err := _case.DoTestSuite(c, *testCode)
		if err != nil {
			panic(err)
		}
		rs := makeResult(resp.IsPassed())
		log.Printf("testcase01-实名认证接口(%d): %s", i+1, rs)
	}
}

func runTestCase02(c *auth.Client, dataDir string) {
	fpath := dataDir + "idcard-check.inprogress.json"
	cases := []idcard.CheckRequest{}
	mustReadJson(fpath, &cases)
	for i, _case := range cases {
		resp, err := _case.DoTestSuite(c, *testCode)
		if err != nil {
			panic(err)
		}
		rs := makeResult(resp.IsInProgress())
		log.Printf("testcase02-实名认证接口(%d): %s", i+1, rs)
	}
}

func runTestCase03(c *auth.Client, dataDir string) {
	fpath := dataDir + "idcard-check.fail.json"
	cases := []idcard.CheckRequest{}
	mustReadJson(fpath, &cases)
	for i, _case := range cases {
		resp, err := _case.DoTestSuite(c, *testCode)
		if err != nil {
			panic(err)
		}
		rs := makeResult(!resp.IsPassed() && !resp.IsInProgress())
		log.Printf("testcase03-实名认证接口(%d): %s", i+1, rs)
	}
}

func runTestCase04(c *auth.Client, dataDir string) {
	fpath := dataDir + "idcard-query.success.json"
	cases := []idcard.QueryRequest{}
	mustReadJson(fpath, &cases)
	for i, _case := range cases {
		if _case.AI == "" {
			panic(fmt.Errorf("invalid data: %s", fpath))
		}
		resp, err := _case.DoTestSuite(c, *testCode)
		if err != nil {
			panic(err)
		}
		rs := makeResult(resp.IsPassed())
		log.Printf("testcase04-实名认证结果查询(%d): %s", i+1, rs)
	}
}

func runTestCase05(c *auth.Client, dataDir string) {
	fpath := dataDir + "idcard-query.inprogress.json"
	cases := []idcard.QueryRequest{}
	mustReadJson(fpath, &cases)
	for i, _case := range cases {
		if _case.AI == "" {
			panic(fmt.Errorf("invalid data: %s", fpath))
		}
		resp, err := _case.DoTestSuite(c, *testCode)
		if err != nil {
			panic(err)
		}
		rs := makeResult(resp.IsInProgress())
		log.Printf("testcase05-实名认证结果查询(%d): %s", i+1, rs)
	}
}

func runTestCase06(c *auth.Client, dataDir string) {
	fpath := dataDir + "idcard-query.fail.json"
	cases := []idcard.QueryRequest{}
	mustReadJson(fpath, &cases)
	for i, _case := range cases {
		if _case.AI == "" {
			panic(fmt.Errorf("invalid data: %s", fpath))
		}
		resp, err := _case.DoTestSuite(c, *testCode)
		if err != nil {
			panic(err)
		}
		rs := makeResult(!resp.IsPassed() && !resp.IsInProgress())
		log.Printf("testcase05-实名认证结果查询(%d): %s", i+1, rs)
	}
}

func runTestCase07(c *auth.Client, dataDir string) {
	fpath := dataDir + "behavior-loginout.json"
	cases := []behavior.LoginOutEvent{}
	mustReadJson(fpath, &cases)
	if len(cases) <= 0 {
		panic(fmt.Errorf("empty testdata fpath=%s", fpath))
	}
	for i := range cases {
		_case := &cases[i]
		_case.Num = i + 1
		_case.SessionID = fmt.Sprintf("%d", i+1)
		_case.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
		_case.BehaviorType = behavior.BehaviorTypes.Online
		_case.UserType = behavior.UserTypes.Guest
		_case.DeviceID = newDeviceID(32)
		_case.PlayerID = ""

	}
	req := behavior.LoginOutRequest{Collections: cases}
	resp, err := req.DoTestSuite(c, *testCode)
	if err != nil {
		panic(err)
	}
	if resp.ErrCode != 0 {
		panic(fmt.Errorf("invalid response:%+v", resp))
	}
	rs := makeResult(resp.ErrCode == 0 && len(resp.Data.Results) <= 0)
	log.Printf("testcase07-游戏用户行为数据上报接口:游客模式: %s", rs)
}

func runTestCase08(c *auth.Client, dataDir string) {
	fpath := dataDir + "behavior-loginout.json"
	cases := []behavior.LoginOutEvent{}
	mustReadJson(fpath, &cases)
	if len(cases) <= 0 {
		panic(fmt.Errorf("empty testdata fpath=%s", fpath))
	}
	for i := range cases {
		_case := &cases[i]
		_case.Num = i + 1
		_case.SessionID = fmt.Sprintf("%d", i+1)
		_case.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
		_case.BehaviorType = behavior.BehaviorTypes.Online
		_case.UserType = behavior.UserTypes.UserAuthed
		_case.DeviceID = newDeviceID(32)
		if _case.PlayerID == "" {
			panic(fmt.Errorf("invalid data: empty pi. fpath: %s", fpath))
		}
	}
	req := behavior.LoginOutRequest{Collections: cases}
	resp, err := req.DoTestSuite(c, *testCode)
	if err != nil {
		panic(err)
	}
	if resp.ErrCode != 0 {
		panic(fmt.Errorf("invalid response:%+v", resp))
	}
	rs := makeResult(resp.ErrCode == 0 && len(resp.Data.Results) <= 0)
	log.Printf("testcase08-游戏用户行为数据上报接口:已认证用户: %s", rs)
}

func newDeviceID(length int) string {
	if length <= 0 && length%2 != 0 {
		panic(fmt.Errorf("invalid length(%d).", length))
	}
	randId := make([]byte, length/2)
	if _, err := io.ReadFull(rand.Reader, randId); err != nil {
		panic(err)
	}
	return hex.EncodeToString(randId)
}

func makeResult(b bool) string {
	if !b {
		return "failed"
	}
	return "passed"
}

func mustReadJson(fpath string, obj interface{}) {
	fd, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(fd)
	if err := dec.Decode(obj); err != nil {
		panic(err)
	}
}

func setupProxy(c *auth.Client, proxy string, fpathOfCertCA string) {
	os.Setenv("https_proxy", proxy)
	os.Setenv("http_proxy", proxy)

	cert, err := ioutil.ReadFile(fpathOfCertCA)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}
	c.ApplyOptions(auth.WithHttpClient(&http.Client{
		Timeout:   10 * time.Second,
		Transport: &transport,
	}))
}
