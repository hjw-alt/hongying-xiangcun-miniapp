# 影绘湘潭后端

Go + MySQL 实现的湘潭非遗纸影戏绘画共创平台 API。

## 启动

如果你的 MySQL root 账号没有密码，直接运行：

```powershell
cd backend
go run .
```

如果 root 有密码，用环境变量传入 DSN：

```powershell
$env:MYSQL_DSN="root:你的密码@tcp(127.0.0.1:3306)/shadowart_xiangtan?charset=utf8mb4&parseTime=true&loc=Local"
go run .
```

服务启动后会自动创建 `shadowart_xiangtan` 数据库、`artworks` 表、`comments` 表，并写入 3 个示例作品。

## 比赛定位

以湘潭非遗纸影戏为核心，围绕“纸影人物”“湘潭印象”“戏台故事”三类作品进行共创。用户可以浏览作品、点赞评论、上传自己的绘画再创作；后台可完成作品和评论的增删改查。

## 接口

| 方法 | 地址 | 说明 |
| --- | --- | --- |
| GET | `/api/health` | 健康检查 |
| GET | `/api/stats` | 作品/点赞/主题统计 |
| GET | `/api/artworks` | 列表，支持 `keyword`、`theme`、`sort=likes` |
| POST | `/api/artworks` | 创建作品 |
| GET | `/api/artworks/{id}` | 作品详情和评论 |
| PUT | `/api/artworks/{id}` | 更新作品 |
| DELETE | `/api/artworks/{id}` | 删除作品 |
| POST | `/api/artworks/{id}/like` | 点赞 |
| GET | `/api/artworks/mine` | 读取当前本机用户作品 |
| GET | `/api/clips` | 在线影戏剧目列表 |
| POST | `/api/clips` | 上传影戏剧目信息 |
| POST | `/api/clips/{id}/like` | 剧目点赞 |
| GET | `/api/favorites/artworks` | 收藏影偶 |
| GET | `/api/favorites/clips` | 收藏剧目 |
| POST | `/api/favorites/{type}/{id}/toggle` | 切换收藏状态 |
| GET | `/api/activities` | 研学 / 展演活动 |
| POST | `/api/bookings` | 创建活动预约订单 |
| GET | `/api/bookings` | 预约记录 |
| POST | `/api/bookings/{id}/cancel` | 取消预约并回补名额 |
| POST | `/api/artworks/{id}/comments` | 新增评论 |
| DELETE | `/api/artworks/{id}/comments/{commentID}` | 删除评论 |
| POST | `/api/upload` | 上传图片，表单字段 `file` |
| POST | `/api/upload/video` | 上传短视频，表单字段 `file`，最大 100MB |

作品图片和视频上传到 `backend/static/uploads`。需要身份的接口会读取 `X-User-Key` 请求头；小程序端会自动生成并携带本机用户标识。

## 微信联调

微信开发者工具可以访问 `http://127.0.0.1:8080`。真机预览时，把 `frontend/config.js` 里的 `BASE_URL` 改成电脑局域网 IP，例如 `http://192.168.1.8:8080`，并在微信开发者工具里勾选“不校验合法域名”。
