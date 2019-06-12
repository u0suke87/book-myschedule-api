package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Schedule struct {
	Name          string
	Purpose       string
	Summary       string
	StartDateTime string
	EndDateTime   string
}

func SetSchedule(ctx *gin.Context) *Schedule {
	name := ctx.Request.FormValue("name")
	purpose := ctx.Request.FormValue("purpose")
	summary := ctx.Request.FormValue("name") + " " + ctx.Request.FormValue("purpose")
	startDateTime := ctx.Request.FormValue("day_s") + "T" + ctx.Request.FormValue("time_s") + "+09:00"
	endDateTime := ctx.Request.FormValue("day_e") + "T" + ctx.Request.FormValue("time_e") + "+09:00"
	return &Schedule{
		Name:          name,
		Purpose:       purpose,
		Summary:       summary,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
	}
}

func CreateRegisterURL(s *Schedule) string {
	var strTime string
	var sTime, eTime []string
	slice := strings.Split(s.StartDateTime, "-")
	strTime = strings.Join(slice, "")
	slice = strings.Split(strTime, ":")
	strTime = strings.Join(slice, "")
	sTime = strings.Split(strTime, "+")

	slice = strings.Split(s.EndDateTime, "-")
	strTime = strings.Join(slice, "")
	slice = strings.Split(strTime, ":")
	strTime = strings.Join(slice, "")
	eTime = strings.Split(strTime, "+")

	Time := sTime[0] + "/" + eTime[0]
	URL := "http://www.google.com/calendar/event?action=TEMPLATE&text=" + s.Purpose + "&dates=" + Time
	return URL
}

func CreateEvent(s *Schedule) {
	b, err := ioutil.ReadFile("client-secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// Scope changed to book schedule by u0suke87 : CalendarReadonlyScope → CalendarScope
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	//call func isValidTime
	//t := time.Now().Format(time.RFC3339)

	_, err = srv.Events.Insert("primary", createEventData(s)).Do()
	if err != nil {
		log.Fatalf("Sorry. Unable to book your events: %v", err)
	}
}

func createEventData(schedule *Schedule) *calendar.Event {
	event := &calendar.Event{
		Summary: schedule.Summary,
		//Location: schedule.Location,
		Start: &calendar.EventDateTime{
			DateTime: schedule.StartDateTime,
			TimeZone: "Asia/Tokyo",
		},
		End: &calendar.EventDateTime{
			DateTime: schedule.EndDateTime,
			TimeZone: "Asia/Tokyo",
		},
	}

	return event
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
