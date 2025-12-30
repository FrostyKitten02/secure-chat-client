package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"secure-chat-client/auth"
)

var users = map[string]string{} // in-memory username -> password

func Start() {
	a := app.New()
	w := a.NewWindow("Secure Chat")
	w.Resize(fyne.NewSize(400, 300))

	// Declare function variables so closures can reference each other
	var showLogin func()
	var showRegisterScreen func()
	var showMainScreen func(username string)

	// Login screen
	showLogin = func() {
		username := widget.NewEntry()
		username.SetPlaceHolder("Username or email")
		password := widget.NewPasswordEntry()
		password.SetPlaceHolder("Password")

		loginBtn := widget.NewButton("Login", func() {
			if pw, ok := users[username.Text]; ok && pw == password.Text {
				showMainScreen(username.Text)
			} else {
				w.SetContent(widget.NewLabel("Invalid username or password"))
			}
		})

		registerBtn := widget.NewButton("New user? Register", func() {
			showRegisterScreen()
		})

		content := container.NewVBox(
			widget.NewLabel("Login"),
			username,
			password,
			loginBtn,
			registerBtn,
		)
		w.SetContent(content)
	}

	// Register screen
	showRegisterScreen = func() {
		email := widget.NewEntry()
		email.SetPlaceHolder("Email")

		username := widget.NewEntry()
		username.SetPlaceHolder("Username")

		password := widget.NewPasswordEntry()
		password.SetPlaceHolder("Password")
		passwordConfirm := widget.NewPasswordEntry()
		passwordConfirm.SetPlaceHolder("Confirm Password")

		registerBtn := widget.NewButton("Register", func() {
			if username.Text == "" || password.Text == "" {
				w.SetContent(widget.NewLabel("Username and password required"))
				return
			}
			registerErr := auth.Register(username.Text, email.Text, password.Text)
			if registerErr != nil {
				w.SetContent(widget.NewLabel(registerErr.Error()))
				return
			}

			w.SetContent(widget.NewLabel("Registration successful! Go back to login."))
		})

		backBtn := widget.NewButton("Back to Login", showLogin)

		content := container.NewVBox(
			widget.NewLabel("Register"),
			email,
			username,
			password,
			passwordConfirm,
			registerBtn,
			backBtn,
		)
		w.SetContent(content)
	}

	// Main screen
	showMainScreen = func(username string) {
		content := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Welcome, %s!", username)),
			widget.NewButton("Logout", showLogin),
		)
		w.SetContent(content)
	}

	// start with login screen
	showLogin()
	w.ShowAndRun()
}
