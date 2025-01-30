package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	Id        int32     `gorm:"primarykey"`
	User      string    `gorm:"not null"`
	TaskName  string    `gorm:"column:task_name"`
	Date      time.Time `gorm:"column:date"`
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
func UsernameExists(db *gorm.DB, username string) bool {
	var task Task
	ok := db.Where("User=?", username).First(&task)
	if ok.RowsAffected > 0 {
		return true
	}
	return false
}

func AddFunc(db *gorm.DB, w *sync.WaitGroup) {

	defer w.Done()
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
func Viewtask(db *gorm.DB, username string, w *sync.WaitGroup) {
	defer w.Done()
	if UsernameExists(db, username) {
		var tasks []Task
		result := db.Where("User=?", username).Find(&tasks)
		if result.Error != nil {
			fmt.Println("Cannot fetch the data")
		}
		fmt.Println("=== All Tasks ===")
		for _, task := range tasks {
			fmt.Println(task.User, task.TaskName, task.Date)
		}
	} else {
		log.Fatal("Cannot find the username")
	}
}
func Update(db *gorm.DB) {
	var tasks Task
	var username string
	var taskname string
	fmt.Println("Enter username")
	fmt.Scan(&username)
	fmt.Println("Enter taskname")
	fmt.Scan(&taskname)
	tochange := db.Where("User=? AND task_name=?", username, taskname).Find(&tasks)
	if tochange.RowsAffected > 0 {
		var newTaskname string
		fmt.Println("Enter new taskname")
		fmt.Scan(&newTaskname)
		db.Model(&tasks).Update("task_name", newTaskname)
		fmt.Println("succesfully Updated")
	} else {
		log.Fatal("User not avilable")
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
	found := db.Where("task_name=?", &taskname).First(&task)
	if found.RowsAffected > 0 {
		fmt.Println("Task Found", task)
	} else {
		log.Fatal("No record found")
	}
}

func main() {
	var wg sync.WaitGroup
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
		fmt.Println("6> Update Exsisting task")
		var userInput string
		fmt.Scanln(&userInput)
		switch userInput {
		case "1":
			wg.Add(1)
			go AddFunc(db, &wg)
		case "2":
			wg.Add(1)
			var username string
			fmt.Println("Enter your username")
			fmt.Scan(&username)
			Viewtask(db, username, &wg)
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
		case "6":
			Update(db)
		default:
			fmt.Println("Select a correct Choice")
		}
		wg.Wait()
	}
}
