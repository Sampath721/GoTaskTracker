package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

type Task struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskTracker struct {
	Tasks []Task `json:"tasks"`
}

func main() {

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	saveCmd := flag.NewFlagSet("save", flag.ExitOnError)
	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)

	addTitle := addCmd.String("title", "", "Task title")
	addDescription := addCmd.String("description", "", "Task description")
	addPriority := addCmd.String("priority", "low", "Priority of the task (low, medium, high)")

	completeIndex := completeCmd.Int("index", -1, "Task index")

	removeIndex := removeCmd.Int("index", -1, "Task index")

	saveFilename := saveCmd.String("filename", "tasks.json", "File name to save tasks")

	loadFilename := loadCmd.String("filename", "tasks.json", "File name to load tasks")
	listPriority := listCmd.String("priority", "", "Filter tasks by priority (low, medium, high)")

	if len(os.Args) < 2 {
		fmt.Println("Please specify a command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		addTask(*addTitle, *addDescription, *addPriority)
	case "list":
		listCmd.Parse(os.Args[2:])
		listTasks(*listPriority)
	case "complete":
		completeCmd.Parse(os.Args[2:])
		completeTask(*completeIndex)
	case "remove":
		removeCmd.Parse(os.Args[2:])
		removeTask(*removeIndex)
	case "save":
		saveCmd.Parse(os.Args[2:])
		saveTasks(*saveFilename)
	case "load":
		loadCmd.Parse(os.Args[2:])
		loadTasks(*loadFilename)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

//add Tasks

func addTask(title, description, Priority string) {
	task := Task{
		Title:       title,
		Description: description,
		Status:      "incomplete",
		Priority:    Priority,
		CreatedAt:   time.Now(),
	}

	tasks := loadTasksFromStorage()
	tasks.Tasks = append(tasks.Tasks, task)

	saveTasksToStorage(tasks)

	fmt.Println("Task added successfully!")
}

//List All Tasks

func listTasks(priorityFilter string) {
	tasks := loadTasksFromStorage()

	// Sort tasks by priority using predefined weight
	priorityWeight := map[string]int{"high": 3, "medium": 2, "low": 1}
	sort.Slice(tasks.Tasks, func(i, j int) bool {
		return priorityWeight[tasks.Tasks[i].Priority] > priorityWeight[tasks.Tasks[j].Priority]
	})

	// Create a new tab writer
	writer := tabwriter.NewWriter(os.Stdout, 8, 8, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(writer, "ID\tPriority\tTitle\tDescription\tStatus\tCreated At")

	// Iterate over tasks and print them
	for i, task := range tasks.Tasks {
		if priorityFilter == "" || task.Priority == priorityFilter {
			fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\t%s\n",
				i+1,
				task.Priority,
				task.Title,
				task.Description,
				task.Status,
				task.CreatedAt.Format("2006-01-02 15:04:05"),
			)
		}
	}
	writer.Flush()
}

// Mark Task As Complete

func completeTask(index int) {
	tasks := loadTasksFromStorage()

	if index < 0 || index >= len(tasks.Tasks) {
		fmt.Println("Invalid task index")
		return
	}

	tasks.Tasks[index].Status = "completed"

	saveTasksToStorage(tasks)

	fmt.Println("Task marked as completed!")
}

//remove Task

func removeTask(index int) {
	tasks := loadTasksFromStorage()

	if index < 0 || index >= len(tasks.Tasks) {
		fmt.Println("Invalid task index")
		return
	}

	tasks.Tasks = append(tasks.Tasks[:index], tasks.Tasks[index+1:]...)

	saveTasksToStorage(tasks)

	fmt.Println("Task removed successfully!")
}

//Save Tasks Useful For Batch Save

func saveTasks(filename string) {
	tasks := loadTasksFromStorage()

	data, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		log.Fatal("Failed to serialize tasks:", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal("Failed to save tasks:", err)
	}
	fmt.Println("Tasks saved to file!")
}

//Load All Tasks

func loadTasks(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed to load tasks:", err)
	}

	var tasks TaskTracker
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		log.Fatal("Failed to deserialize tasks:", err)
	}

	saveTasksToStorage(tasks)

	fmt.Println("Tasks loaded from file!")
}

func loadTasksFromStorage() TaskTracker {
	var tasks TaskTracker

	data, err := os.ReadFile("tasks.json")
	if err != nil {
		return tasks
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return tasks
	}

	return tasks
}

func saveTasksToStorage(tasks TaskTracker) {
	data, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		log.Fatal("Failed to serialize tasks:", err)
	}

	err = os.WriteFile("tasks.json", data, 0644)
	if err != nil {
		log.Fatal("Failed to save tasks:", err)
	}
}
