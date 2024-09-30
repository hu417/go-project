package errs

import "errors"

// ResCode .
type ResCode int

const (
	CodeSuccess       ResCode = 200
	CodeInvalidParams ResCode = 400
	CodeError         ResCode = 500

	CodeConvDataErr       ResCode = 50000
	CodeValidateParamsErr ResCode = 50001
	CodeInvalidToken      ResCode = 50002
	CodeNeedLogin         ResCode = 50003
	CodeInvalidID         ResCode = 50004

	CodeWrongPassword           ResCode = 40301
	CodeWrongUserNameOrPassword ResCode = 40302
	CodeUserNotExist            ResCode = 40401
	CodeUserExist               ResCode = 40902
	CodeEmailExist              ResCode = 40903
)

var (
	ErrUserNotLogin   = errors.New(CodeNeedLogin.Msg())
	ErrUserExist      = errors.New(CodeUserExist.Msg())
	ErrUserNotExist   = errors.New(CodeUserNotExist.Msg())
	ErrEmailExist     = errors.New(CodeEmailExist.Msg())
	ErrInvalidID      = errors.New(CodeInvalidID.Msg())
	ErrConvDataErr    = errors.New(CodeConvDataErr.Msg())
	ErrVoteTimeExpire = errors.New("投票时间已过")
)
