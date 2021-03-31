package checktool

import (
	"fmt"
	"strconv"
)

func (t *ToolData) getSelectLabel(userLabel string) string {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return ""
	}
	selectLabel := ""
	if t.ConfigData.UserLabelFlag {
		t.DebugPrint("selectLabel include UserLabelFlag")
		selectLabel = userLabel
	}
	if t.ConfigData.ElementTypeFlag {
		t.DebugPrint("selectLabel include ElementTypeFlag")
		selectLabel = selectLabel + "-" + t.getElementType()
	}
	if t.ConfigData.ObjectTypeFlag {
		t.DebugPrint("selectLabel include ObjectTypeFlag")
		selectLabel = selectLabel + "-" + t.getObjectType()
	}
	if t.ConfigData.PeriodFlag {
		t.DebugPrint("selectLabel include PeriodFlag")
		selectLabel = selectLabel + "-" + t.getPeriod()
	}
	return selectLabel
}

func (t *ToolData) isSelectTime() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	if (t.ConfigData.selectTimeEnd == 0) && (t.ConfigData.selectTimeBegin == 0) {
		t.DebugPrint("selectTime is default")
		return true
	}
	timeMinute := t.getStartHourMinute()
	if timeMinute == invalidMinuteTime {
		t.ErrorPrint("StartTime get failed")
		return false
	}

	if t.ConfigData.selectTimeBegin < t.ConfigData.selectTimeEnd {
		if (timeMinute >= t.ConfigData.selectTimeBegin) && (timeMinute <= t.ConfigData.selectTimeEnd) {
			return true
		}
		return false
	} else {
		if timeMinute > t.ConfigData.selectTimeEnd && timeMinute < t.ConfigData.selectTimeBegin {
			return false
		}
		return true
	}
}

func (t *ToolData) isSelectType() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	if len(t.ConfigData.ElementTypes) == 0 {
		t.DebugPrint("select ElementTypes is nil,default use all type")
		return true
	}
	ElementType := t.getElementType()
	if ElementType == "" {
		t.DebugPrint(t.fileName + " file get ElementType fail")
		return false
	}
	if _, ok := t.ConfigData.ElementTypes[ElementType]; !ok {
		t.DebugPrint(t.fileName + " file ElementType isn't in the select ElementTypes, ElementType: " +
			ElementType + "select ElementTypes:" + fmt.Sprintf("%v", t.ConfigData.ElementTypes))
		return false
	}
	return true
}

func (t *ToolData) needCheck(fileName string) bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	if !t.isSelectTime() {
		t.DebugPrint(t.fileName + " file isn't in the select time")
		return false
	}
	if !t.isSelectType() {
		t.DebugPrint(t.fileName + " file isn't in the select ElementTypes")
		return false
	}
	t.ResultData.TotalFileNum++
	return true
}

func (t *ToolData) isFileNameOk(fileName string) bool {
	if len(fileName) < 4 || fileName[len(fileName)-4:] != ".xml" {
		return false
	}
	return true
}

func (t *ToolData) setCache() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	if t.MeasureData.PmData == nil {
		t.ErrorPrint("file:" + t.fileName + "PmData is nil")
		return false
	}
	for _, object := range t.MeasureData.PmData.SelectElements("Object") {
		if object == nil {
			t.ErrorPrint("file:" + t.fileName + " object is nil")
			continue
		}
		userLabel := object.SelectAttrValue("UserLabel", "unknown")

		selectLabel := t.getSelectLabel(userLabel)
		if selectLabel == "" {
			t.ErrorPrint("file:" + t.fileName + " selectLabel is nil,the userLabel is " + userLabel)
			continue
		}

		_, ok := t.Cache[selectLabel]
		if !ok {
			t.DebugPrint("make new cache ,key is" + selectLabel)
			t.Cache[selectLabel] = make(map[string]float64)
		}
		filePeriod := t.getStartTime()
		if filePeriod != "" {
			if t.ResultData.containPeriods[selectLabel] == nil {
				t.ResultData.containPeriods[selectLabel] = make(map[string]bool)
			}
			t.ResultData.containPeriods[selectLabel][filePeriod] = true
		}

		for _, name := range object.SelectElements("V") {
			index := name.SelectAttrValue("i", "unknown")
			paraName, ok := t.TransTab[index]
			if !ok {
				t.ErrorPrint("userLabel:" + selectLabel + "para:" + index + "don't include the PmName,use the numKey")
				paraName = index
			}
			valueNum, err := strconv.ParseFloat(name.Text(), 64)
			if err != nil {
				if !t.initNilParaList(selectLabel, paraName) {
					t.ErrorPrint("make map nilParaList fail the selectLabel is " + selectLabel + "," +
						"paraName is " + paraName)
				}
				t.ResultData.NilParaName[selectLabel][paraName] = true
				t.ErrorPrint("atoi fail,the " + paraName + " value isn`t int,real value is:" + name.Text())
				continue
			}
			t.Cache[selectLabel][paraName] += valueNum
			t.DebugPrint(fmt.Sprintf("%s value is %f", paraName, valueNum))
		}

		for _, name := range object.SelectElements("CV") {
			index := name.SelectAttrValue("i", "unknown")
			for i, nameSN := range name.SelectElements("SN") {
				paraName, ok := t.TransTab[index]
				if !ok {
					t.ErrorPrint("userLabel:" + selectLabel + "para:" + index + "don't include the PmName,use the numKey")
					paraName = index
				}
				valueNum, err := strconv.ParseFloat(name.SelectElements("SV")[i].Text(), 64)
				if err != nil {
					if !t.initNilParaList(selectLabel, paraName+"-"+nameSN.Text()) {
						t.ErrorPrint("make map nilParaList fail the selectLabel is " + selectLabel + "," +
							"paraName is " + paraName + "-" + nameSN.Text())
					}
					t.ErrorPrint("atoi fail,the " + paraName + "-" + nameSN.Text() + " value isn`t int," +
						"real value is:" + name.Text())
					continue
				}
				t.Cache[selectLabel][paraName+"-"+nameSN.Text()] += valueNum
				t.DebugPrint(fmt.Sprintf("%s-%s value is %f", paraName, nameSN.Text(), valueNum))
			}
		}

	}

	// 缓存储存完毕，刷新成功文件数和当前处理周期
	t.ResultData.OkFileNum++
	t.ResultData.OkContainFiles[t.fileName] = true
	t.DebugPrint("OkContainFiles include " + t.fileName)
	return true
}

func (t *ToolData) initNilParaList(selectLabel string, paraName string) bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	if t.ResultData.NilParaName[selectLabel] == nil {
		t.ResultData.NilParaName[selectLabel] = make(map[string]bool)
	}
	if _, ok := t.ResultData.NilParaName[selectLabel][paraName]; !ok {
		t.ResultData.NilParaNum++
		t.ResultData.NilParaName[selectLabel][paraName] = true
	}
	return true
}

func (t *ToolData) printFailParaList(failType int) bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	var paraList map[string]map[string]bool
	var typeStr string
	if failType == paraZeroType {
		paraList = t.ResultData.ZeroParaName
		typeStr = "zero para"
	} else {
		paraList = t.ResultData.NilParaName
		typeStr = "nil para"
	}

	if len(paraList) != 0 {
		t.InfoPrintResult("****************************************************************************************************")
		t.InfoPrintResult("the " + typeStr + " list:")
		for label, mapData := range paraList {
			for paraName, _ := range mapData {
				t.InfoPrintResult(fmt.Sprintf("%-50s%s", label, paraName))
			}
		}
	}
	return true
}

func (t *ToolData) printAllResult() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	t.InfoPrintResult("***********************************************RESULT***********************************************")
	t.InfoPrintResult(fmt.Sprintf("check %d files,%d files check fail",
		t.ResultData.TotalFileNum, t.ResultData.TotalFileNum-t.ResultData.OkFileNum))
	if len(t.ResultData.ZeroParaName) == 0 && len(t.ResultData.NilParaName) == 0 {
		t.InfoPrintResult(fmt.Sprintf("all para is ok,total para num is %d", t.ResultData.TotalParaNum))
	} else {
		t.InfoPrintResult(fmt.Sprintf("total para is %d, zero para is %d,nil para is %d",
			t.ResultData.TotalParaNum, t.ResultData.ZeroParaNum, t.ResultData.NilParaNum))
	}
	t.InfoPrintResult("****************************************************************************************************")
	t.InfoPrintResult("contain period list:")
	for label, dataMap := range t.ResultData.containPeriods {
		for periodStr, _ := range dataMap {
			t.InfoPrintResult(label + "   " + periodStr)
		}
	}
	t.InfoPrintResult("****************************************************************************************************")
	t.InfoPrintResult(fmt.Sprintf("need check files list:"))
	var failFileNames []string
	for _, v := range t.ResultData.ContainFiles {
		t.InfoPrintResult(v)
		if _, ok := t.ResultData.OkContainFiles[v]; !ok {
			t.DebugPrint("failFileNames contain " + v)
			failFileNames = append(failFileNames, v)
		}
	}
	if !t.printFailParaList(paraZeroType) {
		t.ErrorPrint("print ZeroPara list fail")
	}
	if !t.printFailParaList(paraNilType) {
		t.ErrorPrint("print NilPara list fail")
	}
	t.InfoPrintResult("****************************************************************************************************")

	return true
}

func (t *ToolData) setResult() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	for label, data := range t.Cache {
		for paraName, valueNum := range data {
			t.ResultData.TotalParaNum++
			if valueNum == 0 {
				if t.ResultData.ZeroParaName[label] == nil {
					t.ResultData.ZeroParaName[label] = make(map[string]bool)
				}
				if _, ok := t.ResultData.ZeroParaName[label][paraName]; !ok {
					t.ResultData.ZeroParaNum++
					t.ResultData.ZeroParaName[label][paraName] = true
				}
			}
		}
	}

	if !t.printAllResult() {
		return false
	}

	return true
}
