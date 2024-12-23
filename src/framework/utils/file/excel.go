package file

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"

	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

// ReadSheet 表格读取数据
//
// filePath 文件路径地址
//
// sheetName 工作簿名称， 空字符默认Sheet1
func ReadSheet(filePath, sheetName string) ([]map[string]string, error) {
	data := make([]map[string]string, 0)
	// 打开 Excel 文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return data, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Errorf("工作表文件关闭失败 : %v", err)
		}
	}()

	// 检查工作簿是否存在
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	if visible, _ := f.GetSheetVisible(sheetName); !visible {
		return data, fmt.Errorf("读取工作簿 %s 失败", sheetName)
	}

	// 获取工作簿记录
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return data, err
	}

	for i, row := range rows {
		// 跳过第一行
		if i == 0 {
			continue
		}
		// 遍历每个单元格
		rowData := map[string]string{}
		for columnIndex, cellValue := range row {
			columnName, _ := excelize.ColumnNumberToName(columnIndex + 1)
			rowData[columnName] = cellValue
		}

		data = append(data, rowData)
	}
	return data, nil
}

// WriteSheet 表格写入数据
//
// headerCells 第一行表头标题 "A1":"?"
//
// dataCells 从第二行开始的数据 "A2":"?"
//
// fileName 文件名称
//
// sheetName 工作簿名称， 空字符默认Sheet1
func WriteSheet(headerCells map[string]string, dataCells []map[string]any, fileName, sheetName string) (string, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger.Errorf("工作表文件关闭失败 : %v", err)
		}
	}()

	// 创建一个工作表
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("创建工作表失败 %v", err)
	}
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)

	// 首个和最后key名称
	firstKey := "A"
	lastKey := "B"

	// 第一行表头标题
	for key, title := range headerCells {
		_ = f.SetCellValue(sheetName, key, title)
		if key[:1] > lastKey {
			lastKey = key[:1]
		}
	}

	// 设置工作表上宽度为 20
	_ = f.SetColWidth(sheetName, firstKey, lastKey, 20)

	// 从第二行开始的数据
	for _, cell := range dataCells {
		for key, value := range cell {
			_ = f.SetCellValue(sheetName, key, value)
		}
	}

	// 上传资源路径
	_, dir := resourceUpload()
	filePath := filepath.Join(constants.UPLOAD_EXPORT, date.ParseDatePath(time.Now()))
	saveFilePath := filepath.Join(dir, filePath, fileName)

	// 创建文件目录
	if err := os.MkdirAll(filepath.Dir(saveFilePath), 0755); err != nil {
		return "", fmt.Errorf("创建保存文件失败 %v", err)
	}

	// 根据指定路径保存文件
	if err := f.SaveAs(saveFilePath); err != nil {
		return "", fmt.Errorf("保存工作表失败 %v", err)
	}
	return saveFilePath, nil
}
