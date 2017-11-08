package logger

import (
	"log"
	"os"
	"PointSystem/conf"
	"runtime"
	"strings"
	"strconv"
)

var l *log.Logger

func Init() {
	logFile, err := os.OpenFile(conf.LogPath, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	l = log.New(logFile, "", log.LstdFlags)
}

func Debug(v ...interface{}) {
	fl := getFileLine()
	l.SetPrefix("DEBUG " + fl)
	l.Println(v)
}

func Info(v ...interface{}) {
	fl := getFileLine()
	l.SetPrefix("INFO " + fl)
	l.Println(v)
}

func Error(v ...interface{}) {
	fl := getFileLine()
	l.SetPrefix("ERROR " + fl)
	l.Println(v)
}

func Fatal(v ...interface{}) {
	fl := getFileLine()
	l.SetPrefix("FATAL " + fl)
	l.Fatalln(v)
}

func getFileLine() (fileLine string) {
	_, file, line, ok := runtime.Caller(4)
	if ok {
		filePart := strings.Split(file, "src/")
		fileLine = filePart[1] + ":" + strconv.Itoa(line) + " "
	}
	return
}