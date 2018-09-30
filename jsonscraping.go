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

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
