package helper

import (
	"encoding/json"
	"fmt"
	"generate-promt-v1/api/models"
	"strings"

	"google.golang.org/api/sheets/v4"
)

func WriteDataToPricingSheet(srv *sheets.Service, spreadsheetID string, sheetName string, resp string) error {
	var estimate models.ProjectEstimate
	if err := json.NewDecoder(strings.NewReader(resp)).Decode(&estimate); err != nil {
		return err
	}

	// Итоговый двумерный массив для записи в Google Sheets
	values := [][]interface{}{}

	//----------------------------------------------------------------------
	// 1. Таблица с командой
	//----------------------------------------------------------------------
	values = append(values, []interface{}{"Сотрудник", "Количество", "Месяц", "Ставка", "Сумма"})
	for _, member := range estimate.Team {
		row := []interface{}{
			member.Role,
			member.Count,
			member.Months,
			member.MonthlySalary,
			member.Sum,
		}
		values = append(values, row)
	}

	// Пустая строка
	values = append(values, []interface{}{})

	//----------------------------------------------------------------------
	// 2. Таблица с модулями
	//----------------------------------------------------------------------
	// Если вам нужно выводить инфо о модулях (Backend, Frontend и т.д.).
	if len(estimate.Modules) > 0 {
		// Шапка для таблицы модулей
		values = append(values, []interface{}{"Модуль", "Часы", "Ставка/час", "Стоимость"})
		for _, mod := range estimate.Modules {
			row := []interface{}{
				mod.ModuleName,
				mod.Hours,
				mod.HourlyRate,
				mod.Cost,
			}
			values = append(values, row)
		}
		// Добавляем пустую строку после таблицы
		values = append(values, []interface{}{})
	}

	//----------------------------------------------------------------------
	// 3. Финансовый план
	//----------------------------------------------------------------------
	fin := estimate.FinancialPlan

	// Строка-заголовок для предоплаты и месяцев
	headerRow := []interface{}{fmt.Sprintf("Предоплата (%.0f%%)", fin.PrepaymentPercent)}
	for i := range fin.MonthlyPayments {
		headerRow = append(headerRow, fmt.Sprintf("%d месяц", i+1))
	}
	values = append(values, headerRow)

	// Строка со значениями предоплаты и месячных платежей
	valueRow := []interface{}{fin.Prepayment}
	for _, payment := range fin.MonthlyPayments {
		valueRow = append(valueRow, payment)
	}
	values = append(values, valueRow)

	// Пустая строка
	values = append(values, []interface{}{})

	// Итого
	values = append(values, []interface{}{"Сумма проекта:", fin.TotalProjectCost})

	//----------------------------------------------------------------------
	// 4. Подготовка ValueRange и запись в Sheets
	//----------------------------------------------------------------------
	valueRange := &sheets.ValueRange{
		Values: values,
	}
	writeRange := fmt.Sprintf("%s!A1", sheetName)

	_, err := srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, valueRange).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		return fmt.Errorf("Sheets update error: %w", err)
	}

	return nil
}
