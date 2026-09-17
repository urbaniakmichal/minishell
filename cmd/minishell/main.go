package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type terminalTheme struct{}

func (m *terminalTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground, theme.ColorNameInputBackground:
		return color.RGBA{R: 15, G: 15, B: 20, A: 255}
	case theme.ColorNameForeground, theme.ColorNameButton:
		return color.RGBA{R: 0, G: 255, B: 128, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 80, G: 100, B: 90, A: 255}
	case theme.ColorNameFocus:
		return color.RGBA{R: 0, G: 200, B: 100, A: 255}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m *terminalTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m *terminalTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *terminalTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func createTerminalTab() fyne.CanvasObject {
	historyView := widget.NewMultiLineEntry()
	historyView.SetPlaceHolder("Command history...")
	historyView.Wrapping = fyne.TextWrapWord
	historyView.Disable()

	outputView := widget.NewMultiLineEntry()
	outputView.SetPlaceHolder("Results...")
	outputView.Wrapping = fyne.TextWrapWord
	outputView.Disable()

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter command...")
	input.OnSubmitted = func(text string) {
		if text == "" {
			return
		}
		historyView.SetText(historyView.Text + fmt.Sprintf("%s > %s\n", time.Now().Format("15:04:05"), text))
		outputView.SetText(outputView.Text + fmt.Sprintf("Executing: %s\n[Mock Output] Success\n", text))
		input.SetText("")
	}

	splitView := container.NewHSplit(outputView, historyView)
	splitView.SetOffset(0.75)

	return container.NewBorder(
		nil,
		input,
		nil,
		nil,
		splitView,
	)
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&terminalTheme{})

	shell := a.NewWindow("Mini Shell")

	if icon, err := fyne.LoadResourceFromPath("terminal.png"); err == nil {
		shell.SetIcon(icon)
	}

	terminalCount := 1

	// Funkcja pomocnicza tworząca ponumerowaną zakładkę
	var createTabItem func() *container.TabItem
	createTabItem = func() *container.TabItem {
		name := fmt.Sprintf("Terminal %d", terminalCount)
		terminalCount++
		return container.NewTabItem(name, createTerminalTab())
	}

	// Specjalna zakładka "+" z pustą zawartością (będzie natychmiast zastępowana)
	plusTab := container.NewTabItem("+", widget.NewLabel(""))

	tabs := container.NewAppTabs(
		createTabItem(),
		plusTab,
	)

	// Nasłuchiwanie kliknięcia w zakładki
	tabs.OnSelected = func(selected *container.TabItem) {
		if selected == plusTab {
			// 1. Tworzymy nowy terminal
			newTab := createTabItem()

			// 2. Wstawiamy go tuż przed zakładkę "+"
			var newItems []*container.TabItem
			for _, item := range tabs.Items {
				if item == plusTab {
					newItems = append(newItems, newTab)
				}
				newItems = append(newItems, item)
			}
			tabs.Items = newItems
			tabs.Refresh()

			// 3. Automatycznie przełączamy na nowo utworzoną zakładkę
			tabs.Select(newTab)
		}
	}

	tabs.SetTabLocation(container.TabLocationTop)

	shell.SetContent(tabs)
	shell.Resize(fyne.NewSize(1000, 500))
	shell.CenterOnScreen()
	shell.ShowAndRun()
}
