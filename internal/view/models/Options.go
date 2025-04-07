package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Options struct {
	View *tview.Flex
}

func (a *App) CreateOptionsMenu() Options {
	list := tview.NewList().
		AddItem("фильтр", "", 0, nil).
		AddItem("Назад", "", 0, nil)

	a.setStyles(list)

	a.setCustomBorder(list.Box, "options")

	list.SetBorderPadding(1, 1, 2, 2)

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
			AddItem(list, 0, 5, true).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false),
			0, 1, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false)

	flex.SetBackgroundColor(tcell.ColorDefault)

	form := a.newOptionForm()

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch index {
		case 0:
			a.View.SetRoot(form, true)
		case 1:
			a.View.SetRoot(a.MainMenu.View, true)
			list.SetCurrentItem(0)
		}
	})

	return Options{
		View: flex,
	}
}

func (a *App) newOptionForm() *tview.Form {
	form := tview.NewForm()
	a.setCustomBorder(form.Box, "test")
	form.SetBorder(true).SetTitle("Edit Option1").SetTitleAlign(tview.AlignLeft)

	name, age, description := a.LoadConfig()

	form.AddInputField("Имя", name, 20, nil, nil).
		AddInputField("Возраст (<,>,=)", age, 20, nil, nil).
		AddInputField("Описание", description, 20, nil, nil).
		AddButton("Save", func() {

			name := form.GetFormItem(0).(*tview.InputField).GetText()
			age := form.GetFormItem(1).(*tview.InputField).GetText()
			description := form.GetFormItem(2).(*tview.InputField).GetText()
			saveConfig(name, age, description)

			a.View.SetRoot(a.Options.View, true)
		}).
		AddButton("Cancel", func() {
			a.View.SetRoot(a.Options.View, true)
		})
	return form
}

type Config struct {
	Name        string `json:"name"`
	AgeFilter   string `json:"age_filter"`
	Description string `json:"description"`
}

func (a *App) LoadConfig() (string, string, string) {
	config, err := loadConfigFromFile("config.json")
	if err != nil {
		return "", "", ""
	}

	return config.Name, config.AgeFilter, config.Description
}

func saveConfig(name string, age string, description string) {
	config := Config{
		Name:        name,
		AgeFilter:   age,
		Description: description,
	}

	err := saveConfigToFile(config, "config.json")
	if err != nil {
		fmt.Println("Error saving config:", err)
	}
}

func loadConfigFromFile(filename string) (*Config, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &config, nil
}

func saveConfigToFile(config Config, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling config to JSON: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}
