# Terminal Based Rolodex Built with [Tview](https://github.com/rivo/tview)

This is just an example app. It does not persist any data. Follow along with the article on the [Earthly blog](https://earthly.dev/blog/tui-app-with-go/)

# Try it out right now with Gitpod

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/jalletto/tui-go-example)

Nothing to pay for, just login with GitHub and you'll be up and running with this tutorial in a full Linux environment with VS Code and the Golang extensions.



```go
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


type Contact struct {
	firstName   string
	lastName    string
	email       string
	phoneNumber string
	state       string
	business    bool
}

var contacts = make([]Contact, 0)

// Tview
var pages = tview.NewPages()
var contactText = tview.NewTextView()
var app = tview.NewApplication()
var form = tview.NewForm()
var contactsList = tview.NewList().ShowSecondaryText(false)
var flex = tview.NewFlex()
var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(a) to add a new contactq \n(q) to quit")

func main() {
	
	contactsList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
		setConcatText(&contacts[index])
	})

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(contactsList, 0, 1, true).
			AddItem(contactText, 0, 4, false), 0, 6, false).
		AddItem(text, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
		} else if event.Rune() == 97 {
			form.Clear(true)
			addContactForm()
			pages.SwitchToPage("Add Contact")
		}
		return event
	})

	pages.AddPage("Menu", flex, true, true)
	pages.AddPage("Add Contact", form, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func addContactList() {
	contactsList.Clear()
	for index, contact := range contacts {
		contactsList.AddItem(contact.firstName+" "+contact.lastName, " ", rune(49+index), nil)
	}
}

func addContactForm() *tview.Form {

	contact := Contact{}

	form.AddInputField("First Name", "", 20, nil, func(firstName string) {
		contact.firstName = firstName
	})

	form.AddInputField("Last Name", "", 20, nil, func(lastName string) {
		contact.lastName = lastName
	})

	form.AddInputField("Email", "", 20, nil, func(email string) {
		contact.email = email
	})

	form.AddInputField("Phone", "", 20, nil, func(phone string) {
		contact.phoneNumber = phone
	})

	form.AddDropDown("State", states, 0, func(state string, index int) {
		contact.state = state
	})

	form.AddCheckbox("Business", false, func(business bool) {
		contact.business = business
	})

	form.AddButton("Save", func() {
		contacts = append(contacts, contact)
		addContactList()
		pages.SwitchToPage("Menu")
	})

	return form
}

func setConcatText(contact *Contact) {
	contactText.Clear()
	text := contact.firstName + " " + contact.lastName + "\n" + contact.email + "\n" + contact.phoneNumber
	contactText.SetText(text)
}

```

## Just a Draft

```go
package main

import (
	"fmt"
	"github.com/rivo/tview"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func getFiles(path string) []os.FileInfo {

	files, err := ioutil.ReadDir(path)
	
	if err != nil {
		log.Fatalf("Error getting files from %s: %v", path, err)
	}

	return files
}


// Add List
var fileNamesList = tview.NewList().ShowSecondaryText(false)

func main() {

	directoryPath, err := os.Getwd()

	if err != nil {
		log.Fatalf("Error getting current directory:", err)
	}
	
	if len(os.Args) > 1 {
		// If a command-line argument is provided, use it as the directory path
		directoryPath = os.Args[1]
	}

	// Get all the files
  files := getFiles(directoryPath)

	// Filenames
	fileNames := make(map[string]string)
	for _, file := range files {
		fileNames[file.Name()] = file.Name()
	}

	// Create a new tview application
	app := tview.NewApplication()

	// form to rename files 
	form := tview.NewForm()

	// Add files names to the form
	for _, file := range files {
		fileName := file.Name()
		form.AddInputField(fileName, fileName, 20, nil, func(newName string) {
			fileNames[fileName] = newName
		})
	}

	form.AddButton("Save", func() {
		reNameFiles(fileNames, directoryPath)
	})

	// Set up the layout
	flex := tview.NewFlex().
		// SetDirection(tview.FlexRow).
		// AddItem(tview.NewTextView().SetText("Files in the current directory").SetTextAlign(tview.AlignCenter), 1, 0, false).
		AddItem(form, 0, 1, true).
		AddItem(fileNamesList, 0, 1, true)

	// Start the application
	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}


func reNameFiles(fileNames map[string]string, directory string) {
	for oldName, newName := range fileNames {
		fileNamesList.AddItem(oldName + " -> " + newName, " ", 43 , nil)
		
		oldFilePath := filepath.Join(directory, oldName)
		newFilePath := filepath.Join(directory, newName)

		err := os.Rename(oldFilePath, newFilePath)

		if err != nil {
			fmt.Printf("Error renaming file: %v\n", err)
		}
	}
}
```