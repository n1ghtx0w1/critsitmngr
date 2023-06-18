package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type BlogPost struct {
	ID                 int
	CustomerName       string
	ContactNumber      string
	SeverityLevel      string
	ProblemStatement   string
	SolutionActionPlan string
}

var blogPosts []BlogPost
var currentID int

const dataFile = "blog_data.json"

func main() {
	loadData()

	myApp := app.New()
	myWindow := myApp.NewWindow("Critical Situation App")

	customerNameEntry := widget.NewEntry()
	customerNameEntry.SetPlaceHolder("Enter Customer Name")

	contactNumberEntry := widget.NewEntry()
	contactNumberEntry.SetPlaceHolder("Enter Customer Contact")

	severityLevelEntry := widget.NewEntry()
	severityLevelEntry.SetPlaceHolder("Enter Severity Level")

	problemStatementEntry := widget.NewEntry()
	problemStatementEntry.SetPlaceHolder("Enter Problem Statement")

	solutionActionPlanEntry := widget.NewEntry()
	solutionActionPlanEntry.SetPlaceHolder("Enter Solution Action Plan")

	submitButton := widget.NewButton("Submit", func() {
		customerName := customerNameEntry.Text
		contactNumber := contactNumberEntry.Text
		severityLevel := severityLevelEntry.Text
		problemStatement := problemStatementEntry.Text
		solutionActionPlan := solutionActionPlanEntry.Text

		blogPost := BlogPost{
			ID:                 currentID,
			CustomerName:       customerName,
			ContactNumber:      contactNumber,
			SeverityLevel:      severityLevel,
			ProblemStatement:   problemStatement,
			SolutionActionPlan: solutionActionPlan,
		}

		blogPosts = append(blogPosts, blogPost)

		currentID++

		customerNameEntry.SetText("")
		contactNumberEntry.SetText("")
		severityLevelEntry.SetText("")
		problemStatementEntry.SetText("")
		solutionActionPlanEntry.SetText("")

		saveData()

		dialog.ShowInformation("Success", "Blog post saved!", myWindow)
	})

	reviewButton := widget.NewButton("Review", func() {
		postListWindow := myApp.NewWindow("Blog Posts")
		postListWindow.Resize(fyne.NewSize(400, 450))

		var postItems []fyne.CanvasObject
		for _, post := range blogPosts {
			post := post // Create a local variable to capture the value of post

			item := widget.NewButton(fmt.Sprintf("ID: %d - %s", post.ID, post.CustomerName), func() {
				showPostDetails(post, postListWindow)
			})

			postItems = append(postItems, item)
		}

		if len(postItems) == 0 {
			dialog.ShowInformation("No Posts", "No blog posts found.", postListWindow)
			return
		}

		postList := container.NewVBox(postItems...)
		scrollContainer := container.NewScroll(postList)

		postListWindow.SetContent(scrollContainer)
		postListWindow.Show()
	})

	deleteButton := widget.NewButton("Delete", func() {
		if len(blogPosts) == 0 {
			dialog.ShowInformation("No Posts", "No blog posts to delete.", myWindow)
			return
		}

		postListWindow := myApp.NewWindow("Delete Blog Posts")
		postListWindow.Resize(fyne.NewSize(400, 450))

		postList := widget.NewList(
			func() int {
				return len(blogPosts)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(index int, item fyne.CanvasObject) {
				if label, ok := item.(*widget.Label); ok {
					post := blogPosts[index]
					label.SetText(fmt.Sprintf("ID: %d - Customer: %s", post.ID, post.CustomerName))
				}
			},
		)

		postList.OnSelected = func(index int) {
			if index >= 0 && index < len(blogPosts) {
				deletePost(blogPosts[index], postListWindow, postList)
			}
		}

		scrollContainer := container.NewScroll(postList)

		postListWindow.SetContent(scrollContainer)
		postListWindow.Show()
	})

	itilButton := widget.NewButton("ITIL", func() {
		urlString := "https://headsec.tech/posts/itil4/" // Replace with your desired URL
		u, err := url.Parse(urlString)
		if err != nil {
			log.Println("Error parsing URL:", err)
			return
		}
		err = fyne.CurrentApp().OpenURL(u)
		if err != nil {
			log.Println("Error opening URL:", err)
			return
		}
	})

	myWindow.Resize(fyne.NewSize(400, 450)) // Resize the window

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Customer Name:"),
		customerNameEntry,
		widget.NewLabel("Contact Number:"),
		contactNumberEntry,
		widget.NewLabel("Severity Level:"),
		severityLevelEntry,
		widget.NewLabel("Problem Statement:"),
		problemStatementEntry,
		widget.NewLabel("Solution Action Plan:"),
		solutionActionPlanEntry,
		submitButton,
		reviewButton,
		deleteButton,
		itilButton, // Add the ITIL button here
	))

	myWindow.ShowAndRun()
}

func loadData() {
	fileData, err := ioutil.ReadFile(dataFile)
	if err != nil {
		log.Println("Error reading data file:", err)
		return
	}

	err = json.Unmarshal(fileData, &blogPosts)
	if err != nil {
		log.Println("Error unmarshaling data:", err)
		return
	}

	if len(blogPosts) > 0 {
		currentID = blogPosts[len(blogPosts)-1].ID + 1
	}

	log.Println("Data loaded successfully")
}

func saveData() {
	data, err := json.Marshal(blogPosts)
	if err != nil {
		log.Println("Error marshaling data:", err)
		return
	}

	err = ioutil.WriteFile(dataFile, data, 0644)
	if err != nil {
		log.Println("Error writing data file:", err)
		return
	}

	log.Println("Data saved successfully")
}

func showPostDetails(post BlogPost, window fyne.Window) {
	postDetails := fmt.Sprintf("ID: %d\nCustomer Name: %s\nContact Number: %s\nSeverity Level: %s\nProblem Statement: %s\nSolution Action Plan: %s",
		post.ID, post.CustomerName, post.ContactNumber, post.SeverityLevel, post.ProblemStatement, post.SolutionActionPlan)

	dialog.ShowInformation("Blog Post Details", postDetails, window)
}

func deletePost(post BlogPost, window fyne.Window, postList *widget.List) {
	confirmation := dialog.NewCustomConfirm("Confirm Delete", "Are you sure you want to delete this blog post?", "Delete", widget.NewLabel("Cancel"), func(result bool) {
		if result {
			for i, p := range blogPosts {
				if p.ID == post.ID {
					blogPosts = append(blogPosts[:i], blogPosts[i+1:]...)
					saveData()

					// Refresh the postList widget
					postList.Refresh()

					dialog.ShowInformation("Success", "Blog post deleted!", window)
					return
				}
			}
		}
	}, window)

	confirmation.Show()
}
