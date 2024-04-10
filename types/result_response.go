package types

import "time"

type ResultResponse struct {
	Hash           string    `json:"hash"`
	Numbers        []int     `json:"numbers"`
	WinningNumbers []int     `json:"winningNumbers"`
	Win            bool      `json:"win"`
	DrawDate       time.Time `json:"drawDate"`
}
