package checktool

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

func (t *ToolData) initOutput() bool {

	timeNow := time.Now().Format("2006-01-02[15,04,05]")
	t.ResultData.logName = timeNow + "result.txt"
	f, err := os.Create(t.ResultData.logName)

	defer f.Close()

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = f.Write([]byte("结果仅供参考，欢迎找虫子\n"))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (t *ToolData) ErrorPrint(text string) {
	if t == nil {
		return
	}
	t.PrintLog(text)
	fmt.Println(text)
}

func (t *ToolData) DebugPrint(text string) {
	if t == nil {
		return
	}
	if t.ConfigData.DebugMode {
		t.PrintLog(text)
	}
}

func (t *ToolData) InfoPrintResult(text string) {
	fmt.Println(text)
	t.printResult(text)
}

func (t *ToolData) PrintLog(text string) {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	defer f.Close()
	text += "\n"
	pc, fileName, line, ok := runtime.Caller(2)
	if ok {
		funName := runtime.FuncForPC(pc)
		text = fileName + ":" + strconv.Itoa(line) + "(" + funName.Name() + ")" + text
	}
	text = time.Now().Format("[2006-01-02 15:04:05]") + text
	if t.ConfigData.DebugMode {
		text = "Debug" + text
	} else {
		text = "Error" + text
	}

	if err != nil {
		fmt.Println(err.Error())

	} else {
		_, err = f.Write([]byte(text))
		if err != nil {
			panic("write failed")
		}
	}
}

func (t *ToolData) printResult(text string) {
	f, err := os.OpenFile(t.ResultData.logName, os.O_RDWR|os.O_APPEND, 0666)

	defer f.Close()

	if err != nil {
		fmt.Println(err.Error())

	} else {
		text = text + "\n"
		_, err = f.Write([]byte(text))
		if err != nil {
			panic("write failed")
		}
	}
}
