package controllers

import (
	"CedAppDeliver/utils"
	"time"

	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(func() {
		InitDB()
		InitLogger()
		jobs.Schedule("0 0 0 * * ?", jobs.Func(utils.CleanOutdateApps))
		now := time.Now()
		logger.Println("CedAppDeliver start at " + now.Format("2006-01-02 15:04:05"))
	})

	revel.TemplateFuncs["toTimeString"] = func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	}
	// revel.InterceptMethod((*GormController).Begin, revel.BEFORE)
	// revel.InterceptMethod((*GormController).Commit, revel.AFTER)
	// revel.InterceptMethod((*GormController).Rollback, revel.FINALLY)

	// logger.Println("hello")
	// logger.Println("oh....")
	// logger.Fatal("test")
	// logger.Fatal("test2")
}
