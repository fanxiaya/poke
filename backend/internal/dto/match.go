package dto

// CreateMatchRequest 创建对局请求。
type CreateMatchRequest struct {
	RoomCode string `json:"roomCode" binding:"required,room_code"`
}

// TransferScoreRequest 记录玩家之间的分数转移。
type TransferScoreRequest struct {
	FromUserID string `json:"fromUserId" binding:"required,min=1,max=64"`
	ToUserID   string `json:"toUserId" binding:"required,min=1,max=64"`
	Score      int    `json:"score" binding:"required,gt=0"`
}

// GrantScoreRequest 记录单个用户凭空加分。
type GrantScoreRequest struct {
	UserID string `json:"userId" binding:"required,min=1,max=64"`
	Score  int    `json:"score" binding:"required,gt=0"`
}

// UpdateScoreEntryRequest 修改一条记分流水的分数。
type UpdateScoreEntryRequest struct {
	Score int `json:"score" binding:"required,gt=0"`
}
