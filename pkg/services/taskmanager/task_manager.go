package taskmanager

import (
	"mini-trh-backend/internal/logger"

	"go.uber.org/zap"
)

type Task struct {
	Name string
	Fn   func() error
}

type TaskManager struct {
	taskQueue   chan Task
	workerCount int
}

func NewTaskManager(workerCount, bufferSize int) *TaskManager {
	return &TaskManager{
		taskQueue:   make(chan Task, bufferSize),
		workerCount: workerCount,
	}
}

func (tm *TaskManager) Start() {
	logger.Log.Info("Starting task manager", zap.Int("workders", tm.workerCount))

	for i := 0; i < tm.workerCount; i++ {
		go tm.worker(i)
	}
}

func (tm *TaskManager) worker(id int) {
	logger.Log.Info("Worker started", zap.Int("workder_id", id))

	for task := range tm.taskQueue {
		logger.Log.Info("Worker processing task",
			zap.Int("worker_id", id),
			zap.String("task_name", task.Name),
		)

		if err := task.Fn(); err != nil {
			logger.Log.Error("Task failed",
				zap.Int("worker_id", id),
				zap.String("task_name", task.Name),
				zap.Error(err),
			)
		} else {
			logger.Log.Info("Task completed successfully",
				zap.Int("worker_id", id),
				zap.String("task_name", task.Name),
			)
		}
	}
}

// AddTask adds a new task to the queue
func (tm *TaskManager) AddTask(task Task) {
	logger.Log.Info("Adding task to queue", zap.String("task_name", task.Name))
	tm.taskQueue <- task
}

// Stop stops the task manager and waits for tasks to complete
func (tm *TaskManager) Stop() {
	logger.Log.Info("Stopping task manager")
	close(tm.taskQueue)
}
