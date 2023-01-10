package lib

type ErrCode int32

func (ec ErrCode) Val() int32 {
	return int32(ec)
}

const (
	Err_Unmarshal ErrCode = -(iota + 10000)
	Err_Bcrypt_Gen
	Err_Acc_Exist
	Err_Parse_Token
	Err_Token_Expired
	Err_Acc_Not_Exist
	Err_Gen_Token
	Err_Base64_Decode
	Err_Bcrypt_Compare
	Err_Forbidden
	Err_Get_Users
)
