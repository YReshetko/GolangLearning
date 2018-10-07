package task

import (
	"sync"
)

type FjkTask interface {
	Execute()
	HasSubtask() bool
	Fork() []FjkTask
	Join([]FjkTask)
}

// TODO create thread pool
func Process(task FjkTask) {
	if task.HasSubtask() {
		var wg sync.WaitGroup
		log("Original task", task)
		subTasks := task.Fork()
		for _, subTask := range subTasks {
			log("Forked subtask", subTask)
			wg.Add(1)
			go func(tsk FjkTask) {
				defer wg.Done()
				Process(tsk)
				log("Processed subtask", tsk)
			}(subTask)
		}
		wg.Wait()
		task.Join(subTasks)
		log("Joined subtask", task)
	}
	task.Execute()
}

func log(msg string, tsk FjkTask) {
	//fmt.Printf("%s: %v\n", msg, tsk)
}
