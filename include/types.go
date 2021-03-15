package include

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

var (
	Log   *logrus.Logger
	db    *gorm.DB
	dbErr error
)

type Model struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt gorm.DeletedAt
	// DeletedBy string
}

func (u *Model) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = createHash()
	return
}

type POSTTaskParameter struct {
	ParameterName  string
	ParameterValue string
}

type DBTaskParameter struct {
	Model
	DBTaskID string
	POSTTaskParameter
}

type RepeatRule struct {
}

type POSTTask struct {
	TaskName        string
	TaskDescription string
	Subject         string
	Message         string
}

type POSTRecipient struct {
	Name  string
	Email string
}

type DBRecipient struct {
	Model
	DBTaskID string
	POSTRecipient
}

type POSTSchedule struct {
	FirstRun       *time.Time
	Repeat         bool
	RepeatEvery    int
	RepeatInterval string
}

type DBTask struct {
	Model
	DBBaseScriptID string
	POSTTask
	POSTSchedule
	Enabled        bool
	Action         string
	Sender         string
	TaskParameters []DBTaskParameter
}

type POSTScriptParameter struct {
	ParameterName        string
	ParameterDescription string
}

type DBScriptParameter struct {
	Model
	DBBaseScriptID string
	POSTScriptParameter
}

type POSTBaseScript struct {
	Name        string
	Description string
}

type DBBaseScript struct {
	Model
	POSTBaseScript
	ScriptParameters []DBScriptParameter
	ScriptFiles      []DBScriptFile
	Task             []DBTask
}

type DBScriptFile struct {
	Model
	ScriptFile     bool
	DBBaseScriptID string
	FileName       string
	PathToFile     string
}

type POSTJobDone struct {
	Files      []string
	DurationMs int
}

type DBJob struct {
	Model
	DBBaseScriptID string
	DBTaskID       string
	DBScriptFileID string
	Source         string
	CommandString  string
	CommandOutput  string
	Error          string
	Reports        []DBReport
	DurationMs     int
}

type DBReport struct {
	Model
	DBJobID   string
	FileName  string
	OpenCount int
}
