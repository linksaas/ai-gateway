# 简介

这个是[凌鲨](https://www.linksaas.pro)的一个副产品。用于给研发团队提供安全可控的AI能力的网关。

## 能力
- 提供内容检查脚本
- 提供http和脚本两种AI转发能力
- 定义了代码相关的AI接口

# 使用

## 部署
//TODO

## 配置
config.yaml是唯一的配置文件

```yaml
port: 8080
ssl:
  enable: false
  cert: your_ssl_cert
  key: your_ssl_key
secret: your_serv_secret
dev: true
logdir: logs
tokenttl: 600
checkscript: script/check.go
provider: 
  coding:
    - backend: script://script/coding_provider.go #suport http:// and script://
      checkscript: script/check.go
      # support lang:python,c,cplusplus,java,csharp,visualbasic,javascript,sql,asm,php,r,go,matlab,swift,delphi,ruby,perl,objc,rust
      complete: []
      convert: []
      explain: []
      fixerror: []
      gentest: []

```
provider 可以提供多个，处理的时候至上而下，找到第一个配置的provider进入后续处理过程。

> 只有在dev为true的时候，可以通过/api/dev/genToken生成访问令牌。

## 脚本
我们的使用[yaegi](https://github.com/traefik/yaegi)作为执行引擎。开发的时候可以当golang进行开发测试。

> 目前只能import标准库，并且需要有golang 1.19以上环境，并设置环境变量GOROOT

### 检查脚本

检查脚本存在两个层级，全局检查脚本和provider层面的检查脚本。在provider层面的检查脚本不存在的情况下，会使用全局检查脚本。

```go
package script

//check content of request
//@return false for forbid process,true for continue process
func CheckContent(apiUrl, content string) bool {
	//TODO
	return true
}
```

### 代码AI功能实现脚本

```go
package script

//process /api/coding/complete/:lang
func Complete(lang, content string) []string {
	//TODO
	return nil
}

//process /api/coding/convert/:lang
func Convert(lang, destLang, content string) []string {
	//TODO
	return nil
}

//process /api/coding/explain/:lang
func Explain(lang, content string) []string {
	//TODO
	return nil
}

//process /api/coding/fixError/:lang
func Fixerror(lang, errStr string) []string {
	//TODO
	return nil
}

//process /api/coding/genTest/:lang
func Gentest(lang, content string) []string {
	//TODO
	return nil
}
```