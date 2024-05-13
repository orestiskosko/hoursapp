package models

import (
	"fmt"
	"math"
	"time"
)

type TaskEntryViewModel struct {
	ID          int
	TaskID      int
	TaskName    string
	ProjectID   int
	ProjectName string
	StartedAt   string
	Duration    string
}

type TrackerViewModel struct {
	IsRunning        bool
	Date             string
	RunningTaskEntry TaskEntryViewModel
	TaskEntries      []TaskEntryViewModel
	ProjectOptions   map[string]string
}

type ViewModelMapper[T interface{}] interface {
	ToViewModel() T
}

func (taskEntry TaskEntry) ToViewModel() TaskEntryViewModel {
	if taskEntry == (TaskEntry{}) {
		return TaskEntryViewModel{}
	}

	duration := "00:00:00"
	if taskEntry.FinishedAt != nil {
		durationInSeconds := taskEntry.FinishedAt.Sub(taskEntry.StartedAt).Seconds()
		duration = fmt.Sprintf(
			"%02d:%02d:%02d",
			int(math.Floor(durationInSeconds)/3600),
			int(math.Floor((float64(int(durationInSeconds)%3600) / 60.0))),
			int(durationInSeconds)%60)
	}

	return TaskEntryViewModel{
		ID:          taskEntry.ID,
		TaskID:      taskEntry.TaskID,
		TaskName:    taskEntry.Task.Name,
		ProjectID:   taskEntry.Task.Project.ID,
		ProjectName: taskEntry.Task.Project.Name,
		StartedAt:   taskEntry.StartedAt.Format(time.RFC3339),
		Duration:    duration,
	}
}

func ToViewModels[S ViewModelMapper[T], T interface{}](s []S) []T {
	targetSlice := make([]T, 0)
	for _, val := range s {
		targetSlice = append(targetSlice, val.ToViewModel())
	}
	return targetSlice
}
