package main

import (
	"fmt"
	"os"
	"time"
)

type Task struct {
	User      string
	TaskName  string
	Date      time.Time
	Completed bool
}

func AddFunc(DB *map[string]Task) {

	var taskname string
	var user string
	fmt.Println("Enter TaskName")
	fmt.Scanln(&taskname)
	fmt.Println("Enter username")
	fmt.Scanln(&user)
	date := time.Now()
	task := Task{
		User:      user,
		TaskName:  taskname,
		Date:      date,
		Completed: false,
	}
	(*DB)[user] = task
}
func Viewtask(DB *map[string]Task) {
	for key, item := range *DB {
		fmt.Printf("Key: %s\n", key) // Print the map key (optional)
		fmt.Println("Task Name:", item.TaskName)
		fmt.Println("Date:", item.Date)
		fmt.Println("Completed:", item.Completed)
		fmt.Println("---------------------------")

	}

}

func main() {
	m := make(map[string]Task)

	fmt.Println("Hello")
	for {
		fmt.Println("Enter your Choices...")
		fmt.Println("1> Add task")
		fmt.Println("2> View Task")
		fmt.Println("3> Quit")
		var userInput string
		fmt.Scanln(&userInput)
		switch userInput {
		case "1":
			AddFunc(&m)
		case "2":
			Viewtask(&m)
		case "3":
			os.Exit(0)
		default:
			fmt.Println("Select a correct Choice")
		}
	}
}
