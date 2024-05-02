package handlers

import (
	"fmt"
	"go_final_project/constants"
	"go_final_project/services"
	"net/http"
	"time"
)

func HandleNextDate(res http.ResponseWriter, req *http.Request) {
	nowDate, _ := time.Parse(constants.DateLayout, req.URL.Query().Get("now"))
	taskDate := req.URL.Query().Get("date")
	repeatRule := req.URL.Query().Get("repeat")
	nextDate, _ := services.NextDate(nowDate, taskDate, repeatRule)
	fmt.Printf("Now %s, taskDate %s, repeat %s, nextDate %s\n", nowDate.Format(constants.DateLayout), taskDate, repeatRule, nextDate)
	res.Write([]byte(nextDate))
}
