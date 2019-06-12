package controllers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/u0suke87/book-myschedule-api/models"
)

func HomeHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func AddCalender(ctx *gin.Context) {
	ctx.Request.ParseForm()

	s := models.SetSchedule(ctx)
	registerURL := template.HTML(models.CreateRegisterURL(s))
	models.CreateEvent(s)
	ctx.HTML(http.StatusOK, "thanks.tmpl", gin.H{
		"name":          s.Name,
		"startDateTime": s.StartDateTime,
		"endDateTime":   s.EndDateTime,
		"registerURL":   registerURL,
	})
}
