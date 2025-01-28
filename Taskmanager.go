package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

type Task struct {
	Id        int32  `gorm:"primarykey"`
	User      string `gorm:"not null"`
	TaskName  string
	Date      time.Time
	Completed bool
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("taskmanager_db"), &gorm.Config{})
	if err != nil {
		fmt.Println("can not open the DB", err)
	}
	err = db.AutoMigrate(&Task{})
	if err != nil {
		fmt.Println("Cannot create the db")
	}
	return db
}

func AddFunc(db *gorm.DB) {

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
	result := db.Create(&task)
	if result.Error != nil {
		fmt.Println("error adding task")
	}
	fmt.Println("Task added succesfully")
}
func Viewtask(db *gorm.DB) {
	var tasks []Task
	result := db.Find(&tasks)
	if result.Error != nil {
		fmt.Println("Cannot fetch the data")
	}
	fmt.Println("=== All Tasks ===")
	for _, task := range tasks {
		fmt.Println(task.Id, task.TaskName, task.Date)
	}
}
func Deletetask(db *gorm.DB, todelId int) {
	var task Task

	_ = db.First(&task, todelId)
	//if err != nil {
	//fmt.Println("Error in finding task")
	//}
	_ = db.Delete(&task)
	//if err != nil {
	//fmt.Println("Error in deleting the task")
	//}
	fmt.Println("Task deleted succesfully", task)

}
func Search(db *gorm.DB, taskname string) {
	var task Task
	_ = db.First(&task)
	fmt.Println("Task Found", task)
}

func main() {
	db := initDB()
	//m := make(map[string]Task)

	fmt.Println("Hello")
	for {
		fmt.Println("Enter your Choices...")
		fmt.Println("1> Add task")
		fmt.Println("2> View Task")
		fmt.Println("3> Quit")
		fmt.Println("4> Delete Task")
		fmt.Println("5> Search Task")
		var userInput string
		fmt.Scanln(&userInput)
		switch userInput {
		case "1":
			AddFunc(db)
		case "2":
			Viewtask(db)
		case "3":
			os.Exit(0)
		case "4":
			var todel int
			fmt.Println("Enter the task Id to delete")
			fmt.Scan(&todel)
			Deletetask(db, todel)
		case "5":
			fmt.Println("Task you want to Search")
			var taskname string
			fmt.Scan(&taskname)
			Search(db, taskname)
		default:
			fmt.Println("Select a correct Choice")
		}
	}
}
