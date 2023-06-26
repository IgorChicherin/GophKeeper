package models

type RequestUserModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseUserLogin struct {
	Token string `json:"token"`
	Cert  string `json:"cert"`
}
