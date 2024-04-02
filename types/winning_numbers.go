package types

import "time"

type WinningNumbers struct {
	Numbers  []int     `json:"numbers" bson:"numbers"`
	DrawDate time.Time `json:"draw_date" bson:"draw_date"`
}
