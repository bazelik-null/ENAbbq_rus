package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Установщик русификатора для ENA: Dream BBQ")
	w.Resize(fyne.NewSize(800, 600))
	a.Settings().SetTheme(theme.DarkTheme())

	mainLabel := widget.NewLabel("Пожалуйста подождите, идёт загрузка ресурсов установщика...")
	init := container.New(layout.NewCenterLayout(), mainLabel)
	w.SetContent(init)
	w.Show()

	installFunc := func() {
		err := install()
		if err != 0 {
			w.SetContent(pageERR(w, err))
		} else {
			w.SetContent(page0(w))
		}
	}

	go installFunc()
	fyne.CurrentApp().Run()
}

func page0(w fyne.Window) *fyne.Container {
	mainLabel := widget.NewLabel("Добро пожаловать в установщик русификатора для ENA: Dream BBQ")

	teamLabel := canvas.NewText("  by BARBEQUE TEAM", color.RGBA{169, 169, 169, 255})
	teamLabel.TextSize = 12

	errorLabel := canvas.NewText("", color.RGBA{255, 0, 0, 255})

	btnContinue := widget.NewButton("Продолжить", func() {
		w.SetContent(pageInstall(w))
	})

	page0 := container.New(layout.NewVBoxLayout(),
		mainLabel,
		teamLabel,
		btnContinue,
		errorLabel,
	)

	page0 = container.New(layout.NewCenterLayout(), page0)

	checkIntegrity(btnContinue, errorLabel)

	return page0
}

func pageERR(_ fyne.Window, err int) *fyne.Container {
	pageERRContainer := container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			canvas.NewText("[FATL]: Произошла критическая ошибка при загрузке файлов.", color.RGBA{255, 0, 0, 255}),
			canvas.NewText("[FATL]: Error "+fmt.Sprint(err), color.RGBA{255, 0, 0, 255}),
			widget.NewButtonWithIcon("Закрыть", theme.WindowCloseIcon(), func() {
				fyne.CurrentApp().Quit()
			}),
		),
	)
	return pageERRContainer
}

func pageInstall(w fyne.Window) *fyne.Container {
	var path string

	appDir, _ := os.Getwd()
	steamIcon, _ := fyne.LoadResourceFromPath(filepath.Join(appDir, "resources", "steamIcon.png"))

	labelPath := widget.NewLabel("")
	errorLabel := canvas.NewText("", color.RGBA{255, 0, 0, 255})

	btnContinue := widget.NewButtonWithIcon("Установить", theme.DownloadIcon(), func() {
		w.SetContent(pageEnd(path))
	})
	btnContinue.Disable()

	buttonsContainer := container.New(layout.NewVBoxLayout())

	btnSteam := widget.NewButtonWithIcon("Steam", steamIcon, func() {
		if runtime.GOOS == "windows" {
			path = filepath.Join("C:\\", "Program Files (x86)", "Steam", "steamapps", "common", "ENA Dream BBQ")
		} else {
			homeDir := os.Getenv("HOME")
			path = filepath.Join(homeDir, ".steam", "root", "steamapps", "common", "ENA Dream BBQ")
		}
		checkExecutable(path, btnContinue, errorLabel)
		labelPath.SetText("Выбранный путь: " + path)
	})

	btnBrowse := widget.NewButtonWithIcon("Открыть", theme.SearchIcon(), func() {
		browseFile(w, func(selectedPath string) {
			path = selectedPath
			checkExecutable(path, btnContinue, errorLabel)
			labelPath.SetText("Выбранный путь: " + path)
		})
	})

	buttonsContainer.Add(btnSteam)
	buttonsContainer.Add(btnBrowse)

	pageInstall := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("Выберите путь до игры. Если она установлена по стандартному пути нажмите на кнопку Steam."),
		buttonsContainer,
		labelPath,
		errorLabel,
		btnContinue,
	)

	pageInstall = container.New(layout.NewCenterLayout(), pageInstall)

	return pageInstall
}

func checkIntegrity(btnContinue *widget.Button, errorLabel *canvas.Text) {
	appDir, _ := os.Getwd()
	resourcesPath := filepath.Join(appDir, "resources", "meta.json")
	if _, err := os.Stat(resourcesPath); os.IsNotExist(err) {
		btnContinue.Disable()
		errorLabel.Text = "[FATL]: \"resources\" не найдено."
		errorLabel.Refresh()
	} else {
		btnContinue.Enable()
		errorLabel.Text = ""
		errorLabel.Refresh()
	}
}

func checkExecutable(selectedPath string, btnContinue *widget.Button, errorLabel *canvas.Text) {
	executablePath := filepath.Join(selectedPath, "ENA-4-DreamBBQ.exe")
	if _, err := os.Stat(executablePath); os.IsNotExist(err) {
		btnContinue.Disable()
		errorLabel.Text = "[ERROR]: \"ENA-4-DreamBBQ.exe\" не найден, выберите папку с исполняемым файлом игры"
		errorLabel.Refresh()
	} else {
		btnContinue.Enable()
		errorLabel.Text = ""
		errorLabel.Refresh()
	}
}

func browseFile(w fyne.Window, onPathSelected func(string)) {
	dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if folder != nil {
			onPathSelected(folder.Path())
		}
	}, w)
}

func pageEnd(path string) *fyne.Container {
	appDir, _ := os.Getwd()
	err := inject(path)
	os.RemoveAll(filepath.Join(appDir, "resources"))
	if err != nil {
		pageEndContainer := container.New(layout.NewCenterLayout(),
			container.New(layout.NewVBoxLayout(),
				canvas.NewText("[FATL]: Произошла критическая ошибка при инъекции ассетов.", color.RGBA{255, 0, 0, 255}),
				canvas.NewText("[FATL]: Error "+fmt.Sprint(err), color.RGBA{255, 0, 0, 255}),
				widget.NewButtonWithIcon("Закрыть", theme.WindowCloseIcon(), func() {
					fyne.CurrentApp().Quit()
				}),
			),
		)
		return pageEndContainer
	} else {
		pageEndContainer := container.New(layout.NewCenterLayout(),
			container.New(layout.NewVBoxLayout(),
				widget.NewLabel("Спасибо за установку"),
				widget.NewButtonWithIcon("Закрыть", theme.WindowCloseIcon(), func() {
					fyne.CurrentApp().Quit()
				}),
			),
		)
		return pageEndContainer
	}
}
