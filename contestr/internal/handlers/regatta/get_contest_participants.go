package regatta

import (
	"context"
	"net/http"

	"contestr/internal/repository"
	"contestr/pkg/logger"

	"github.com/labstack/echo/v4"
)

type ParticipantItem struct {
	ParticipantID string `json:"participant_id"`
	DisplayName   string `json:"display_name"`
}

type HandleLister interface {
	ListByContestID(ctx context.Context, contestID int) ([]repository.CodeforcesHandleMapping, error)
}

type ParticipantsHandle struct {
	handles HandleLister
}

func NewParticipantsHandle(handles HandleLister) *ParticipantsHandle {
	return &ParticipantsHandle{handles: handles}
}

func (h *ParticipantsHandle) GetContestParticipants(ctx echo.Context, contestId int) error {
	mappings, err := h.handles.ListByContestID(ctx.Request().Context(), contestId)
	if err != nil {
		logger.Errorf(ctx.Request().Context(), "error while listing contest participants: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	items := make([]ParticipantItem, 0, len(mappings))
	for _, m := range mappings {
		displayName := m.Name
		if displayName == "" {
			displayName = m.Handle
		}
		items = append(items, ParticipantItem{
			ParticipantID: m.Handle,
			DisplayName:   displayName,
		})
	}
	return ctx.JSON(http.StatusOK, items)
}
