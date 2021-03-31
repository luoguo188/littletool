package checktool

import (
	etree "modchecktool/github.com/etree-master"
	"strconv"
)

//定义
const (
	readConfigOk    = 1
	readConfigNotOk = 0
	readConfigSkip  = -1
)

const (
	paraNilType  = 0
	paraZeroType = 1
)

const (
	invalidMinuteTime = -1
	maxMinuteTime     = 60 * 24
)

// ToolData 工具数据，包括整个工程的数据和单个文件数据
type ToolData struct {
	*PmFileData
	PmFile *etree.Element

	ConfigData Config

	Cache      map[string]map[string]float64
	ResultData Result
}

//根据读到的字符决定目前使用哪些标签作为索引，1-网元ID（UserLabel），2-模块类型（ObjectType）, 3-网元形态（ElementType）,4-周期（Period）
type Config struct {
	UserLabelFlag   bool
	ObjectTypeFlag  bool
	ElementTypeFlag bool
	PeriodFlag      bool

	selectTimeBegin int
	selectTimeEnd   int

	ElementTypes map[string]bool

	DebugMode bool
}

type Result struct {
	TotalFileNum   int
	OkFileNum      int
	OkContainFiles map[string]bool
	ContainFiles   []string

	TotalParaNum int

	ZeroParaName map[string]map[string]bool
	ZeroParaNum  int

	NilParaName map[string]map[string]bool
	NilParaNum  int

	logName string

	containPeriods map[string]map[string]bool
}

// Measure 测量对象数据节点信息存储
type Measure struct {
	ObjectType *etree.Element
	PmName     *etree.Element
	PmData     *etree.Element
}

// FileHeader 文件头数据信息存储
type FileHeader struct {
	ElementType *etree.Element
	StartTime   *etree.Element
	Period      *etree.Element
}

// PmFileData 单个文件数据存储
type PmFileData struct {
	fileName string

	Measure     *etree.Element
	MeasureData Measure

	FileHeaderData FileHeader
	FileHeader     *etree.Element

	TransTab map[string]string
}

func (t *ToolData) getStartTime() string {
	if t == nil {
		t.ErrorPrint("t is nil")
		return ""
	}
	if t.FileHeaderData.StartTime == nil {
		t.ErrorPrint("StartTime tag is nil")
		return ""
	}
	return t.FileHeaderData.StartTime.Text()
}

func (t *ToolData) getStartHourMinute() int {
	if t == nil {
		t.ErrorPrint("t is nil")
		return invalidMinuteTime
	}

	timeStr := t.getStartTime()
	var hour, minute int
	var err error
	hour, err = strconv.Atoi(timeStr[len(timeStr)-8 : len(timeStr)-6])
	if err != nil {
		t.DebugPrint("atoi failed, StartTime hour isn`t number,real value is " + timeStr)
		return invalidMinuteTime
	}
	if minute, err = strconv.Atoi(timeStr[len(timeStr)-5 : len(timeStr)-3]); err != nil {
		t.DebugPrint("atoi failed, StartTime minute isn`t number，real value is " + timeStr)
		return invalidMinuteTime
	}
	return hour*60 + minute
}

func (t *ToolData) getElementType() string {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return ""
	}
	if t.FileHeaderData.ElementType == nil {
		t.DebugPrint("ElementType is nil")
		return ""
	}

	return t.FileHeaderData.ElementType.Text()
}

func (t *ToolData) getPeriod() string {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return ""
	}
	if t.FileHeaderData.Period == nil {
		t.DebugPrint("ObjectType is nil")
		return ""
	}
	return t.FileHeaderData.Period.Text()
}

func (t *ToolData) getObjectType() string {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return ""
	}
	if t.MeasureData.ObjectType == nil {
		t.DebugPrint("ObjectType is nil")
		return ""
	}
	return t.MeasureData.ObjectType.Text()
}
