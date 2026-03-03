package errno

const (
	// CodeSuccess 代表请求成功。
	CodeSuccess = 0

	// 通用错误码。
	CodeInternalError  = 10000
	CodeInvalidParam   = 10001
	CodeUnauthorized   = 10002

	// 房间相关错误码。
	CodeRoomNotFound = 20001
	CodeRoomClosed   = 20002

	// 对局相关错误码。
	CodeMatchNotFound = 30001
	CodeScoreEntryNotFound = 30002
)

var messages = map[int]string{
	CodeSuccess:      "success",
	CodeInternalError: "internal error",
	CodeInvalidParam: "invalid parameters",
	CodeUnauthorized: "unauthorized",
	CodeRoomNotFound: "room not found",
	CodeRoomClosed:   "room is closed",
	CodeMatchNotFound: "match not found",
	CodeScoreEntryNotFound: "score entry not found",
}

// Message 返回错误码对应的默认消息。
func Message(code int) string {
	msg, ok := messages[code]
	if !ok {
		return messages[CodeInternalError]
	}
	return msg
}
