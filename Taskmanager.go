package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	Id        int32  `gorm:"primarykey"`
	User      string `gorm:"not null"`
	TaskName  string `gorm:"column:task_name"`
	Date      string `gorm:"column:startdate"`
	EndDate   string `gorm:"column:endDate"`
	Completed bool   `gorm:"column:Completed"`
	Priority  string `gorm:"column:priority"`
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("taskmanager_db"), &gorm.Config{})
	if err != nil {
		fmt.Println("can not open the DB", err)
	}
	err = db.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal("Cannot create the db")
		//fmt.Println("Cannot create the db")
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

func AddFunc(db *gorm.DB, w *sync.WaitGroup) { //Adding task into the Database

	defer w.Done()
	var taskname string
	var user string
	var priority string
	var endDate int
	fmt.Println("Enter TaskName")
	fmt.Scanln(&taskname)
	fmt.Println("Enter username")
	fmt.Scanln(&user)
	fmt.Println("Set Priority (low/high)")
	fmt.Scan(&priority)
	fmt.Println("Enter End Date")
	fmt.Scan(&endDate)
	date := time.Now().Format("2006-01-02")
	dueDate := time.Now().AddDate(0, 0, endDate).Format("2006-01-02")
	task := Task{
		User:      user,
		TaskName:  taskname,
		Date:      date,
		Completed: false,
		EndDate:   dueDate,
		Priority:  priority,
	}
	result := db.Create(&task)
	if result.Error != nil {
		log.Error("error adding task")
	}
	fmt.Println("Task added succesfully")
}
func Viewtask(db *gorm.DB, username string, w *sync.WaitGroup) {
	defer w.Done()
	if UsernameExists(db, username) {
		var tasks []Task
		result := db.Where("User=?", username).Find(&tasks)
		if result.Error != nil {
			log.Error("Cannot fetch the data")
		}
		fmt.Println("=== All Tasks ===")
		for _, task := range tasks {
			output := fmt.Sprintf("Username:= %v, TaskName:= %v, Completed:= %v, TaskPriority:= %v, StartDate:= %v, EndDate:= %v", task.User, task.TaskName, task.Completed, task.Priority, task.Date, task.EndDate)
			fmt.Println(output)
		}
	} else {
		log.Error("Cannot find the username")
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

func Deletetask(db *gorm.DB) {
	var task Task
	var username string
	var taskname string
	fmt.Println("Enter username")
	fmt.Scan(&username)
	fmt.Println("Enter taskname")
	fmt.Scan(&taskname)
	found := db.Where("User=? AND task_name=?", username, taskname).First(&task)
	if found.RowsAffected > 0 {
		db.Delete(&task)
		fmt.Println("Task deleted succesfully", found)
	}
	log.Error("User or task name not found")
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
		fmt.Println("1> ADD TASK")
		fmt.Println("2> VIEW TASK")
		fmt.Println("3> QUIT")
		fmt.Println("4> DELETE TASK")
		fmt.Println("5> SEARCH TASK")
		fmt.Println("6> UPDATE EXSISTING TASK")
		fmt.Println("7> MARK TASK AS COMPLETED")
		fmt.Println("8> SORT TASKS PRIORITY WISE")
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
			Deletetask(db)
		case "5":
			fmt.Println("Task you want to Search")
			var taskname string
			fmt.Scan(&taskname)
			Search(db, taskname)
		case "6":
			Update(db)
		case "7":
			wg.Add(1)
			go MarkComplete(db, &wg)
		case "8":
			SortbasedPriority(db)
		default:
			fmt.Println("Select a correct Choice")
		}
		wg.Wait()
	}
}

func SortbasedPriority(db *gorm.DB) {
	var tasks []Task
	var username string
	var priority string
	fmt.Println("Enter the Username")
	fmt.Scan(&username)
	fmt.Println("Enter the priority(high/low)")
	fmt.Scan(&priority)
	found := db.Find(&tasks, "User=? AND Priority=?", username, priority)
	if found.RowsAffected > 0 {
		fmt.Println(tasks)
	} else {
		log.Error("No user found")
	}
}

func MarkComplete(db *gorm.DB, w *sync.WaitGroup) {
	defer w.Done()
	var username string
	var taskname string
	var tasks Task
	fmt.Println("enter username")
	fmt.Scan(&username)
	fmt.Println("enter taskname")
	fmt.Scan(&taskname)
	result := db.Where("user=? and task_name=?", username, taskname).First(&tasks)
	if result.RowsAffected > 0 {
		db.Model(&tasks).Update("completed", true)
		fmt.Println("task completed good job")
	} else {
		log.Error("task cannot be mark as completed")
	}
}
