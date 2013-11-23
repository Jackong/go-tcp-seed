/**
 * User: jackong
 * Date: 11/6/13
 * Time: 4:28 PM
 */
package global

import (
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
