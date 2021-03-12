package include

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Log   *logrus.Logger
	db    *gorm.DB
	dbErr error
)

type DBScryptParameter struct {
	gorm.Model
	DBBaseScriptID uint
	ParameterName  string
	ParameterValue string
	Runtime        bool
	DefaultValue   string
}

type DBScheduling struct {
	gorm.Model
	DBBaseScriptID uint

	DayOfMonth  int
	DayOfWeek   int
	WeekOfMonth int

	Enabled bool
}

type DBBaseScript struct {
	gorm.Model
	Hash           string
	ScryptFolder   string
	MainScriptFile string
	Enabled        bool

	ScriptParameters []DBScryptParameter
	ScriptFiles      []DBScriptFile
	ScriptScheduling []DBScheduling
}

type DBScriptFile struct {
	gorm.Model
	DBBaseScriptID uint
	FileName       string
}

type DBScriptHistory struct {
	gorm.Model
	DBBaseScriptID uint
	Action         string
	Source         string
	Object         string
}
