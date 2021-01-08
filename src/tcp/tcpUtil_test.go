package tcp

import (
	"stress-tool/model"
	"testing"
)

var testValues [][]map[string]int
var testResult []bool

func addTestItem(dataTypeMap map[string]int, result bool) {
	arr := make([]map[string]int, len(dataTypeMap))
	index := 0
	for k, v := range dataTypeMap {
		m := make(map[string]int, 1)
		m[k] = v
		arr[index] = m
		index++
	}
	testValues = append(testValues, arr)
	testResult = append(testResult, result)
}

func TestCheckDataRangeIsContinuous(t *testing.T) {
	addTestItem(map[string]int{
		"0~3":   model.STRING,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, true)
	addTestItem(map[string]int{
		"1~3":   model.STRING,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"-1~3":  model.STRING,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"3~0":   model.STRING,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"3~0":   model.STRING,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~3":   model.STRING,
		"4~-12": model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~3":   model.STRING,
		"4~12":  model.STRING,
		"11~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~11":  model.STRING,
		"12~12": model.STRING,
		"13~20": model.STRING,
	}, true)
	addTestItem(map[string]int{
		"0~4":   model.NUMBER,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~4":   model.FLOAT,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~4":   model.FLOAT,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~a":   model.FLOAT,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0ss":   model.FLOAT,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)
	addTestItem(map[string]int{
		"0~2~3": model.FLOAT,
		"4~12":  model.STRING,
		"13~20": model.STRING,
	}, false)

	i := 0
	var msg string
	for length := len(testValues); i < length; i++ {
		m := testValues[i]
		_, msg = checkDataRangeIsContinuous(m)
		expected := ""
		if !testResult[i] {
			expected = "Not null string"
		}
		t.Logf("Test CheckDataRangeIsContinuous: param: [%v], expected: [%s], get: [%s]", m, expected, msg)
		if msg == "" && !testResult[i] || msg != "" && testResult[i] {
			t.Errorf("TestCheckDataRangeIsContinuous fail")
		}
	}
}
