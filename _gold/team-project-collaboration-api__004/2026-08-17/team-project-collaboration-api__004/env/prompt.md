# 项目：多人项目任务管理API

从0做一个Go多人项目任务管理API，用Gin和PostgreSQL实现。用户可以创建项目并邀请成员，项目内可以维护任务，任务包含标题、描述、优先级、状态、负责人、截止日期和标签；支持修改任务状态、变更负责人、添加评论、查看项目动态、按状态和负责人筛选任务、查看自己负责的过期任务。代码按Go企业项目结构组织：cmd/api/main.go、cmd/migrate、internal/config、internal/model、internal/repository、internal/service、internal/handler、internal/middleware、internal/router、internal/pkg、migrations。只有项目成员可以访问项目内任务，状态变更和负责人变更要写入活动记录，任务分页查询在repository层完成。
