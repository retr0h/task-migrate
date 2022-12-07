package migrate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Migration is the top-level migration instance.
type Migration struct {
	Logger      *logrus.Logger
	repo        Repository
	versionsDir string
	verbose     bool

	ctx context.Context
}

// PlannedMigration migrations to be run.
type PlannedMigration struct {
	Tasks []*Task
}

// NewMigration creates an instance of Migration.
func NewMigration(
	logger *logrus.Logger,
	repository *Repository,
	versionsDir string,
	verbose bool,
) *Migration {
	ctx := context.Background()

	return &Migration{
		Logger:      logger,
		repo:        *repository,
		versionsDir: versionsDir,
		verbose:     verbose,

		ctx: ctx,
	}
}

// Status show migration status.
func (m *Migration) Status() (table.Writer, error) {
	records, err := m.repo.Find()
	if err != nil {
		return nil, err
	}

	t := table.NewWriter()
	t.Style().Color.Header = text.Colors{text.FgCyan}
	t.Style().Color.RowAlternate = text.Colors{text.FgBlack}
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"migration", "applied"})
	for _, record := range records {
		t.AppendRow([]interface{}{record.Name, record.CreatedAt})
	}

	return t, nil
}

// Up migrates to the most recent version available.
func (m *Migration) Up() error {
	plannedMigration, err := m.FindMigrations()
	if err != nil {
		return err
	}

	if len(plannedMigration.Tasks) == 0 {
		m.Logger.Info("no migrations to run")

		return nil
	}

	for _, task := range plannedMigration.Tasks {
		err := task.Run()
		if err != nil {
			return err
		}

		baseName := filepath.Base(task.entrypoint)
		err = m.repo.Insert(baseName)
		if err != nil {
			return err
		}
	}

	return nil
}

// FindMigrations determine planned migrations.
func (m *Migration) FindMigrations() (*PlannedMigration, error) {
	var lastMigrationVersionID int64

	lastMigrationRecord, err := m.repo.Last()
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	if lastMigrationRecord != nil {
		lastMigrationVersionID = lastMigrationRecord.ID
	}

	files, err := os.ReadDir(m.versionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %s %w", m.versionsDir, err)
	}

	tasks := make([]*Task, len(files))
	for i, file := range files {
		versionID, err := strconv.ParseInt(file.Name()[0:1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert str to int: %w", err)
		}

		if versionID <= lastMigrationVersionID {
			continue
		}

		filePath := filepath.Join(m.versionsDir, file.Name())
		tasks[i] = NewTask(filePath, m.verbose)

	}

	return &PlannedMigration{
		Tasks: tasks,
	}, nil
}
