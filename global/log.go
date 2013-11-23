/**
 * User: jackong
 * Date: 11/6/13
 * Time: 4:28 PM
 */
package global

import (
	"net/http"
	"sync"
	"github.com/Jackong/log/writer"
	"github.com/Jackong/log"
)

type asyncMail struct {
	*writer.Email
}


func (this *asyncMail) Write(p []byte) (n int, err error) {
	go this.Email.Write(p)
	return
}


type dateLog struct {
	mu     sync.Mutex
	date string
	getLog func(string) log.Logger
	log.Logger
}

func (this *dateLog) Output(level, depth int, s string) {
	this.ensureDate()
	this.Logger.Output(level, depth + 1, s)
}

func (this *dateLog) ensureDate() {
	this.mu.Lock()
	defer this.mu.Unlock()
	today := Today()
	if this.date != today {
		this.date = today
		this.Logger = this.getLog(this.date)
	}
}

func newDateLog(date string, getLog func(date string) log.Logger) *dateLog {
	return &dateLog{date: date, Logger: getLog(date)}
}

type accessLog struct {
	logger log.Logger
}

func (this *accessLog) Debug(code int, req *http.Request, v ... interface {}) {
    this.logger.Debug(this.logAppend(code, req, v)...)
}

func (this *accessLog) Info(code int, req *http.Request, v ... interface {}) {
	this.logger.Info(this.logAppend(code, req, v)...)
}

func (this *accessLog) Config(code int, req *http.Request, v ... interface {}) {
	this.logger.Config(this.logAppend(code, req, v)...)
}

func (this *accessLog) Warn(code int, req *http.Request, v ... interface {}) {
	this.logger.Warn(this.logAppend(code, req, v)...)
}

func (this *accessLog) Error(code int, req *http.Request, v ... interface {}) {
	this.logger.Error(this.logAppend(code, req, v)...)
}

func (this *accessLog) Alert(code int, req *http.Request, v ... interface {}) {
	this.logger.Alert(this.logAppend(code, req, v)...)
}

func (this *accessLog) Fatal(code int, req *http.Request, v ... interface {}) {
	this.logger.Fatal(this.logAppend(code, req, v)...)
}

func (this *accessLog) logAppend(code int, req *http.Request, v ... interface {}) (more []interface {}) {
	const (
		DELIMITER = "|"
	)
	more = append(more,
		code, DELIMITER,
		req.RemoteAddr, DELIMITER,
		req.Method, DELIMITER,
		req.URL.Path, DELIMITER,
		req.UserAgent(), DELIMITER)
	more = append(more, v...)
	return
}
