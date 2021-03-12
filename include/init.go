package include

import "os"

func InitApplication(f *os.File) {

	initLog(f)
	initDatabase()

}
