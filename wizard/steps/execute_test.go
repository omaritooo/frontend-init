package steps_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/omaritooo/frontend-init/wizard/steps"
)

func TestExecuteStep_ViewShowsAllTaskLabels(t *testing.T) {
	tasks := []steps.TaskStatus{
		{Label: "Scaffold project", State: steps.TaskDone},
		{Label: "Install packages", State: steps.TaskRunning},
		{Label: "Configure Tailwind", State: steps.TaskPending},
	}
	s := steps.NewExecuteStep(tasks)
	view := s.View()
	assert.Contains(t, view, "Scaffold project")
	assert.Contains(t, view, "Install packages")
	assert.Contains(t, view, "Configure Tailwind")
}

func TestExecuteStep_NotDoneInitially(t *testing.T) {
	tasks := []steps.TaskStatus{
		{Label: "Install packages", State: steps.TaskPending},
	}
	s := steps.NewExecuteStep(tasks)
	assert.False(t, s.IsDone())
}

func TestExecuteStep_DoneWhenAllTasksComplete(t *testing.T) {
	tasks := []steps.TaskStatus{
		{Label: "Step A", State: steps.TaskPending},
		{Label: "Step B", State: steps.TaskPending},
	}
	s := steps.NewExecuteStep(tasks)
	// mark step 0 done
	s2, _ := s.Update(steps.TaskProgressMsg{Index: 0, State: steps.TaskDone})
	assert.False(t, s2.IsDone())
	// mark step 1 done
	s3, _ := s2.Update(steps.TaskProgressMsg{Index: 1, State: steps.TaskDone})
	assert.True(t, s3.IsDone())
}

func TestExecuteStep_LabelIsExecute(t *testing.T) {
	s := steps.NewExecuteStep(nil)
	assert.Equal(t, "execute", s.Label())
}
