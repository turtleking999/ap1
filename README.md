# 航空訂票系統

這是一個使用 Go 語言開發的航空訂票系統後端 API。

## 項目架構
airline-booking/
├── config/
│ └── config.go # 配置文件，包含數據庫和 Redis 初始化
├── controllers/
│ └── flight_controller.go # 處理 HTTP 請求的控制器
├── models/
│ └── flight.go # 數據模型定義
├── repositories/
│ └── flight_repository.go # 數據庫操作邏輯
├── services/
│ └── flight_service.go # 業務邏輯層
├── main.go # 應用程序入口
├── go.mod # Go 模塊定義
├── go.sum # Go 模塊依賴版本
├── Dockerfile # Docker 構建文件
├── docker-compose.yml # Docker Compose 配置文件
└── README.md # 項目說明文檔

## 技術棧

- Go 1.22
- PostgreSQL
- Redis
- Docker
- fasthttp

## 主要功能

- 航班搜索（基於起點、目的地和日期）
- 分頁返回航班列表
- Redis 緩存以提高性能

## 打包和啟動方式

### 使用 Docker Compose（推薦）

1. 確保你已經安裝了 Docker 和 Docker Compose。

2. 在項目根目錄下運行以下命令：

   ```bash
   docker-compose up -d
   ```

   這將啟動應用程序、PostgreSQL 和 Redis 服務。

3. 要停止服務，運行：

   ```bash
   docker-compose down
   ```

### 使用 Docker

1. 構建 Docker 鏡像:

   ```bash
   docker build -t airline-booking .
   ```

2. 運行 Docker 容器:

   ```bash
   docker run -p 8080:8080 airline-booking
   ```

### 本地運行

1. 確保已安裝 Go 1.22 或更高版本。

2. 安裝依賴:

   ```bash
   go mod download
   ```

3. 編譯並運行應用:

   ```bash
   go build -o airline-booking
   ./airline-booking
   ```

## 配置

在運行應用之前，請確保正確設置了以下環境變量或在 `config/config.go` 中修改相應的值:

- `DBHost`: PostgreSQL 主機地址
- `DBPort`: PostgreSQL 端口
- `DBUser`: PostgreSQL 用戶名
- `DBPassword`: PostgreSQL 密碼
- `DBName`: PostgreSQL 數據庫名稱
- `ServerPort`: 應用服務器端口
- `RedisAddr`: Redis 服務器地址
- `RedisPassword`: Redis 密碼（如果有）
- `RedisDB`: Redis 數據庫編號

使用 Docker Compose 時，這些配置已經在 `docker-compose.yml` 文件中設置好了。

## API 端點

- `POST /search`: 搜索航班
  - 請求體示例:
    ```json
    {
      "origin": "New York",
      "destination": "London",
      "date": "2023-05-01T00:00:00Z",
      "page": 1,
      "page_size": 10
    }
    ```

## 注意事項

- 使用 Docker Compose 時，確保沒有其他服務佔用了 8080（應用）、5432（PostgreSQL）和 6379（Redis）端口。
- 在生產環境中，請適當調整 `docker-compose.yml` 中的配置以確保安全性和性能。
- 本項目使用 fasthttp 作為 HTTP 服務器，提供高性能的請求處理。

## 貢獻

歡迎提交 issues 和 pull requests 來改進這個項目。

## 許可證

[MIT License](LICENSE)
