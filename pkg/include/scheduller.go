package include

import (
	"github.com/go-co-op/gocron"
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

	Sch.StartAsync()

}

func RunScheduler(task *DBTask) {

	var jb *gocron.Job
	var err error

	switch task.RepeatInterval {
	case "hour":
		jb, err = Sch.Every(task.RepeatEvery).Hour().StartAt(*task.FirstRun).Do(StartJob, task, "timer", false)
	case "minutes":
		jb, err = Sch.Every(task.RepeatEvery).Minute().StartAt(*task.FirstRun).Do(StartJob, task, "timer", false)
	case "minute":
		jb, err = Sch.Every(task.RepeatEvery).Minute().StartAt(*task.FirstRun).Do(StartJob, task, "timer", false)
	case "day":
		jb, err = Sch.Every(task.RepeatEvery).Day().StartAt(*task.FirstRun).Do(StartJob, task, "timer", false)
	case "month":
		jb, err = Sch.Every(task.RepeatEvery).Month(task.FirstRun.Day()).StartAt(*task.FirstRun).Do(StartJob, task, "timer")
	}

	if err != nil {
		Log.Error(err.Error())
	} else {
		jb.Tag(task.ID)
		Log.Tracef("Task is run")
	}

}
