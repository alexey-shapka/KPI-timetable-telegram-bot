package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func getinfo(group string, week int, day int) [][]string {
	var result [][]string

	groupLink := "https://api.rozklad.org.ua/v2/groups/" + group + "/lessons?filter={%22lesson_week%22:" +
		strconv.Itoa(week) + ",%20%22day_number%22:" + strconv.Itoa(day) + "}"

	response, _ := http.Get(groupLink)
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)

	type BlockInfo struct {
		Day_number    string
		Lesson_number string
		Lesson_name   string
		Lesson_type   string
		Lesson_room   string
		Teacher_name  string
		Day_name      string
	}

	type Link struct {
		Data []BlockInfo
	}

	var link Link

	json.Unmarshal([]byte(contents), &link)
	jsonarray := link.Data

	for i, _ := range jsonarray {
		info := []string{jsonarray[i].Day_name, jsonarray[i].Day_number, jsonarray[i].Lesson_number, jsonarray[i].Lesson_name, jsonarray[i].Lesson_type,
			jsonarray[i].Lesson_room, jsonarray[i].Teacher_name}
		result = append(result, info)
	}

	if len(result) == 0 {
		result = append(result, []string{strconv.Itoa(day)})
	}

	return result
}

func getRating(group string) string {

	link := "https://api.rozklad.org.ua/v2/groups/" + group + "/lessons"
	response, _ := http.Get(link)
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)

	type Teacher struct {
		Teacher_rating string
		Teacher_name   string
	}

	type Teachers struct {
		Teachers []Teacher
	}

	type All struct {
		Data []Teachers
	}

	var all All
	json.Unmarshal([]byte(contents), &all)
	jsonarray := all.Data

	resultString := "*Рейтинг викладачів " + strings.ToUpper(group) + ":*\n"

	teacherMap := make(map[string]string)
	var score []string
	for i, _ := range jsonarray {
		j := 0
		for range jsonarray[i].Teachers {
			if contains(score, jsonarray[i].Teachers[j].Teacher_rating) == false {
				score = append(score, jsonarray[i].Teachers[j].Teacher_rating)
			}
			teacherMap[jsonarray[i].Teachers[j].Teacher_name] = jsonarray[i].Teachers[j].Teacher_rating
			j++
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(score)))

	for i, _ := range score {
		for k, v := range teacherMap {
			if score[i] == v {
				resultString += k + " " + v + "\n"
			}
		}
	}

	return fmt.Sprintf(resultString)
}

func getTeacherschedule(teacherName string) string {

	link := "https://api.rozklad.org.ua/v2/teachers/" + teacherName + "/lessons"

	response, _ := http.Get(link)
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)

	type GroupInfo struct{
		Group_full_name string
	}

	type TeacherInfo struct{
		Day_number    string
		Lesson_number string
		Lesson_name   string
		Lesson_type   string
		Lesson_room   string
		Day_name      string
		Lesson_week   string
		Groups []GroupInfo
	}

	type Link struct {
		Data []TeacherInfo
	}

	var getLink Link
	json.Unmarshal([]byte(contents), &getLink)
	jsonarray := getLink.Data

	week := make(map[string]int)
	week["1"] = 0
	week["2"] = 0
	week["3"] = 0
	week["4"] = 0
	week["5"] = 0
	week["6"] = 0
	

	resultString := ""
	check := 0

	for i, _ := range jsonarray {
		groups := ""
		for j, _ := range jsonarray[i].Groups{
			groups += jsonarray[i].Groups[j].Group_full_name + " "
		}

		if jsonarray[i].Lesson_week == "2" && check == 0{
			week["1"] = 0
			week["2"] = 0
			week["3"] = 0
			week["4"] = 0
			week["5"] = 0
			week["6"] = 0
			check++
			resultString+="\n"
		}
			

		if week[jsonarray[i].Day_number] == 0{
			resultString += fmt.Sprintf("*%s:*\n",jsonarray[i].Day_name)
			week[jsonarray[i].Day_number] ++
		}
			
		resultString+=fmt.Sprintf("%s) %s %s %s %s\n",jsonarray[i].Lesson_number, jsonarray[i].Lesson_name, jsonarray[i].Lesson_type,
		jsonarray[i].Lesson_room, groups)
	}

	return resultString
}


func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
