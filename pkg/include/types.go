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

type getResponse struct {
	Total    int64
	Current  int
	Entities interface{}
}

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
	Enabled         bool
	//Action          string
	Sender string
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
	TaskParameters []DBTaskParameter
	Jobs           []DBJob
	Recipients     []DBRecipient
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

type POSTJobDoneFile struct {
	FileName   string
	ReportName string
}

type POSTJobDone struct {
	Files      []POSTJobDoneFile
	DurationMs int
	ResultType string
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
	Mails          []DBOutgoingMails
	DurationMs     int
	TestRun        bool
}

type DBReportDownloadHistory struct {
	Model
	DBReportID    string
	DBRecipientID string
}

type DBReport struct {
	Model
	DBJobID string
	POSTJobDoneFile
	OpenCount       int
	DownloadHistory []DBReportDownloadHistory
}

type DBOutgoingMails struct {
	Model
	DBJobID    string
	ToEmail    string
	Subject    string
	Message    string
	OutMessage string
	Status     string
	History    []DBOutgoingMailHistory
}

type DBOutgoingMailHistory struct {
	Model
	DBOutgoingMailsID string `gorm:"index"`
	RecType           string
	HistoryMessage    string
}
