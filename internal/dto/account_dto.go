package dto

type DepositRequest struct {
	Id     uint    `json:"id" binding:"required,gt=0"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type WithdrawRequest struct {
	Id     uint    `json:"id" binding:"required,gt=0"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
