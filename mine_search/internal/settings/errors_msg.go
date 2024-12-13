package settings

type ErrMsgBase struct {
	Code int
	Msg  string
}
type ErrMsgs struct {
	DbFailed      *ErrMsgBase
	Params        *ErrMsgBase
	RateLimit     *ErrMsgBase
	TokenRequired *ErrMsgBase
	TokenExpired  *ErrMsgBase
	InvalidToken  *ErrMsgBase
	CallFail      *ErrMsgBase
	ErrorAccess   *ErrMsgBase
	ErrorAddress  *ErrMsgBase
	HideBlock     *ErrMsgBase // Error của app sử dụng cho block Swap wifi 6 - 10405: ẩn block
}

func NewErrMsgs() *ErrMsgs {
	return &ErrMsgs{
		DbFailed: &ErrMsgBase{
			Code: 300,
			Msg:  "Kết nối hệ thống lỗi, vui lòng thử lại sau ít phút",
		},
		Params: &ErrMsgBase{
			Code: 400,
			Msg:  "Sai thông tin đầu vào, vui lòng kiểm tra lại thông tin",
		},
		RateLimit: &ErrMsgBase{
			Code: 400,
			Msg:  "Bạn truy cập quá nhanh.",
		},
		TokenRequired: &ErrMsgBase{
			Code: 1003,
			Msg:  "Token không tồn tại",
		},
		TokenExpired: &ErrMsgBase{
			Code: 1001,
			Msg:  "Token hết hạn",
		},
		InvalidToken: &ErrMsgBase{
			Code: 1002,
			Msg:  "Token không hợp lệ",
		}, CallFail: &ErrMsgBase{
			Code: 300,
			Msg:  "Hệ thống đang bận, vui lòng thử lại sau ít phút",
		}, ErrorAccess: &ErrMsgBase{
			Code: 301,
			Msg:  "Có lỗi trong quá trình xử lý, vui lòng thử lại sau ít phút",
		}, ErrorAddress: &ErrMsgBase{
			Code: 302,
			Msg:  "Hệ thống đang bận, vui lòng thử lại sau ít phút :(",
		}, HideBlock: &ErrMsgBase{
			Code: 10405,
			Msg:  "Có lỗi trong quá trình xử lý, vui lòng thử lại sau ít phút",
		},
	}
}
