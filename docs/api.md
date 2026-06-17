# Thrum API 文档

## 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/register | 注册 |
| POST | /api/login | 登录 |
| GET | /api/users?q= | 搜索用户 |
| GET | /api/users/:id | 用户主页 |
| GET | /api/posts?page=&page_size=&category_id=&q= | 帖子列表 |
| GET | /api/posts/:id | 帖子详情 |
| GET | /api/posts/:id/comments | 评论列表 |
| GET | /api/announcement | 站内公告 |
| GET | /api/categories | 分类列表 |

## 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/posts | 发帖 |
| PUT | /api/posts/:id | 编辑帖子 |
| DELETE | /api/posts/:id | 删除帖子 |
| POST | /api/posts/:id/comments | 发表评论 |
| PUT | /api/comments/:id | 编辑评论 |
| DELETE | /api/comments/:id | 删除评论 |
| POST | /api/threads | 创建对话主题 |
| GET | /api/threads?with= | 获取主题列表 |
| DELETE | /api/threads/:id | 删除主题 |
| POST | /api/messages | 发送私信 |
| GET | /api/messages | 会话列表 |
| GET | /api/messages/:user_id?thread= | 对话详情 |
| POST | /api/follow | 关注 |
| DELETE | /api/follow/:user_id | 取消关注 |
| GET | /api/following | 我关注的 |
| GET | /api/followers | 关注我的 |
| POST | /api/upload/image | 上传图片 |
| POST | /api/upload/file | 上传文件 |

## 管理员接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/admin/users?q= | 用户列表 |
| PUT | /api/admin/users/:id/ban | 封禁 |
| PUT | /api/admin/users/:id/unban | 解封 |
| POST | /api/admin/announcement | 发布公告 |
| DELETE | /api/admin/announcement | 删除公告 |

## 超管接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/admin/users | 创建用户 |
| DELETE | /api/admin/users/:id | 删除用户 |
| PUT | /api/admin/users/:id/roles | 角色管理 |
| POST | /api/admin/categories | 创建分类 |
| DELETE | /api/admin/categories/:id | 删除分类 |
