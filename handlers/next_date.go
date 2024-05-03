package handlers

import (
	"go_final_project/services"
	"go_final_project/utils"
	"net/http"
	"time"
)

func HandleNextDate(res http.ResponseWriter, req *http.Request) {
	nowDate, _ := time.Parse(utils.DateLayout, req.URL.Query().Get("now"))
	taskDate := req.URL.Query().Get("date")
	repeatRule := req.URL.Query().Get("repeat")
	nextDate, _ := services.NextDate(nowDate, taskDate, repeatRule)
	res.Write([]byte(nextDate))
}
