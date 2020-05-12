# Go-Routine-Control

The issue with Go routines and their use in tasks (in sneaker bots) is stopping them when a user requests it. As you know modern bots have stop functions for each task, which sparked me to create this for Go.

## Use of channels

I have used channels to signal a task, or go routine to stop. Two new channels are created when you start a task:
 - TaskIDStop
 - TaskIDStopped

These channels are unique. They are then stored in an array, known as "TaskArray". I use a struct to form this array, which consists of your Task ID and your two unique channels. 

To stop the task you call: 
```go
stop(TaskID)
```

Beneath the hood, this is:
```go
Find index of TaskID in TaskArray (i)
close(taskArray[i].taskIDStop)
<-taskArray[i].taskIDStopped
```
