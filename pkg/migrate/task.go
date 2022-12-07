package migrate

import (
	"context"
	"fmt"
	"os"

	"github.com/go-task/task/v3"
	"github.com/go-task/task/v3/taskfile"
)

// Task entrypoint executor.
type Task struct {
	entrypoint string
	verbose    bool
	color      bool
	stdin      *os.File
	stdout     *os.File
	stderr     *os.File
}

// NewTask creates an instance of Task.
func NewTask(
	entrypoint string,
	verbose bool,
) *Task {
	color := true
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	return &Task{
		entrypoint: entrypoint,
		verbose:    verbose,
		color:      color,
		stdin:      stdin,
		stdout:     stdout,
		stderr:     stderr,
	}
}

// Run execute the given task file.
func (t *Task) Run() error {
	e := &task.Executor{
		Color:      t.color,
		Verbose:    t.verbose,
		Entrypoint: t.entrypoint,
		Stdin:      t.stdin,
		Stdout:     t.stdout,
		Stderr:     t.stderr,
	}

	if err := e.Setup(); err != nil {
		return fmt.Errorf("failed to setup task: %w", err)
	}

	ctx := context.Background()
	if err := e.Run(ctx, taskfile.Call{Task: "up"}); err != nil {
		if err, ok := err.(*task.TaskRunError); ok {
			return fmt.Errorf("failed to execute task: %w", err)
		}
	}

	return nil
}
