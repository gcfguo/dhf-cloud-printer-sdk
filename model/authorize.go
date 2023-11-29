package model

type (
	VerifyReq struct {
		Email string `json:"email"`
	}
	VerifyRes struct {
		VerifyCode string `json:"verify_code"`
	}
)

type (
	SignInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	SignInRes struct {
		TokenValue string `json:"token_value"`
		TokenType  string `json:"token_type"`
	}
)

type (
	SignUpReq struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		VerifyCode string `json:"verify_code"`
	}
	SignUpRes struct {
	}
)
