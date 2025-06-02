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
	// Init app
	a := app.New()
	w := a.NewWindow("Установщик русификатора для ENA: Dream BBQ")
	w.Resize(fyne.NewSize(800, 600))

	mainLabel := widget.NewLabel("Пожалуйста подождите, идёт загрузка ресурсов установщика...")
	loadingWidget := widget.NewActivity()
	loadingWidget.Start()

	init := container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			mainLabel,
			loadingWidget,
		),
	)

	w.SetContent(init)
	w.Show()

	// Download files in goroutine
	go func() {
		err := download()
		if err != 0 {
			w.SetContent(pageERR(w, err))
		} else {
			w.SetContent(page0(w))
		}
	}()

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

	page0 := container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			mainLabel,
			teamLabel,
			btnContinue,
			errorLabel,
		),
	)

	// Check integrity of downloaded files
	checkIntegrity(btnContinue, errorLabel)

	return page0
}

func pageERR(_ fyne.Window, err int) *fyne.Container {
	errLabel := canvas.NewText("[FATL]: Произошла критическая ошибка при загрузке файлов.", color.RGBA{255, 0, 0, 255})
	errCode := canvas.NewText("[FATL]: Error "+fmt.Sprint(err), color.RGBA{255, 0, 0, 255})

	buttonClose := widget.NewButtonWithIcon("Закрыть", theme.WindowCloseIcon(), func() {
		fyne.CurrentApp().Quit()
	})

	pageERRContainer := container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			errLabel,
			errCode,
			buttonClose,
		),
	)
	return pageERRContainer
}

func pageInstall(w fyne.Window) *fyne.Container {
	var path string
	var steamIcon fyne.Resource

	appDir, _ := os.Getwd()
	currentTheme := fyne.CurrentApp().Settings().Theme()

	// Load an icon depending on theme
	if currentTheme == theme.DarkTheme() {
		steamIcon, _ = fyne.LoadResourceFromPath(filepath.Join(appDir, "resources", "steamIconW.png"))
	} else {
		steamIcon, _ = fyne.LoadResourceFromPath(filepath.Join(appDir, "resources", "steamIcon.png"))
	}

	label := widget.NewLabel("Выберите путь до игры. Если она установлена по стандартному пути нажмите на кнопку Steam.")
	labelPath := widget.NewLabel("")
	errorLabel := canvas.NewText("", color.RGBA{255, 0, 0, 255})

	btnContinue := widget.NewButtonWithIcon("Установить", theme.DownloadIcon(), func() {
		w.SetContent(pageEnd(path))
	})
	btnContinue.Disable()

	btnSteam := widget.NewButtonWithIcon("Steam", steamIcon, func() {
		// Choose default path to game depending on OS
		if runtime.GOOS == "windows" {
			path = filepath.Join("C:\\", "Program Files (x86)", "Steam", "steamapps", "common", "ENA Dream BBQ")
		} else {
			homeDir := os.Getenv("HOME")
			path = filepath.Join(homeDir, ".steam", "root", "steamapps", "common", "ENA Dream BBQ")
		}
		// Check if there is executable game file
		checkExecutable(path, btnContinue, errorLabel)
		// Display chosen path
		labelPath.SetText("Выбранный путь: " + path)
	})

	btnBrowse := widget.NewButtonWithIcon("Открыть", theme.SearchIcon(), func() {
		browseFile(w, func(selectedPath string) {
			path = selectedPath
			// Check if there is executable game file
			checkExecutable(path, btnContinue, errorLabel)
			// Display chosen path
			labelPath.SetText("Выбранный путь: " + path)
		})
	})

	pageInstall := container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			label,
			btnSteam,
			btnBrowse,
			labelPath,
			errorLabel,
			btnContinue,
		),
	)

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

	err := install(path)

	os.RemoveAll(filepath.Join(appDir, "resources"))

	if err != nil {
		errLabel := canvas.NewText("[FATL]: Произошла критическая ошибка при инъекции ассетов.", color.RGBA{255, 0, 0, 255})
		errCode := canvas.NewText("[FATL]: Error "+fmt.Sprint(err), color.RGBA{255, 0, 0, 255})

		buttonClose := widget.NewButtonWithIcon("Закрыть", theme.WindowCloseIcon(), func() {
			fyne.CurrentApp().Quit()
		})

		pageEndContainer := container.New(layout.NewCenterLayout(),
			container.New(layout.NewVBoxLayout(),
				errLabel,
				errCode,
				buttonClose,
			),
		)
		return pageEndContainer
	} else {
		label := widget.NewLabel("Спасибо за установку")
		buttonClose := widget.NewButtonWithIcon("Закрыть", theme.WindowCloseIcon(), func() {
			fyne.CurrentApp().Quit()
		})

		pageEndContainer := container.New(layout.NewCenterLayout(),
			container.New(layout.NewVBoxLayout(),
				label,
				buttonClose,
			),
		)
		return pageEndContainer
	}
}
