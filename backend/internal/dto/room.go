package dto

// CreateRoomRequest 创建房间请求。
type CreateRoomRequest struct {
	TTLHours int `json:"ttlHours" binding:"omitempty,gte=1,lte=720"`
}

// JoinRoomRequest 通过房间码加入房间请求。
type JoinRoomRequest struct {
	RoomCode string `json:"roomCode" binding:"required,room_code"`
}

// LeaveRoomRequest 通过房间码离开房间请求。
type LeaveRoomRequest struct {
	RoomCode string `json:"roomCode" binding:"required,room_code"`
}

