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
var taskMap = make(map[int]taskControlStore)

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
	if val, ok := taskMap[taskID]; !ok {
		fmt.Println("Task does not exist.")
	} else if ok {
		//Close the stop channel which stops the task
		close(val.taskIDStop)

		//Wait for the stopped channel, which confirms it has closed
		<-val.taskIDStopped
		fmt.Println("Task stopped")
		//Remove from task array
		delete(taskMap, taskID)
	}
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
	taskMap[task.taskID] = task

	//Create a go routine for task
	waitgroup.Add(1)
	go Task(task.taskIDStop, task.taskIDStopped, &waitgroup)

	//Wait and then stop task
	time.Sleep(3000 * time.Millisecond)
	StopTask(task.taskID)

	waitgroup.Wait()
}
