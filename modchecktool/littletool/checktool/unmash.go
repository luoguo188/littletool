package checktool

import (
	"fmt"
	etree "modchecktool/github.com/etree-master"
)

func (t *ToolData) unMashRoot(fileName string) bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	t.fileName = fileName
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(fileName); err != nil {
		t.ErrorPrint("read file " + fileName + " error detail is" + err.Error())
		return false
	}
	t.PmFile = doc.SelectElement("PmFile")
	return true
}

func (t *ToolData) unMashElement(root *etree.Element, tag string) *etree.Element {
	if root == nil {
		t.DebugPrint(fmt.Sprintf("%s parent node does't find", tag))
		return nil
	}
	return root.SelectElement(tag)
}

func (t *ToolData) unMashTransTab() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}

	if t.MeasureData.PmName == nil {
		fmt.Println("Measure tag isn`t find")
		return false
	}

	for _, name := range t.MeasureData.PmName.SelectElements("N") {
		if name == nil {
			t.DebugPrint("PmName tag isn`t exist")
			continue
		}
		index := name.SelectAttrValue("i", "unknown")
		t.TransTab[index] = name.Text()
	}
	return true

}

func (t *ToolData) unMashPmFileHeaderData() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	t.FileHeader = t.unMashElement(t.PmFile, "FileHeader")

	t.FileHeaderData.ElementType = t.unMashElement(t.FileHeader, "ElementType")
	t.FileHeaderData.StartTime = t.unMashElement(t.FileHeader, "StartTime")
	t.FileHeaderData.Period = t.unMashElement(t.FileHeader, "Period")

	return true
}

func (t *ToolData) unMashPmFileData() bool {
	t.DebugPrint("enter")
	if t == nil {
		t.ErrorPrint("t is nil")
		return false
	}
	t.Measure = t.unMashElement(t.PmFile, "Measurements")

	if t.Measure == nil {
		t.ErrorPrint(t.fileName + ":Measure tag  isn`t exist")
		return false
	}

	t.MeasureData.PmName = t.unMashElement(t.Measure, "PmName")
	t.MeasureData.ObjectType = t.unMashElement(t.Measure, "ObjectType")
	t.MeasureData.PmData = t.unMashElement(t.Measure, "PmData")

	if t.MeasureData.PmData == nil {
		t.ErrorPrint(t.fileName + ":PmData tag  isn`t exist")
		return false
	}

	if !t.unMashTransTab() {
		t.ErrorPrint(t.fileName + "trans tabs get fail(PmName format is wrong)")
	}
	return true
}
