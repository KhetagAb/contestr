package handlers

import (
	"contestr/internal/handlers/admin"
	"contestr/internal/handlers/codeforces"
	"contestr/internal/handlers/contests"
	"contestr/internal/handlers/regatta"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	contestHandle        *codeforces.ContestHandle
	regattaContestHandle *regatta.ContestHandle
	contestsListHandle   *contests.ListHandle
	adminLoginHandle     *admin.LoginHandle
	adminMeHandle        *admin.MeHandle
	adminTimetableHandle *admin.TimetableHandle
	adminContestsHandle  *admin.ContestsHandle
}

func (h *Handlers) GetRegattaContestStandings(ctx echo.Context, contestId int) error {
	return h.regattaContestHandle.GetContest(ctx, contestId)
}

func NewHandlers(
	contestHandle *codeforces.ContestHandle,
	regattaContestHandle *regatta.ContestHandle,
	contestsListHandle *contests.ListHandle,
	adminLoginHandle *admin.LoginHandle,
	adminMeHandle *admin.MeHandle,
	adminTimetableHandle *admin.TimetableHandle,
	adminContestsHandle *admin.ContestsHandle,
) *Handlers {
	return &Handlers{
		contestHandle:        contestHandle,
		regattaContestHandle: regattaContestHandle,
		contestsListHandle:   contestsListHandle,
		adminLoginHandle:     adminLoginHandle,
		adminMeHandle:        adminMeHandle,
		adminTimetableHandle: adminTimetableHandle,
		adminContestsHandle:  adminContestsHandle,
	}
}

func (h *Handlers) GetContest(ctx echo.Context, contestId int) error {
	return h.contestHandle.GetContest(ctx, contestId)
}

func (h *Handlers) GetContests(ctx echo.Context) error {
	return h.contestsListHandle.GetContests(ctx)
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

func (h *Handlers) PatchAdminTimetableActiveTourDuration(ctx echo.Context, contestId int) error {
	return h.adminTimetableHandle.PatchAdminTimetableActiveTourDuration(ctx, contestId)
}

func (h *Handlers) PostAdminTimetableAdvance(ctx echo.Context, contestId int) error {
	return h.adminTimetableHandle.PostAdminTimetableAdvance(ctx, contestId)
}

func (h *Handlers) GetAdminContests(ctx echo.Context) error {
	return h.adminContestsHandle.GetAdminContests(ctx)
}

func (h *Handlers) PostAdminContest(ctx echo.Context) error {
	return h.adminContestsHandle.PostAdminContest(ctx)
}

func (h *Handlers) DeleteAdminContest(ctx echo.Context, contestId int) error {
	return h.adminContestsHandle.DeleteAdminContest(ctx, contestId)
}

func (h *Handlers) PatchAdminContestSettings(ctx echo.Context, contestId int) error {
	return h.adminContestsHandle.PatchAdminContestSettings(ctx, contestId)
}

func (h *Handlers) PostAdminContestRefresh(ctx echo.Context, contestId int) error {
	return h.adminContestsHandle.PostAdminContestRefresh(ctx, contestId)
}

func (h *Handlers) GetAdminContestHandles(ctx echo.Context, contestId int) error {
	return h.adminContestsHandle.GetAdminContestHandles(ctx, contestId)
}

func (h *Handlers) PutAdminContestHandles(ctx echo.Context, contestId int) error {
	return h.adminContestsHandle.PutAdminContestHandles(ctx, contestId)
}

func (h *Handlers) DeleteAdminContestHandle(ctx echo.Context, contestId int, handle string) error {
	return h.adminContestsHandle.DeleteAdminContestHandle(ctx, contestId, handle)
}
