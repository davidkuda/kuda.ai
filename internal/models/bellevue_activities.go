package models

import (
	"database/sql"
	"time"
)

type BellevueActivities []*BellevueActivity

func NewBellevueActivity() *BellevueActivity {
	return &BellevueActivity{
		Date: time.Now(),
	}
}

// This struct keeps track of things you need to pay for.
type BellevueActivity struct {
	Date         time.Time
	Breakfasts   int
	LunchDinners int
	Coffees      int
	Saunas       int
	Comment      string
	TotalCost    int
}

type BellevueActivityModel struct {
	DB *sql.DB
}
