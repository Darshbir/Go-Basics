package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
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
		if strings.EqualFold(i, item) {
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

func Check_Valid_Day(day string, days []string) (bool, error) {
	valid_Day := false
	for _, day_name := range days {
		if strings.EqualFold(day_name, day) {
			valid_Day = true
		}
	}
	var err error
	if !valid_Day {
		err = errors.New("Not A Valid Day")
	}
	return valid_Day, err
}

func Check_Valid_Meal(meal string, meals []string) (bool, error) {
	valid_Meal := false
	for _, meal_name := range meals {
		if strings.EqualFold(meal_name, meal) {
			valid_Meal = true
		}
	}
	var err error
	if !valid_Meal {
		err = errors.New("Not A Valid Meal")
	}
	return valid_Meal, err
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
	fmt.Print("Choose an action to perform:\n Type 1 to get the items in a particular meal\n Type 2 for number of items in a particular meal\n Type 3 to check for presence of an item in the meal \n Type 4 to generate corresponding JSON files with all the items\n")
	var input string
	fmt.Scanln(&input)
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Error: Please input an integer between 1 and 4")
		return
	}
	if num < 1 || num > 4 {
		fmt.Println("Error: Please enter a number between 1 and 4.")
		return
	}
	if num == 1 {
		var day, meal_name string
		fmt.Print("Please Enter Day: ")
		fmt.Scanln(&day)
		fmt.Print("Please Enter Name of The meal: ")
		fmt.Scanln(&meal_name)
		valid_Meal, err := Check_Valid_Meal(meal_name, meal_Opt[:])
		if !valid_Meal {
			fmt.Print(err)
			return
		}
		valid_Day, err := Check_Valid_Day(day, days[:])
		if !valid_Day {
			fmt.Print(err)
			return
		}
		items := getMealItems(day, meal_name, xlsxFile)
		fmt.Println("Items:")
		for _, item := range items {
			fmt.Printf("%s\n", item)
		}
	} else if num == 2 {
		var day, meal_name string
		fmt.Print("Please Enter Day: ")
		fmt.Scanln(&day)
		fmt.Print("Please Enter Name of The meal: ")
		fmt.Scanln(&meal_name)
		valid_Meal, err := Check_Valid_Meal(meal_name, meal_Opt[:])
		if !valid_Meal {
			fmt.Print(err)
			return
		}
		valid_Day, err := Check_Valid_Day(day, days[:])
		if !valid_Day {
			fmt.Print(err)
			return
		}
		fmt.Printf("%v", getMealItemsCount(day, meal_name, xlsxFile))
	} else if num == 3 {
		var day, meal_name string
		fmt.Print("Please Enter Day: ")
		fmt.Scanln(&day)
		fmt.Print("Please Enter Name of The meal: ")
		fmt.Scanln(&meal_name)
		fmt.Print("Please Enter Name of The item you want to check: ")
		reader := bufio.NewReader(os.Stdin)
		item_name, _ := reader.ReadString('\n')
		item_name = strings.TrimSpace(item_name)
		valid_Meal, err := Check_Valid_Meal(meal_name, meal_Opt[:])
		if !valid_Meal {
			fmt.Print(err)
			return
		}
		valid_Day, err := Check_Valid_Day(day, days[:])
		if !valid_Day {
			fmt.Print(err)
			return
		}
		if isItemInMeal(day, meal_name, item_name, xlsxFile) {
			fmt.Print("This items is present in the meal")
		} else {
			fmt.Print("This item is not present in the meal")
		}
	} else if num == 4 {
		err = saveMenuAsJSON(menu)
		if err != nil {
			fmt.Println(err)
		}
	}
}
