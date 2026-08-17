# Bug 复现说明

## Bug 是什么

错误包装丢失原始 sentinel，错误码始终返回 internal，超时错误不会被识别为可重试。

## 如何触发

运行 `go test ./internal/taskrules -run TestErrorClassification -count=1`。

## 错误信息

```text
--- FAIL: TestErrorClassification
    rules_test.go:11: wrapped error should preserve the not-found sentinel
```
