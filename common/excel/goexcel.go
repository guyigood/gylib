package excel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/guyigood/gylib/common"
	"strconv"
)

func getcellcolsname(cols_row, col1 int) string {
	cols_sub := cols_row - 65 - 26
	result := string(rune(cols_row)) + strconv.Itoa(col1)
	if cols_sub >= 0 {
		result = "A" + string(rune(cols_sub+65)) + strconv.Itoa(col1)
	}
	return result
}

func Buildexcel(savename string, fdtitle, fdname []string, data []map[string]string) string {
	filename := common.Get_Upload_filename(savename, "")
	xlsx := excelize.NewFile()
	// Create a new sheet.
	xlsx.NewSheet("Sheet1")
	// Set value of a cell.
	j := 1
	cols_row := 65
	for k := 0; k < len(fdtitle); k++ {
		//fmt.Println(this.getcellcolsname(cols_row,1))
		xlsx.SetCellValue("Sheet1", getcellcolsname(cols_row, 1), fdtitle[k]) //string(rune(cols_row))+strconv.Itoa(j), fdtitle[k])
		cols_row++
	}
	j = 2
	for _, val := range data {
		cols_row = 65
		for i := 0; i < len(fdtitle); i++ {
			e_val, ok := val[fdname[i]]
			if ok {
				if val[fdname[i]] != "" {
					xlsx.SetCellValue("Sheet1", getcellcolsname(cols_row, j), e_val) //string(rune(cols_row))+strconv.Itoa(j), e_val)
				}
			}
			cols_row++
		}
		j++
	}
	xlsx.SetActiveSheet(2)
	err := xlsx.SaveAs(filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return filename[1:]
}

func Export_Excel(filename string, title_data []string, data []map[string]interface{}, title []string) bool {
	xlsx := excelize.NewFile()
	// Create a new sheet.
	xlsx.NewSheet("Sheet1")
	// Set value of a cell.
	cols_row := 65
	for _, val := range title_data {
		xlsx.SetCellValue("Sheet1", string(rune(cols_row))+"1", val)
		cols_row++
	}
	i := 2
	//var title =[]string{"bh","zh_name","duty_no","address","bank_no","skr","shr","kpr","memo","spname","ggxh","jldw","quantity","price","ws_price","sl","ssbm","is_kp"}
	for _, val_data := range data {
		cols_row = 65
		for j := 0; j < len(title); j++ {
			//fmt.Println(title[j],val_data[title[j]])
			xlsx.SetCellValue("Sheet1", string(rune(cols_row))+strconv.Itoa(i), val_data[title[j]])
			cols_row++
		}
		i++
	}
	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(2)
	// Save xlsx file by the given path.
	err := xlsx.SaveAs(filename)
	if err != nil {
		fmt.Println(err)
		return false

	}
	return true
}

func Import_Excel(filename string) {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	cell := xlsx.GetCellValue("Sheet1", "B2")
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows := xlsx.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
