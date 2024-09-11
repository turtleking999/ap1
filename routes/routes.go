package routes

import (
	"airline-booking/controllers"

	"github.com/fasthttp/router"
)

// SetupRoutes 配置所有的路由
func SetupRoutes(r *router.Router, fc *controllers.FlightController) {
	// POST /flights/search: 發起航班搜索
	// 設計要點：
	// 1. 異步處理：立即返回請求ID，提高系統響應性和並發處理能力
	// 2. 資源管理：更好地控制和分配系統資源，尤其是在高負載情況下
	// 3. 可擴展性：便於實現負載均衡和水平擴展
	r.POST("/flights/search", fc.SearchFlights)

	// GET /flights/results: 獲取搜索結果
	// 設計要點：
	// 1. 靈活性：支持長時間運行的複雜查詢，避免連接超時
	// 2. 用戶體驗：允許客戶端實現進度展示和部分結果顯示
	// 3. 錯誤處理：提供更好的重試機制和錯誤恢復能力
	r.GET("/flights/results", fc.GetSearchResults)

}
