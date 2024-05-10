package main

import (
	"fmt"
	"time"
)

func ParseTime(str string) time.Time {
	layout := "15:04"
	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println("Ошибка при преобразовании строки во время ", err)
	}
	return t

}
