# go-certvault

Go 实现的 TLS 证书库：PEM 入库、CSR 签发、到期扫描、轮换回滚、吊销清单与 trust bundle 导出。配套 `certd` 管理服务与静态页。

## 运行

```bash
go test ./... -count=1
go run ./cmd/certd -addr :8091 -web web
```

浏览器打开 `http://localhost:8091/`：导入证书、查看到期、触发扫描与轮换预览。

## 库面

```go
v := certvault.New()
id, err := v.ImportPEM(certPEM, keyPEM, certvault.Meta{Name: "api"})
_ = v.ScanExpiring(30 * 24 * time.Hour)
_ = v.Rotate(id, newCert, newKey)
_ = v.Close()
```
