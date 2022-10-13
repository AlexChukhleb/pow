package logger

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func InitLog() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			f := frame.Function
			ss := strings.Split(f, ".")
			f = ss[len(ss)-1]
			fl := path.Base(frame.File)
			fl, _, _ = strings.Cut(fl, ".")
			fileName := "(" + fl + ":" + strconv.Itoa(frame.Line) + ":" + f + ")"
			return "", fileName
		},
	})

	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)

	{
		dlevel := os.Getenv("debuglevel")
		lvl, _ := strconv.Atoi(dlevel)
		if lvl > 0 {
			log.SetLevel(log.Level(lvl))
		}
	}
}
