package main

import (
	"fmt"
	"image/color"
	"os/exec"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	shell := startAppAndWindow()
	terminalCount := 1

	var createTabItem func() *container.TabItem
	createTabItem = func() *container.TabItem {
		name := fmt.Sprintf("Terminal %d", terminalCount)
		terminalCount++
		return container.NewTabItem(name, createTerminalTab())
	}

	plusTab := container.NewTabItem("+", widget.NewLabel(""))
	tabs := container.NewAppTabs(
		createTabItem(),
		plusTab,
	)
	tabs.OnSelected = func(selected *container.TabItem) {
		if selected == plusTab {
			newTab := createTabItem()

			var newItems []*container.TabItem
			for _, item := range tabs.Items {
				if item == plusTab {
					newItems = append(newItems, newTab)
				}
				newItems = append(newItems, item)
			}
			tabs.Items = newItems
			tabs.Refresh()

			tabs.Select(newTab)
		}
	}

	tabs.SetTabLocation(container.TabLocationTop)
	showShell(shell, tabs)
}

func startAppAndWindow() fyne.Window {
	a := app.New()
	a.Settings().SetTheme(&terminalTheme{})
	shell := a.NewWindow("Mini Shell")
	if icon, err := fyne.LoadResourceFromPath("terminal.png"); err == nil {
		shell.SetIcon(icon)
	}
	return shell
}

func showShell(shell fyne.Window, tabs *container.AppTabs) {
	shell.SetContent(tabs)
	shell.Resize(fyne.NewSize(1000, 500))
	shell.CenterOnScreen()
	shell.ShowAndRun()
}

func runCommand(command string) ([]byte, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", command)
	case "darwin", "linux":
		cmd = exec.Command("sh", "-c", command)
	default:
		cmd = exec.Command("cmd", "/c", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}

	return output, nil
}

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

	resultsBox := container.NewVBox()
	scrollableResults := container.NewScroll(resultsBox)

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter command...")
	input.OnSubmitted = func(text string) {
		if text == "" {
			return
		}

		addBlock := func(text string, bgColor color.RGBA) {
			content := widget.NewLabel(text)
			content.Wrapping = fyne.TextWrapWord
			copyBtn := widget.NewButton("Copy", func() {
				fyne.CurrentApp().Clipboard().SetContent(text)
			})
			copyBtn.Importance = widget.HighImportance
			copyBtn.Resize(fyne.NewSize(60, 30))
			rowContent := container.NewBorder(nil, nil, nil, copyBtn, content)
			bg := canvas.NewRectangle(bgColor)
			resultsBox.Add(container.NewStack(bg, rowContent))
			resultsBox.Refresh()
			scrollableResults.ScrollToBottom()
		}

		output, err := runCommand(text)
		if err != nil {
			addBlock(err.Error(), color.RGBA{R: 90, G: 25, B: 25, A: 255})
		} else {
			var bgCol color.RGBA
			if len(resultsBox.Objects)%2 == 0 {
				bgCol = color.RGBA{R: 20, G: 20, B: 30, A: 255}
			} else {
				bgCol = color.RGBA{R: 30, G: 35, B: 50, A: 255}
			}
			addBlock(string(output), bgCol)
		}

		historyView.SetText(historyView.Text + fmt.Sprintf("%s > %s\n", time.Now().Format("15:04:05"), text))
		input.SetText("")
	}

	splitView := container.NewHSplit(scrollableResults, historyView)
	splitView.SetOffset(0.75)

	return container.NewBorder(
		nil,
		input,
		nil,
		nil,
		splitView,
	)
}
