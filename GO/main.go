package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Meal struct {
	Day   string   `json:"day"`
	Date  string   `json:"date"`
	Meal  string   `json:"meal"`
	Items []string `json:"items"`
}

func NewMeal(day, date, meal string, items []string) *Meal {
	return &Meal{
		Day:   day,
		Date:  date,
		Meal:  meal,
		Items: items,
	}
}

func (m *Meal) PrintDetails() {
	fmt.Printf("Day: %s, Date: %s, Meal: %s\n", m.Day, m.Date, m.Meal)
	fmt.Println("Items:")
	for _, item := range m.Items {
		fmt.Printf("- %s\n", item)
	}
}

func getMealItems(day, meal string, sheet *excelize.File) []string {
	cols, err := sheet.GetCols("Sheet1")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, col := range cols {
		if len(col) > 0 && strings.EqualFold(col[0], day) {
			i := 1
			for i < len(col) && !strings.EqualFold(col[i], meal) {
				i++
			}

			i++
			var mealItems []string
			for i < len(col) && col[i] != "" && !strings.EqualFold(col[i], day) {
				fmt.Printf("%s\n", col[i])
				mealItems = append(mealItems, col[i])
				i++
			}

			return mealItems
		}
	}
	return nil
}

func getMealItemsCount(day, meal string, sheet *excelize.File) int {
	items := getMealItems(day, meal, sheet)
	return len(items)
}

func isItemInMeal(day, meal, item string, sheet *excelize.File) bool {
	items := getMealItems(day, meal, sheet)
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func saveMenuAsJSON(menu []*Meal) error {
	jsonData, err := json.MarshalIndent(menu, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("menu.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	xlsxFile, err := excelize.OpenFile("data/Sample-Menu.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	days := [7]string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"}
	meal_Opt := [3]string{"BREAKFAST", "LUNCH", "DINNER"}
	var menu []*Meal

	for i := 0; i < 7; i++ {
		var date string
		if i <= 4 {
			date = fmt.Sprintf("0%v-Feb-24", i+5)
		} else {
			date = fmt.Sprintf("%v-Feb-24", i+5)
		}
		for j := 0; j < 3; j++ {
			menu = append(menu, NewMeal(days[i], date, meal_Opt[j], getMealItems(days[i], meal_Opt[j], xlsxFile)))
		}
	}

	for _, meal := range menu {
		meal.PrintDetails()
	}

	err = saveMenuAsJSON(menu)
	if err != nil {
		fmt.Println(err)
	}
}
