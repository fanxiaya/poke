package model

import "time"

const MaxRecentMatches = 10

// User 表示局外用户信息。
type User struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	OpenID         string    `bson:"openid" json:"openid"`
	Nickname       string    `bson:"nickname" json:"nickname"`
	AvatarURL      string    `bson:"avatarUrl" json:"avatarUrl"`
	RecentMatchIDs []string  `bson:"recentMatchIds" json:"recentMatchIds"`
	CreatedAt      time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time `bson:"updatedAt" json:"updatedAt"`
}

// NewUser 创建一个初始化完成的用户实例。
func NewUser(openID, nickname, avatarURL string) *User {
	u := &User{
		OpenID:    openID,
		Nickname:  nickname,
		AvatarURL: avatarURL,
	}
	u.InitForCreate()
	return u
}

// InitForCreate 用于新建用户时初始化默认字段。
func (u *User) InitForCreate() {
	now := time.Now()
	if u.RecentMatchIDs == nil {
		u.RecentMatchIDs = make([]string, 0, MaxRecentMatches)
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	u.NormalizeRecentMatches()
}

// Touch 在更新用户时刷新更新时间。
func (u *User) Touch() {
	u.UpdatedAt = time.Now()
}

// NormalizeRecentMatches 保证最多保留最近 10 局。
func (u *User) NormalizeRecentMatches() {
	if len(u.RecentMatchIDs) <= MaxRecentMatches {
		return
	}
	u.RecentMatchIDs = u.RecentMatchIDs[:MaxRecentMatches]
}

// PrependRecentMatch 在历史最前方插入一局，并自动去重和截断。
func (u *User) PrependRecentMatch(matchID string) {
	if matchID == "" {
		return
	}

	filtered := make([]string, 0, len(u.RecentMatchIDs)+1)
	filtered = append(filtered, matchID)
	for _, id := range u.RecentMatchIDs {
		if id == "" || id == matchID {
			continue
		}
		filtered = append(filtered, id)
	}

	u.RecentMatchIDs = filtered
	u.NormalizeRecentMatches()
	u.Touch()
}
