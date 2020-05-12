# Go-Routine-Control

The issue with Go routines and their use in tasks (in sneaker bots) is stopping them when a user requests it. As you know modern bots have stop functions for each task, which sparked me to create this for Go.

## Use of channels

I have used channels to signal a task, or go routine to stop. Two new channels are created when you start a task:
 - TaskIDStop
 - TaskIDStopped

These channels are unique. They are then stored in an array, known as "taskArray". I use a struct to form this array, which consists of your Task ID and your two unique channels. 

## Starting the task
To start a task, you would create both channels and a Task ID. This item would then be appended to taskArray, and then you call your function:
```go
go Task(task.taskIDStop, task.taskIDStopped, &waitgroup)
```
Both task.taskIDStop and task.taskIDStopped are **chan struct{}**. You do not have to use struct, and can instead use **chan string** if you want. 

After you have started your task, you can now stop it at any time.

## Calling the stop function

To stop the task you call: 
```go
stop(TaskID)
```

Beneath the hood, this is:
```go
Find index of TaskID in taskArray (i)
close(taskArray[i].taskIDStop)
<-taskArray[i].taskIDStopped
Return task stopped to UI
Remove task from taskArray
```
This finds the specified task in your task array, and the related channels. It then closes the taskIDStop channel. 

## Handling the stop
To know when the channel has closed, for each operation in your function you would check using a select statement, like so:
```go

defer close(taskIDStopped)

Select {
default:
 Do task etc
case <-taskIDStop:
 waitgroup.Done()
 return
 //Task exits
}
```
You must defer the closing of the taskIDStopped channel. Defering the closing will allow it to close the channel when your go routine exits, therefore signaling your go routine has stopped successfully.
