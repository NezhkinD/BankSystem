package dto

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
