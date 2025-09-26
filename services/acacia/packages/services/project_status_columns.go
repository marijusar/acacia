package services

import (
	"acacia/packages/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrCannotDeleteLastColumn = errors.New("cannot delete the last column in a project")

type ProjectStatusColumnService struct {
	queries *db.Queries
	db      *sql.DB
}

func NewProjectStatusColumnService(queries *db.Queries, database *sql.DB) *ProjectStatusColumnService {
	return &ProjectStatusColumnService{
		queries: queries,
		db:      database,
	}
}

func (s *ProjectStatusColumnService) DeleteProjectStatusColumnWithReorder(ctx context.Context, columnID int64) (*db.ProjectStatusColumn, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	// Get column info before deletion
	columnInfo, err := qtx.GetProjectStatusColumnByID(ctx, columnID)
	if err != nil {
		return nil, err
	}

	// Check if this is the last column in the project
	columnCount, err := qtx.GetProjectStatusColumnCountByProjectID(ctx, columnInfo.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project column count: %w", err)
	}

	if columnCount <= 1 {
		return nil, ErrCannotDeleteLastColumn
	}

	// Find next column for issue reassignment
	nextColumnID, err := qtx.GetNextColumnForReassignment(ctx, db.GetNextColumnForReassignmentParams{
		ProjectID:     columnInfo.ProjectID,
		PositionIndex: columnInfo.PositionIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get next column: %w", err)
	}

	err = qtx.ReassignIssuesFromColumn(ctx, db.ReassignIssuesFromColumnParams{
		SourceColumn: columnID,
		TargetColumn: nextColumnID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to reassign issues: %w", err)
	}

	// Shift columns left
	err = qtx.ShiftColumnsLeft(ctx, db.ShiftColumnsLeftParams{
		ProjectID:     columnInfo.ProjectID,
		PositionIndex: columnInfo.PositionIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to shift columns: %w", err)
	}

	// Delete the column
	deletedColumn, err := qtx.DeleteProjectStatusColumn(ctx, columnID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete column: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &deletedColumn, nil
}

