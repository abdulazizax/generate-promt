package helper

import (
	"encoding/json"
	"fmt"
	"sort"

	"google.golang.org/api/sheets/v4"
)

func WriteDataToCompetitorsSheet(srv *sheets.Service, spreadsheetID string, sheetName string, resp string) error {
	var data []map[string]interface{}

	err := json.Unmarshal([]byte(resp), &data)
	if err != nil {
		fmt.Println("Ошибка при парсинге JSON:", err)
		return err
	}

	// Если данных нет, возвращаем ошибку или nil
	if len(data) == 0 {
		return nil
	}

	// Формируем динамический header — объединяем ключи из всех строк
	headerMap := make(map[string]bool)
	for _, row := range data {
		for key := range row {
			headerMap[key] = true
		}
	}

	// Собираем ключи в срез для упорядочивания
	var keys []string
	for key := range headerMap {
		keys = append(keys, key)
	}
	// Сортируем ключи для консистентного порядка (опционально, можно задать свой порядок)
	sort.Strings(keys)

	// Преобразуем отсортированные ключи в срез interface{}
	header := make([]interface{}, len(keys))
	for i, key := range keys {
		header[i] = key
	}

	// Инициализируем двумерный срез значений с header-строкой
	values := [][]interface{}{header}

	// Преобразуем каждый элемент data в срез значений в том же порядке, что и header
	for _, row := range data {
		rowData := make([]interface{}, len(keys))
		for i, key := range keys {
			if val, ok := row[key]; ok {
				rowData[i] = val
			} else {
				rowData[i] = ""
			}
		}
		values = append(values, rowData)
	}

	// Диапазон записи начинается с ячейки A1 указанного листа
	writeRange := fmt.Sprintf("%s!A1", sheetName)

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Обновляем значения на листе (используем опцию RAW)
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, valueRange).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		return fmt.Errorf("не удалось обновить данные: %v", err)
	}
	return nil
}
