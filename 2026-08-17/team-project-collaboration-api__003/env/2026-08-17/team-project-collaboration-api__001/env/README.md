# Team Project Task API

多人项目任务管理 API，使用 Go 1.26、Gin 和 PostgreSQL。

## 环境变量

默认值可直接用于本地验证：

- `DATABASE_URL`: `postgres://postgres:postgres@127.0.0.1:55432/team_project_task?sslmode=disable`
- `HTTP_ADDR`: `:18108`
- `JWT_SECRET`: `dev-only-change-me`
- `TOKEN_TTL`: `72h`

## 启动

```bash
go run ./cmd/migrate
go run ./cmd/api
```

## 测试

```bash
go test ./...
```

## API 概览

认证：

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`

项目：

- `POST /api/v1/projects`
- `GET /api/v1/projects`
- `GET /api/v1/projects/:projectID`
- `POST /api/v1/projects/:projectID/members`
- `GET /api/v1/projects/:projectID/members`
- `GET /api/v1/projects/:projectID/activities`

任务：

- `POST /api/v1/projects/:projectID/tasks`
- `GET /api/v1/projects/:projectID/tasks`
- `GET /api/v1/projects/:projectID/tasks/:taskID`
- `PATCH /api/v1/projects/:projectID/tasks/:taskID`
- `DELETE /api/v1/projects/:projectID/tasks/:taskID`
- `PATCH /api/v1/projects/:projectID/tasks/:taskID/status`
- `PATCH /api/v1/projects/:projectID/tasks/:taskID/assignee`
- `POST /api/v1/projects/:projectID/tasks/:taskID/comments`
- `GET /api/v1/projects/:projectID/tasks/:taskID/comments`
- `GET /api/v1/tasks/overdue`

所有项目接口需要请求头 `Authorization: Bearer <token>`，并且仅项目成员可访问。状态和负责人变更会写入项目动态，任务分页在 repository 层通过 `LIMIT/OFFSET` 完成。
