package include

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
)

var Sch *gocron.Scheduler

func ApplicationStartAllTasks() {

	Log.Trace("Open all active tasks and run it")

	var tasks []DBTask
	db.Where("enabled = ?", true).Find(&tasks)
	for _, task := range tasks {
		Log.Tracef("run task (%s) '%s' which is enabled '%t' ", task.ID, task.TaskName, task.Enabled)
		RunScheduler(&task)
	}

	var jb *gocron.Job
	jb, _ = Sch.Every(time.Hour * 24).Do(SyncBuildings)
	fmt.Println(jb.NextRun())

	Sch.StartAsync()
	Log.Info("Scheduler is running : ", Sch.IsRunning(), "\n")
}

func RunScheduler(task *DBTask) {

	var jb *gocron.Job
	var err error

	switch task.RepeatInterval {
	case "hour":
		jb, err = Sch.Every(task.RepeatEvery).Hour().StartAt(*task.FirstRun).Do(StartJob, *task, "timer", false)
	case "minutes":
		jb, err = Sch.Every(task.RepeatEvery).Minute().StartAt(*task.FirstRun).Do(StartJob, *task, "timer", false)
	case "minute":
		jb, err = Sch.Every(task.RepeatEvery).Minute().StartAt(*task.FirstRun).Do(StartJob, *task, "timer", false)
	case "day":
		jb, err = Sch.Every(task.RepeatEvery).Day().StartAt(*task.FirstRun).Do(StartJob, *task, "timer", false)
	case "month":
		day := task.FirstRun.Day()
		jb, err = Sch.Every(task.RepeatEvery).Months(day).Do(StartJob, *task, "timer", false)
	}

	if err != nil {
		Log.Error(err.Error())
	} else {
		jb.Tag(task.ID)
		Log.Info("Task, first run ", task.FirstRun, " next run ", jb.NextRun(), "task id : ", task.ID)
	}

}
