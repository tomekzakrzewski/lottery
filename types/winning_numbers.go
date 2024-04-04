package types

import "time"

type WinningNumbers struct {
	Numbers  []int     `json:"numbers" bson:"numbers"`
	DrawDate time.Time `json:"drawDate" bson:"drawDate"`
}
