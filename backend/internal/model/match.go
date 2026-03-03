package model

import "time"

type MatchStatus string

const (
	MatchStatusPlaying  MatchStatus = "playing"
	MatchStatusFinished MatchStatus = "finished"
)

type MatchEntryType string

const (
	MatchEntryTypeTransfer MatchEntryType = "transfer"
	MatchEntryTypeGrant    MatchEntryType = "grant"
)

// MatchPlayer 保存对局内玩家快照信息。
type MatchPlayer struct {
	UserID        string `bson:"userId" json:"userId"`
	NicknameSnap  string `bson:"nicknameSnapshot" json:"nicknameSnapshot"`
	AvatarURLSnap string `bson:"avatarSnapshot" json:"avatarSnapshot"`
}

// MatchScoreEntry 保存一条记分流水。
type MatchScoreEntry struct {
	EntryNo        int            `bson:"entryNo" json:"entryNo"`
	Type           MatchEntryType `bson:"type" json:"type"`
	FromUserID     string         `bson:"fromUserId,omitempty" json:"fromUserId,omitempty"`
	ToUserID       string         `bson:"toUserId,omitempty" json:"toUserId,omitempty"`
	UserID         string         `bson:"userId,omitempty" json:"userId,omitempty"`
	Score          int            `bson:"score" json:"score"`
	OperatorUserID string         `bson:"operatorUserId" json:"operatorUserId"`
	CreatedAt      time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time      `bson:"updatedAt" json:"updatedAt"`
}

// Match 表示一场可持续记分的对局。
type Match struct {
	ID          string         `bson:"_id,omitempty" json:"id"`
	RoomCode    string         `bson:"roomCode" json:"roomCode"`
	CreatedBy   string         `bson:"createdBy" json:"createdBy"`
	Status      MatchStatus    `bson:"status" json:"status"`
	Players     []MatchPlayer  `bson:"players" json:"players"`
	ScoreEntries []MatchScoreEntry `bson:"scoreEntries" json:"scoreEntries"`
	TotalScores map[string]int `bson:"totalScores" json:"totalScores"`
	StartedAt   *time.Time     `bson:"startedAt,omitempty" json:"startedAt,omitempty"`
	EndedAt     *time.Time     `bson:"endedAt,omitempty" json:"endedAt,omitempty"`
	CreatedAt   time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time      `bson:"updatedAt" json:"updatedAt"`
}

// NewMatch 创建一个初始化完成的对局实例。
func NewMatch(roomCode, createdBy string, players []MatchPlayer) *Match {
	m := &Match{
		RoomCode:  roomCode,
		CreatedBy: createdBy,
		Players:   players,
	}
	m.InitForCreate()
	return m
}

// InitForCreate 用于新建对局时设置默认值。
func (m *Match) InitForCreate() {
	now := time.Now()
	m.Start()

	if m.ScoreEntries == nil {
		m.ScoreEntries = make([]MatchScoreEntry, 0)
	}
	if m.TotalScores == nil {
		m.TotalScores = make(map[string]int)
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now
}

// Start 将对局状态切换为进行中。
func (m *Match) Start() {
	now := time.Now()
	m.Status = MatchStatusPlaying
	if m.StartedAt == nil {
		m.StartedAt = &now
	}
	m.UpdatedAt = now
}

// Finish 将对局状态切换为已结束。
func (m *Match) Finish() {
	now := time.Now()
	m.Status = MatchStatusFinished
	m.EndedAt = &now
	m.UpdatedAt = now
}

// TransferScore 记录一笔玩家之间的分数转移：from 扣分，to 加分。
func (m *Match) TransferScore(fromUserID, toUserID string, score int, operatorUserID string) (*MatchScoreEntry, error) {
	if m.Status == MatchStatusFinished {
		return nil, ErrMatchFinished
	}
	if score <= 0 {
		return nil, ErrScoreMustBePositive
	}
	if fromUserID == toUserID {
		return nil, ErrTransferSameUser
	}
	if !m.isUserInMatch(fromUserID) || !m.isUserInMatch(toUserID) {
		return nil, ErrUserNotInMatch
	}

	now := time.Now()
	entry := MatchScoreEntry{
		EntryNo:        len(m.ScoreEntries) + 1,
		Type:           MatchEntryTypeTransfer,
		FromUserID:     fromUserID,
		ToUserID:       toUserID,
		Score:          score,
		OperatorUserID: operatorUserID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	m.ScoreEntries = append(m.ScoreEntries, entry)
	m.RecalculateTotalScores()
	m.UpdatedAt = now
	return &entry, nil
}

// GrantScore 记录一笔用户凭空加分流水。
func (m *Match) GrantScore(userID string, score int, operatorUserID string) (*MatchScoreEntry, error) {
	if m.Status == MatchStatusFinished {
		return nil, ErrMatchFinished
	}
	if score <= 0 {
		return nil, ErrScoreMustBePositive
	}
	if !m.isUserInMatch(userID) {
		return nil, ErrUserNotInMatch
	}

	now := time.Now()
	entry := MatchScoreEntry{
		EntryNo:        len(m.ScoreEntries) + 1,
		Type:           MatchEntryTypeGrant,
		UserID:         userID,
		Score:          score,
		OperatorUserID: operatorUserID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	m.ScoreEntries = append(m.ScoreEntries, entry)
	m.RecalculateTotalScores()
	m.UpdatedAt = now
	return &entry, nil
}

// UpdateEntryScore 修改某条记分流水的分数。
func (m *Match) UpdateEntryScore(entryNo int, score int, operatorUserID string) error {
	if entryNo <= 0 {
		return ErrEntryNotFound
	}
	if score <= 0 {
		return ErrScoreMustBePositive
	}

	index := -1
	for i := range m.ScoreEntries {
		if m.ScoreEntries[i].EntryNo == entryNo {
			index = i
			break
		}
	}
	if index < 0 {
		return ErrEntryNotFound
	}

	target := m.ScoreEntries[index]
	if target.Type == MatchEntryTypeTransfer && target.FromUserID == target.ToUserID {
		return ErrTransferSameUser
	}
	if target.Type == MatchEntryTypeTransfer &&
		(!m.isUserInMatch(target.FromUserID) || !m.isUserInMatch(target.ToUserID)) {
		return ErrUserNotInMatch
	}
	if target.Type == MatchEntryTypeGrant && !m.isUserInMatch(target.UserID) {
		return ErrUserNotInMatch
	}

	now := time.Now()
	m.ScoreEntries[index].Score = score
	m.ScoreEntries[index].OperatorUserID = operatorUserID
	m.ScoreEntries[index].UpdatedAt = now
	m.RecalculateTotalScores()
	m.UpdatedAt = now
	return nil
}

// RecalculateTotalScores 根据流水重算累计分。
func (m *Match) RecalculateTotalScores() {
	totals := make(map[string]int, len(m.Players))
	for _, player := range m.Players {
		totals[player.UserID] = 0
	}

	for _, entry := range m.ScoreEntries {
		switch entry.Type {
		case MatchEntryTypeTransfer:
			totals[entry.FromUserID] -= entry.Score
			totals[entry.ToUserID] += entry.Score
		case MatchEntryTypeGrant:
			totals[entry.UserID] += entry.Score
		}
	}
	m.TotalScores = totals
}

func (m *Match) isUserInMatch(userID string) bool {
	if userID == "" {
		return false
	}
	validUsers := make(map[string]struct{}, len(m.Players))
	for _, player := range m.Players {
		validUsers[player.UserID] = struct{}{}
	}
	_, ok := validUsers[userID]
	return ok
}
