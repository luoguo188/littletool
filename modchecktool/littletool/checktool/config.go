package checktool

import (
	"fmt"
	"strconv"
	"strings"
)

func (t *ToolData) setDefaultConfig() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	t.ConfigData.ElementTypeFlag = true
	t.ConfigData.UserLabelFlag = true
	t.ConfigData.ObjectTypeFlag = true
	t.ConfigData.PeriodFlag = false
	t.DebugPrint("set default Config success")
	return true
}

func (t *ToolData) ReadDebugMode() int {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return readConfigNotOk
	}
	var debugMode uint8
	fmt.Println("please set the debug mode,1 mean true,0 mean false" +
		"(if put in only enter,will skip debug mode,select mode)")
	_, err := fmt.Scanf("%d\n", &debugMode)
	if err != nil {
		if ok := t.setDefaultConfig(); !ok {
			t.ErrorPrint("set config fail")
			return readConfigNotOk
		}
		return readConfigSkip
	}
	if debugMode == 1 {
		t.ConfigData.DebugMode = true
		t.DebugPrint("debug mode on")
	}
	return readConfigOk
}

func (t *ToolData) ReadSelectLabel() int {
	t.DebugPrint("enter ReadSelectLabel")
	var selectLabel []byte
	fmt.Println("1-网元ID（UserLabel），2-模块类型（ObjectType）, 3-网元形态（ElementType）,4-周期（Period）\n" +
		"please set the select Label,such as \"123\", default is \"123\"")
	_, err := fmt.Scanf("%s\n", &selectLabel)
	if err != nil {
		if ok := t.setDefaultConfig(); !ok {
			t.ErrorPrint("set config fail")
			return readConfigNotOk
		}
		t.DebugPrint("skip select labels")
		return readConfigSkip
	}
	if len(selectLabel) < 1 {
		if ok := t.setDefaultConfig(); !ok {
			t.ErrorPrint("set config fail")
			return readConfigNotOk
		}
	}
	//根据读到的字符决定目前使用哪些标签作为索引，1-网元ID（UserLabel），2-模块类型（ObjectType）, 3-网元形态（ElementType）,4-周期（Period）
	for _, label := range selectLabel {
		switch label {
		case '1':
			t.ConfigData.UserLabelFlag = true
		case '2':
			t.ConfigData.ObjectTypeFlag = true
		case '3':
			t.ConfigData.ElementTypeFlag = true
		case '4':
			t.ConfigData.PeriodFlag = true
		default:

		}
	}
	t.DebugPrint("select label is " + string(selectLabel))
	return readConfigOk
}

func (t *ToolData) ReadTime() int {
	t.DebugPrint("enter")
	fmt.Println("please set the ElementTypes,such as 1000-1045," +
		"Indicates that the file between 10:15 and 10:45 will be processed")
	var timeString string
	_, err := fmt.Scanf("%s\n", &timeString)
	if err != nil {
		t.DebugPrint("ReadTime config skip")
		return readConfigSkip
	}
	if len(timeString) < 9 {
		t.DebugPrint("timeString length is wrong,timeString is " + timeString)
		return readConfigNotOk
	}
	var hour, minute int
	if hour, err = strconv.Atoi(timeString[:2]); err != nil {
		t.DebugPrint("begin hour is not num")
		return readConfigNotOk
	}
	minute, err = strconv.Atoi(timeString[2:4])
	if err != nil {
		t.DebugPrint("begin minute is not num")
		return readConfigNotOk
	}
	t.ConfigData.selectTimeBegin = hour*60 + minute
	if hour, err = strconv.Atoi(timeString[5:7]); err != nil {
		t.DebugPrint("end hour is not num")
		return readConfigNotOk
	}
	minute, err = strconv.Atoi(timeString[7:9])
	if err != nil {
		t.DebugPrint("end minute is not num")
		return readConfigNotOk
	}
	t.ConfigData.selectTimeEnd = hour*60 + minute
	t.DebugPrint("set period is " + timeString + "convent to minute num is " +
		strconv.Itoa(t.ConfigData.selectTimeBegin) + "-" + strconv.Itoa(t.ConfigData.selectTimeEnd))
	return readConfigOk
}

func (t *ToolData) ReadElementType() int {
	t.DebugPrint("enter")
	fmt.Println("please set the select ElementTypes,such as udm,hss use \",\" split")
	var ElementTypeString string
	_, err := fmt.Scanf("%s\n", &ElementTypeString)
	if err != nil {
		t.DebugPrint("select ElementTypes config skip")
		return readConfigSkip
	}
	ElementTypes := strings.Split(ElementTypeString, ",")
	for _, v := range ElementTypes {
		t.ConfigData.ElementTypes[v] = true
	}
	return readConfigOk
}

func (t *ToolData) ReadConfig() bool {
	t.DebugPrint("enter")
	ReadStatus := t.ReadDebugMode()
	if ReadStatus == readConfigSkip {
		t.InfoPrintResult("skip the debug mode")
		return true
	}
	if ReadStatus == readConfigNotOk {
		t.ErrorPrint("skip debug mode and select labels,but set default config fail")
		return false
	}

	ReadStatus = t.ReadSelectLabel()
	if ReadStatus == readConfigNotOk {
		t.ErrorPrint("skip select labels,but set default config fail")
		return false
	}

	ReadStatus = t.ReadTime()
	if ReadStatus == readConfigNotOk {
		t.ConfigData.selectTimeBegin = 0
		t.ConfigData.selectTimeEnd = 0
		t.DebugPrint("default use all period files")
	}

	ReadStatus = t.ReadElementType()

	return true
}
