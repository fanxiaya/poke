package model

import "errors"

var (
	// Room 相关错误。
	ErrRoomClosed       = errors.New("房间已关闭")
	ErrUserIDEmpty      = errors.New("用户ID不能为空")
	ErrUserAlreadyInRoom = errors.New("用户已在房间中")
	ErrUserNotInRoom    = errors.New("用户不在房间中")

	// Match 相关错误。
	ErrMatchFinished   = errors.New("对局已结束")
	ErrUserNotInMatch  = errors.New("用户不在对局中")
	ErrScoreMustBePositive = errors.New("分数必须为正数")
	ErrEntryNotFound = errors.New("记分流水不存在")
)
