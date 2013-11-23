/**
 * User: jackong
 * Date: 11/5/13
 * Time: 4:50 PM
 */
package global

import (
	"github.com/gorilla/mux"
	"github.com/Jackong/log"
	"os"
	"fmt"
	"morning-dairy/config"
	"time"
	"github.com/Jackong/db"
	_ "github.com/Jackong/db/mysql"
	"github.com/Jackong/log/writer"
	"net/http"
	"morning-dairy/err"
)

var (
	GoPath string
	Now func() string
	Time func() time.Time
	Today func() string
	Project config.Config
	Access *accessLog
	Log log.Logger
	Conn db.Database
	Router *router
)

func init() {
	fmt.Println("init env...")
	baseEnv()

	fmt.Println("loading config...")
	loadConfig()

	fmt.Println("opening db...")
	openDb()

	fmt.Println("init router...")
	Router = &router{mux.NewRouter()}
	Router.NotFoundHandler = &notFoundHandler{}
	initLog()
}

type notFoundHandler struct {

}

func (this *notFoundHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	panic(err.AccessError{Status: http.StatusNotFound, Msg: "page not found"})
}

func initLog() {
	today := Today()
	fmt.Println("getting mail log...")
	mailLog := newDateLog(today, mailLog)

	fileLevel := int(Project.Get("log", "file", "level").(float64))

	fmt.Println("getting access log...")
	Access = &accessLog{logger: log.MultiLogger(newDateLog(today, func(date string) log.Logger {
				return fileLog("access", date, fileLevel)
			}), mailLog)}


	fmt.Println("getting action log...")
	actionLog := newDateLog(today, func(date string) log.Logger {
			return fileLog("action", date, fileLevel)
		})

	Log = log.MultiLogger(actionLog, mailLog)

	OnBefore(logBefore)
}

func logBefore(writer http.ResponseWriter, req *http.Request) (err error) {
	Access.Info(0, req, "access")
	return
}

func baseEnv() {
	Time = func() time.Time {
		return time.Now()
	}

	Now = func() string {
		return Time().Format(FORMAT_DATE_TIME)
	}

	Today = func() string {
		return Time().Format(FORMAT_DATE)
	}

	GoPath = os.Getenv("GOPATH")
}

func loadConfig() {
	Project = config.NewConfig(GoPath  + "/src/morning-dairy/config/project.json")
}

func fileLog(dir, date string, level int) log.Logger {
	dir = Project.String("log", "dir") + "/" +  dir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	file, err := os.OpenFile(dir + "/" + date + ".log", os.O_RDWR | os.O_CREATE | os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		ShutDown()
	}
	logger := log.NewLogger(file, Project.String("server", "name"), level)
	return logger
}

func mailLog(date string) log.Logger {
	mail := &writer.Email{
		User: Project.String("log", "email", "user"),
		Password: Project.String("log", "email", "password"),
		Host : Project.String("log", "email", "host"),
		To : Project.String("log", "email", "to"),
		Subject: Project.String("log", "email", "subject"),
	}

	server := Project.String("server", "name")
	mailLevel := int(Project.Get("log", "email", "level").(float64))
	check := Project.Get("formal").(bool)
	if check {
		if err := mail.SendMail("starting server " + server + "..."); err != nil {
			fmt.Println(err)
			ShutDown()
		}
		return log.NewLogger(&asyncMail{mail}, server, mailLevel)
	}
	return fileLog("email", date, mailLevel)
}

func openDb() {
	settings := db.DataSource{
		Socket:   Project.String("db", "socket"),
		Database: Project.String("db", "database"),
		User:     Project.String("db", "user"),
		Password: Project.String("db", "password"),
	}

	var err error
	if Conn, err = db.Open("mysql", settings); err != nil {
		fmt.Println(err)
		ShutDown()
	}

	OnShutDown(func() {
		Conn.Close()
	})
}
