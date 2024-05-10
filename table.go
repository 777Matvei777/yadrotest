package main

import "time"

type Table struct {
	StartTime  time.Time
	CostPerDay int
	IsBusy     bool
	TotalTime  time.Time
}

func (t *Table) CountingCostAndTime(EndTime time.Time, PriceForHour int) {
	duration := EndTime.Sub(t.StartTime)
	min := int(duration.Minutes())
	hours := min / 60
	if min%60 != 0 {
		hours++
	}
	t.CostPerDay += PriceForHour * hours
	t.TotalTime = t.TotalTime.Add(duration)
}

func (t *Table) TakeTheTable(StartTime time.Time) {
	t.StartTime = StartTime
	t.IsBusy = true
}
