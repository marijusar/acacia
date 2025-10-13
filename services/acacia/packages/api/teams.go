package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"acacia/packages/auth"
	"acacia/packages/db"
	"acacia/packages/httperr"
	"acacia/packages/schemas"

	"github.com/go-playground/validator/v10"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
)

type TeamsController struct {
	queries   *db.Queries
	logger    *logrus.Logger
	validator *validator.Validate
}

func NewTeamsController(queries *db.Queries, logger *logrus.Logger) *TeamsController {
	return &TeamsController{
		queries:   queries,
		logger:    logger,
		validator: validator.New(),
	}
}

func (c *TeamsController) CreateTeam(w http.ResponseWriter, r *http.Request) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := auth.GetUserID(r)
	if !ok {
		return httperr.WithStatus(errors.New("Unauthorized"), http.StatusUnauthorized)
	}

	var req schemas.CreateTeamInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	// Validate input
	if err := c.validator.Struct(req); err != nil {
		return httperr.WithStatus(schemas.HandleTeamValidationErrors(err), http.StatusBadRequest)
	}

	// Create team
	params := db.CreateTeamParams{
		Name:        req.Name,
		Description: null.String{}, // Empty description
	}

	team, err := c.queries.CreateTeam(r.Context(), params)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create team")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Add creator as team member
	_, err = c.queries.AddTeamMember(r.Context(), db.AddTeamMemberParams{
		TeamID: team.ID,
		UserID: userID,
	})
	if err != nil {
		c.logger.WithError(err).Error("Failed to add team member")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
	return nil
}

func (c *TeamsController) GetUserTeams(w http.ResponseWriter, r *http.Request) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := auth.GetUserID(r)
	if !ok {
		return httperr.WithStatus(errors.New("Unauthorized"), http.StatusUnauthorized)
	}

	teams, err := c.queries.GetUserTeams(r.Context(), userID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get user teams")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Ensure we always return an array, not null
	// When no teams are found, teams will be nil, so initialize empty slice
	if teams == nil {
		teams = []db.Team{}
	}

	json.NewEncoder(w).Encode(teams)
	return nil
}
