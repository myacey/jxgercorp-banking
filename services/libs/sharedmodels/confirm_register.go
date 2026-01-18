package sharedmodels

type RegisterConfirmMsgEmail struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	ConfirmCode string `json:"confirm_code"`
}
