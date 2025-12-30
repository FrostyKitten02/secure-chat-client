package gui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"secure-chat-client/auth"
	"secure-chat-client/client"
	"secure-chat-client/service"
	"secure-chat-client/state"
)

func Start() {
	a := app.New()
	w := a.NewWindow("Secure Chat")
	w.Resize(fyne.NewSize(400, 300))

	// Declare function variables so closures can reference each other
	var showLogin func()
	var showRegisterScreen func()
	var showMainScreen func()
	var showAddUserScreen func()

	// Login screen
	showLogin = func() {
		username := widget.NewEntry()
		username.SetPlaceHolder("Username or email")
		password := widget.NewPasswordEntry()
		password.SetPlaceHolder("Password")

		loginBtn := widget.NewButton("Login", func() {
			loginErr := auth.Login(username.Text, password.Text)
			if loginErr != nil {
				dialog.ShowInformation("Login", "Login failed", w)
				return
			}
			showMainScreen()
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
				dialog.ShowError(errors.New("username or password required"), w)
				return
			}
			registerErr := auth.Register(username.Text, email.Text, password.Text)
			if registerErr != nil {
				dialog.ShowError(registerErr, w)
				return
			}

			dialog.ShowInformation("Register", "Registration successful", w)
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
	showMainScreen = func() {
		// ---- Left: user list ----
		userList := widget.NewList(
			func() int {
				return len(state.Chats)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				chat := state.Chats[i]
				o.(*widget.Label).SetText(chat.User.Username)
			},
		)

		// ---- Right: chat area ----
		chatLabel := widget.NewLabel("Select a user to start chatting")
		chatLabel.Wrapping = fyne.TextWrapWord

		chatScroll := container.NewVScroll(chatLabel)

		messageEntry := widget.NewEntry()
		messageEntry.SetPlaceHolder("Type a message...")

		sendButton := widget.NewButton("Send", func() {
			//TODO: implement!!!
			if messageEntry.Text == "" {
				return
			}
			err := service.SendMessage(state.CurrentChatUserId, "test message")
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			messageEntry.SetText("")
		})

		chatRight := container.NewBorder(
			nil,
			container.NewBorder(
				nil,
				nil,
				nil,
				sendButton,
				messageEntry,
			),
			nil,
			nil,
			chatScroll,
		)

		// ---- User selection ----
		userList.OnSelected = func(i widget.ListItemID) {
			chatLabel.SetText("Chat with " + state.Chats[i].User.Username)
			err := state.SetCurrentChat(state.Chats[i])
			if err != nil {
				dialog.ShowError(err, w)
			}
		}

		// ---- Split view ----
		split := container.NewHSplit(
			userList,
			chatRight,
		)
		split.SetOffset(0.25) // 25% left, 75% right

		// ---- Top bar (optional) ----
		top := container.NewHBox(
			widget.NewLabel(fmt.Sprintf("Welcome to secure chat!")),
			layout.NewSpacer(),
			widget.NewButton("Add user", func() {
				showAddUserScreen()
			}),
			widget.NewButton("Logout", func() {
				showLogin()
				state.ClearAuthState()
			}),
		)

		w.SetContent(container.NewBorder(top, nil, nil, nil, split))
	}

	showAddUserScreen = func() {
		searchEntry := widget.NewEntry()
		searchEntry.SetPlaceHolder("Enter username")

		var results []client.UserListItemDto

		resultList := widget.NewList(
			func() int { return len(results) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(results[i].Username)
			},
		)

		selectedIndex := -1
		resultList.OnSelected = func(id widget.ListItemID) {
			selectedIndex = id
		}

		searchBtn := widget.NewButton("Search", func() {
			query := searchEntry.Text
			if query == "" { // TODO: call backend to add user

				dialog.ShowInformation("Search", "Please enter a username", w)
				return
			}

			res, err := service.Search(query)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			results = res
			resultList.Refresh()
		})

		addBtn := widget.NewButton("Add to chat", func() {
			if selectedIndex < 0 || selectedIndex >= len(results) {
				dialog.ShowInformation("Select user", "Please select a user from the list", w)
				return
			}

			addErr := service.AddUser(results[selectedIndex])
			if addErr != nil {
				dialog.ShowError(addErr, w)
				return
			}
			dialog.ShowInformation("Add user",
				fmt.Sprintf("User %s added to chat", results[selectedIndex].Username), w)
			showMainScreen() // go back to main screen
		})

		backBtn := widget.NewButton("Back", func() {
			showMainScreen() // go back to main screen
		})

		// ---- Top bar with back button ----
		topBar := container.NewHBox(layout.NewSpacer(), backBtn)

		// ---- Scrollable list of results ----
		scroll := container.NewVScroll(resultList)
		scroll.SetMinSize(fyne.NewSize(0, 300)) // sets initial height

		// ---- Main content VBox ----
		mainContent := container.NewVBox(
			widget.NewLabel("Search users"),
			searchEntry,
			searchBtn,
			scroll,
		)

		// ---- Final layout: top bar + main content + add button pinned at bottom ----
		content := container.NewBorder(
			topBar,      // top
			addBtn,      // bottom
			nil,         // left
			nil,         // right
			mainContent, // center
		)

		w.SetContent(content)
	}

	// start with login screen
	showLogin()
	w.ShowAndRun()
}
