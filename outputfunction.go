package main

import (
	"fmt"
)

func getDay(data[][]string) string{
	day := ""
	if len(data[0])>1 {
		day+=fmt.Sprintf("*%s:*\n",data[0][0])
	for i, _ := range data{
		day+= fmt.Sprintf("%s) %s %s %s %s\n", data[i][2],data[i][3],data[i][6],data[i][4],data[i][5])
	}
		return day

	}	else {
		
		switch data[0][0]{
		case "0": 
			day+=fmt.Sprintf("*%s:*\n","Неділя")
		case "1":
			day+=fmt.Sprintf("*%s:*\n","Понеділок")
		case "2":
			day+=fmt.Sprintf("*%s:*\n","Вівторок")
		case "3":
			day+=fmt.Sprintf("*%s:*\n","Середа")
		case "4":
			day+=fmt.Sprintf("*%s:*\n","Четвер")
		case "5":
			day+=fmt.Sprintf("*%s:*\n","П'ятниця")
		case "6":
			day+=fmt.Sprintf("*%s:*\n","Субота")
		case "7":
			day+=fmt.Sprintf("*%s:*\n","Неділя")
		}
		
		day += "Вихідний."
	}

	return day
}