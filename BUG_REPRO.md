# Bug 复现说明

## Bug 是什么

任务标签整理会把空字符串、重复标签和仅大小写不同的标签原样保留。

## 如何触发

运行 `go test ./internal/taskrules -run TestLabelRules -count=1`，检查 `NormalizeLabels`、`DeduplicateLabels`、`HasLabel` 和 `CanonicalLabels` 的标签规则。

## 错误信息

```text
--- FAIL: TestLabelRules
    rules_test.go:12: NormalizeLabels(["  urgent " "" "urgent" "high"]) = []string{"  urgent ", "", "urgent", "high"}, want []string{"urgent", "high"}
```
