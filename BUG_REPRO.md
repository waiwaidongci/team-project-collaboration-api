# Bug 复现说明

## Bug 是什么

清理顺序错误、父上下文取消没有传给子上下文、资源关闭被跳过，panic 也没有被恢复。

## 如何触发

运行 `go test ./internal/taskrules -run 'Test(CleanupOrder|ContextPropagation|CleanupResource|RecoverPanic)' -count=1`。

## 错误信息

```text
--- FAIL: TestCleanupOrder
--- FAIL: TestContextPropagation
--- FAIL: TestCleanupResource
panic: boom
```
