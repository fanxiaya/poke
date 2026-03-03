package model

import (
	"crypto/rand"
	"math/big"
	"time"
)

type RoomStatus string

const (
	RoomStatusOpen   RoomStatus = "open"
	RoomStatusClosed RoomStatus = "closed"
)

const (
	roomCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	roomCodeLength  = 5
)

// Room 表示局外房间，用于承载用户集合与当前对局。
type Room struct {
	ID             string     `bson:"_id,omitempty" json:"id"`
	RoomCode       string     `bson:"roomCode" json:"roomCode"`
	HostUserID     string     `bson:"hostUserId" json:"hostUserId"`
	MemberUserIDs  []string   `bson:"memberUserIds" json:"memberUserIds"`
	CurrentMatchID string     `bson:"currentMatchId,omitempty" json:"currentMatchId,omitempty"`
	Status         RoomStatus `bson:"status" json:"status"`
	CreatedAt      time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `bson:"updatedAt" json:"updatedAt"`
	ExpireAt       time.Time  `bson:"expireAt" json:"expireAt"`
}

// NewRoom 创建一个初始化完成的房间实例。
func NewRoom(hostUserID string, ttl time.Duration) *Room {
	r := &Room{
		HostUserID: hostUserID,
	}
	r.InitForCreate(ttl)
	return r
}

// InitForCreate 用于新建房间时设置默认值。
func (r *Room) InitForCreate(ttl time.Duration) {
	now := time.Now()
	if ttl <= 0 {
		ttl = 24 * time.Hour * 30
	}
	if r.RoomCode == "" {
		r.RoomCode = GenerateRoomCode()
	}
	if r.Status == "" {
		r.Status = RoomStatusOpen
	}
	if r.MemberUserIDs == nil {
		r.MemberUserIDs = make([]string, 0, 8)
	}
	if r.HostUserID != "" && !containsString(r.MemberUserIDs, r.HostUserID) {
		r.MemberUserIDs = append(r.MemberUserIDs, r.HostUserID)
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}
	r.UpdatedAt = now
	r.ExpireAt = now.Add(ttl)
}

// Touch 在更新房间时刷新更新时间。
func (r *Room) Touch() {
	r.UpdatedAt = time.Now()
}

// Close 将房间状态置为关闭。
func (r *Room) Close() {
	r.Status = RoomStatusClosed
	r.Touch()
}

// CanJoin 判断房间当前是否允许加入。
func (r *Room) CanJoin() bool {
	return r.Status == RoomStatusOpen
}

// Join 将用户加入房间成员列表。
func (r *Room) Join(userID string) error {
	if userID == "" {
		return ErrUserIDEmpty
	}
	if !r.CanJoin() {
		return ErrRoomClosed
	}
	if containsString(r.MemberUserIDs, userID) {
		return ErrUserAlreadyInRoom
	}

	r.MemberUserIDs = append(r.MemberUserIDs, userID)
	r.Touch()
	return nil
}

// Leave 将用户从房间成员列表中移除。
func (r *Room) Leave(userID string) error {
	if userID == "" {
		return ErrUserIDEmpty
	}

	index := -1
	for i, id := range r.MemberUserIDs {
		if id == userID {
			index = i
			break
		}
	}
	if index < 0 {
		return ErrUserNotInRoom
	}

	r.MemberUserIDs = append(r.MemberUserIDs[:index], r.MemberUserIDs[index+1:]...)
	r.Touch()
	return nil
}

// GenerateRoomCode 生成 5 位字母数字房间码。
func GenerateRoomCode() string {
	b := make([]byte, roomCodeLength)
	max := big.NewInt(int64(len(roomCodeCharset)))
	for i := 0; i < roomCodeLength; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			b[i] = roomCodeCharset[0]
			continue
		}
		b[i] = roomCodeCharset[n.Int64()]
	}
	return string(b)
}

func containsString(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
