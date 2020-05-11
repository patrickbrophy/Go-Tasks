package main

import (
	"fmt"
	"sync"
	"time"
)

//Struct for storing task details
type taskControlStore struct {
	taskID        int
	taskIDStop    chan struct{}
	taskIDStopped chan struct{}
}

//Array for tasks
var taskArray []taskControlStore

//Example
func Task(taskIDStop chan struct{}, taskIDStopped chan struct{}, waitgroup *sync.WaitGroup) {

	//Defer the closing of the "stopped" channel
	defer close(taskIDStopped)

	//Each select checks if the taskIDStop channel has been closed.
	select {
	default:
		fmt.Println("Started Task")
	case <-taskIDStop:
		waitgroup.Done()
		return
	}
	time.Sleep(1 * time.Second)

	select {
	default:
		fmt.Println("Adding to cart...")
	case <-taskIDStop:
		waitgroup.Done()
		return
	}
	time.Sleep(1 * time.Second)

	select {
	default:
		fmt.Println("Added to cart")
	case <-taskIDStop:
		waitgroup.Done()
		return
	}
	time.Sleep(1 * time.Second)

	//Call done
	waitgroup.Done()
}

//Stop task function
func StopTask(taskID int) {

	//Find task ID -> channels
	index, taskIDFound := FindInTaskArray(taskArray, taskID)
	if !taskIDFound {
		fmt.Println("Task does not exist.")
	} else if taskIDFound {
		//Close the stop channel which stops the task
		close(taskArray[index].taskIDStop)

		//Wait for the stopped channel, which confirms it has closed
		<-taskArray[index].taskIDStopped
		fmt.Println("Task stopped")
		//Remove from task array
		taskArray = append(taskArray[:index], taskArray[index+1:]...)
	}
}

//Find index of task ID for later use
func FindInTaskArray(slice []taskControlStore, val int) (int, bool) {
	for i, item := range slice {
		if item.taskID == val {
			return i, true
		}
	}
	return -1, false
}

func main() {

	var waitgroup sync.WaitGroup

	//Create channels for control
	newTaskIDStop := make(chan struct{})
	newTaskIDStopped := make(chan struct{})

	//Create new task
	task := taskControlStore{
		taskID:        5657675675,
		taskIDStop:    newTaskIDStop,
		taskIDStopped: newTaskIDStopped}

	//Add to array
	taskArray = append(taskArray, task)

	//Create a go routine for task
	waitgroup.Add(1)
	go Task(task.taskIDStop, task.taskIDStopped, &waitgroup)

	//Wait and then stop task
	time.Sleep(1200 * time.Millisecond)
	StopTask(task.taskID)

	waitgroup.Wait()
}
