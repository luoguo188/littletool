package checktool

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

// fileDataInit 单个文件初始化函数
func (t *ToolData) fileDataInit(fileName string) bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	t.PmFileData = &PmFileData{}
	t.TransTab = make(map[string]string)
	if !t.unMashRoot(fileName) {
		t.ErrorPrint(fileName + ":PmFile tag  isn`t exist")
		return false
	}

	return true
}
func (t *ToolData) toolInit() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	t.Cache = make(map[string]map[string]float64, 0)
	t.ConfigData.ElementTypes = make(map[string]bool)
	t.ResultData.containPeriods = make(map[string]map[string]bool)
	t.ResultData.OkContainFiles = make(map[string]bool)
	t.ResultData.ZeroParaName = make(map[string]map[string]bool)
	t.ResultData.NilParaName = make(map[string]map[string]bool)
	return true
}

func LoopCheck() {

	cmd := ""
	for {
		tD := ToolData{}
		tD.StartCheck()
		fmt.Printf("press enter continue,put in exit to stop\n")
		if _, err := fmt.Scanf("%s\n", &cmd); err != nil {
			continue
		}
		if cmd == "exit" {
			break
		}
	}
}

func (t *ToolData) StartCheck() {

	defer func() {
		if r := recover(); r != nil {
			buff := make([]byte, 1<<10)
			runtime.Stack(buff, false)
			t.ErrorPrint(fmt.Sprintf("%v %v", r, string(buff)))
		}
	}()
	t.toolInit()
	t.ReadConfig()

	pwd, _ := os.Getwd() //获取当前目录
	if !t.initOutput() {
		t.ErrorPrint("输出文件创建失败，检查文件夹权限，系统时间格式！")
		return
	}

	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		t.ErrorPrint("read file lst fail")
		log.Fatal(err)
	}
	for i := range fileInfoList {
		fileName := fileInfoList[i].Name()
		if !t.isFileNameOk(fileName) {
			continue
		}
		if !t.fileDataInit(fileName) {
			t.ErrorPrint("one file init fail,name is:" + fileName)
			continue
		}

		if !t.unMashPmFileHeaderData() {
			continue
		}

		if !t.needCheck(fileName) {
			t.DebugPrint(fileName + " need't to check")
			continue
		}
		t.ResultData.ContainFiles = append(t.ResultData.ContainFiles, fileName)
		// 每读取一个文件刷新节点缓存
		// 解析文件节点信息
		if !t.unMashPmFileData() {
			t.ErrorPrint(fileName + "unMash node fail")
			continue
		}
		// 单个文件结果存入缓存
		if !t.setCache() {
			t.ErrorPrint(fileName + ":set Cache fail")
		}

	}

	if !t.setResult() {
		t.ErrorPrint("part of task check success")
		return
	}
	t.ErrorPrint("all file check finish")

}
