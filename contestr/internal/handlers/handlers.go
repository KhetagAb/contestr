package handlers

import (
	"contestr/internal/handlers/admin"
	"contestr/internal/handlers/codeforces"
	"contestr/internal/handlers/regatta"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	contestHandle        *codeforces.ContestHandle
	regattaContestHandle *regatta.ContestHandle
	adminLoginHandle     *admin.LoginHandle
	adminMeHandle        *admin.MeHandle
	adminTimetableHandle *admin.TimetableHandle
}

func (h *Handlers) GetRegattaContestStandings(ctx echo.Context, contestId int) error {
	return h.regattaContestHandle.GetContest(ctx, contestId)
}

func NewHandlers(
	contestHandle *codeforces.ContestHandle,
	regattaContestHandle *regatta.ContestHandle,
	adminLoginHandle *admin.LoginHandle,
	adminMeHandle *admin.MeHandle,
	adminTimetableHandle *admin.TimetableHandle,
) *Handlers {
	return &Handlers{
		contestHandle:        contestHandle,
		regattaContestHandle: regattaContestHandle,
		adminLoginHandle:     adminLoginHandle,
		adminMeHandle:        adminMeHandle,
		adminTimetableHandle: adminTimetableHandle,
	}
}

func (h *Handlers) GetContest(ctx echo.Context, contestId int) error {
	return h.contestHandle.GetContest(ctx, contestId)
}

func (h *Handlers) PostAdminAuthLogin(ctx echo.Context) error {
	return h.adminLoginHandle.PostAdminAuthLogin(ctx)
}

func (h *Handlers) GetAdminMe(ctx echo.Context) error {
	return h.adminMeHandle.GetAdminMe(ctx)
}

func (h *Handlers) GetAdminTimetable(ctx echo.Context, contestId int) error {
	return h.adminTimetableHandle.GetAdminTimetable(ctx, contestId)
}

func (h *Handlers) PutAdminTimetable(ctx echo.Context, contestId int) error {
	return h.adminTimetableHandle.PutAdminTimetable(ctx, contestId)
}

func (h *Handlers) DeleteAdminTimetable(ctx echo.Context, contestId int) error {
	return h.adminTimetableHandle.DeleteAdminTimetable(ctx, contestId)
}

func (h *Handlers) PostAdminTimetableTourStart(ctx echo.Context, contestId int, tourNumber int) error {
	return h.adminTimetableHandle.PostAdminTimetableTourStart(ctx, contestId, tourNumber)
}
