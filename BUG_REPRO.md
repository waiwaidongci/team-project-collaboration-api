# Bug 复现说明

## Bug 是什么

任务负责人为空值或零值时会被当成有效负责人，空值格式化会引发 panic，负责人比较结果也反了。

## 如何触发

运行 `go test ./internal/taskrules -run 'TestAssignee' -count=1`。

## 错误信息

```text
--- FAIL: TestAssigneeRules
    rules_test.go:11: zero assignee should be invalid
--- FAIL: TestAssigneeHelpers
    rules_test.go:29: SameAssignee should return true for equal values
```
