# Bug 复现说明

## Bug 是什么

并发计数存在数据竞争，普通测试也会出现计数不准。

## 如何触发

运行 `go test -race ./internal/taskrules -run TestCounterConcurrentAdds -count=1`。

## 错误信息

```text
WARNING: DATA RACE
--- FAIL: TestCounterConcurrentAdds
    rules_test.go:9: counter value = ..., want 8000
```
