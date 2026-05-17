package tgbot

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"contestr/internal/configs"
	regattasvc "contestr/internal/services/regatta"
	"contestr/pkg/logger"
	regattapkg "contestr/pkg/regatta"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TimetableService interface {
	CreateTimetable(ctx context.Context, timetable regattapkg.ToursTimetable) (*regattapkg.ToursTimetable, error)
	GetTimetable(ctx context.Context, contestID int) (*regattapkg.ToursTimetable, error)
	UpdateTimetable(ctx context.Context, timetable regattapkg.ToursTimetable) (*regattapkg.ToursTimetable, error)
	DeleteTimetable(ctx context.Context, contestID int) error
	MoveTimetableTour(ctx context.Context, contestID int, tourNumber int, newStartTime int) (*regattapkg.ToursTimetable, error)
	GetFirstNotStartedTimetableTour(ctx context.Context, contestID int) (*regattapkg.TourConfig, error)
}

type TimetableHandle struct {
	cfg       *configs.Config
	timetable TimetableService
}

func NewTimetableHandle(
	cfg *configs.Config,
	timetable TimetableService,
) *TimetableHandle {
	return &TimetableHandle{
		cfg:       cfg,
		timetable: timetable,
	}
}

func (h *TimetableHandle) Register() (bot.HandlerType, string, bot.MatchType, bot.HandlerFunc) {
	return bot.HandlerTypeMessageText, "timetable", bot.MatchTypeCommand, h.Handle
}

func (h *TimetableHandle) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	if !h.isAdmin(update.Message.From.ID) {
		h.sendText(ctx, b, chatID, "Only admins can edit tour timetables.")
		return
	}

	parts := strings.Fields(update.Message.Text)
	if len(parts) < 2 {
		h.sendText(ctx, b, chatID, timetableUsage())
		return
	}

	message, err := h.handleCommand(ctx, parts)
	if err != nil {
		logger.Errorf(ctx, "error handling timetable command: %v", err)
		message = formatTimetableError(err)
	}
	h.sendText(ctx, b, chatID, message)
}

func (h *TimetableHandle) handleCommand(ctx context.Context, parts []string) (string, error) {
	switch parts[1] {
	case "create":
		return h.handleCreate(ctx, parts)
	case "get":
		return h.handleGet(ctx, parts)
	case "update":
		return h.handleUpdate(ctx, parts)
	case "delete":
		return h.handleDelete(ctx, parts)
	case "move":
		return h.handleMove(ctx, parts)
	case "next":
		return h.handleNext(ctx, parts)
	default:
		return timetableUsage(), nil
	}
}

func (h *TimetableHandle) handleCreate(ctx context.Context, parts []string) (string, error) {
	if len(parts) < 4 {
		return timetableUsage(), nil
	}

	contestID, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid contest_id: %w", err)
	}

	tours, err := parseTourConfigs(parts[3:])
	if err != nil {
		return "", err
	}

	timetable, err := h.timetable.CreateTimetable(ctx, regattapkg.ToursTimetable{
		ContestId: contestID,
		TourTimes: tours,
	})
	if err != nil {
		return "", err
	}

	return "Расписание создано.\n" + formatTimetable(timetable), nil
}

func (h *TimetableHandle) handleGet(ctx context.Context, parts []string) (string, error) {
	if len(parts) != 3 {
		return timetableUsage(), nil
	}

	contestID, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid contest_id: %w", err)
	}

	timetable, err := h.timetable.GetTimetable(ctx, contestID)
	if err != nil {
		return "", err
	}

	return formatTimetable(timetable), nil
}

func (h *TimetableHandle) handleUpdate(ctx context.Context, parts []string) (string, error) {
	if len(parts) < 4 {
		return timetableUsage(), nil
	}

	contestID, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid contest_id: %w", err)
	}

	tours, err := parseTourConfigs(parts[3:])
	if err != nil {
		return "", err
	}

	timetable, err := h.timetable.UpdateTimetable(ctx, regattapkg.ToursTimetable{
		ContestId: contestID,
		TourTimes: tours,
	})
	if err != nil {
		return "", err
	}

	return "Расписание обновлено.\n" + formatTimetable(timetable), nil
}

func (h *TimetableHandle) handleDelete(ctx context.Context, parts []string) (string, error) {
	if len(parts) != 3 {
		return timetableUsage(), nil
	}

	contestID, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid contest_id: %w", err)
	}

	if err := h.timetable.DeleteTimetable(ctx, contestID); err != nil {
		return "", err
	}

	return fmt.Sprintf("Расписание для контеста %d удалено.", contestID), nil
}

func (h *TimetableHandle) handleMove(ctx context.Context, parts []string) (string, error) {
	if len(parts) != 5 {
		return timetableUsage(), nil
	}

	contestID, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid contest_id: %w", err)
	}
	tourNumber, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", fmt.Errorf("invalid tour_number: %w", err)
	}
	newStartTime, err := strconv.Atoi(parts[4])
	if err != nil {
		return "", fmt.Errorf("invalid start_time: %w", err)
	}

	timetable, err := h.timetable.MoveTimetableTour(ctx, contestID, tourNumber, newStartTime)
	if err != nil {
		return "", err
	}

	return "Расписание перенесено.\n" + formatTimetable(timetable), nil
}

func (h *TimetableHandle) handleNext(ctx context.Context, parts []string) (string, error) {
	if len(parts) != 3 {
		return timetableUsage(), nil
	}

	contestID, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid contest_id: %w", err)
	}

	tour, err := h.timetable.GetFirstNotStartedTimetableTour(ctx, contestID)
	if err != nil {
		return "", err
	}

	return "Первый незапущенный тур:\n" + formatTour(*tour), nil
}

func parseTourConfigs(rawTours []string) ([]regattapkg.TourConfig, error) {
	tours := make([]regattapkg.TourConfig, 0, len(rawTours))
	for _, rawTour := range rawTours {
		parts := strings.Split(rawTour, ":")
		if len(parts) < 2 || len(parts) > 3 {
			return nil, fmt.Errorf("invalid tour config %q, expected start:duration[:started]", rawTour)
		}

		startTime, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid start_time in %q: %w", rawTour, err)
		}
		duration, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid duration in %q: %w", rawTour, err)
		}

		started := false
		if len(parts) == 3 {
			started, err = parseStarted(parts[2])
			if err != nil {
				return nil, err
			}
		}

		tours = append(tours, regattapkg.TourConfig{
			StartTime: startTime,
			Duration:  duration,
			Started:   started,
		})
	}
	return tours, nil
}

func parseStarted(value string) (bool, error) {
	switch strings.ToLower(value) {
	case "true", "1", "yes", "started":
		return true, nil
	case "false", "0", "no", "not_started", "pending":
		return false, nil
	default:
		return false, fmt.Errorf("invalid started value %q", value)
	}
}

func formatTimetable(timetable *regattapkg.ToursTimetable) string {
	lines := []string{
		fmt.Sprintf("Расписание контеста %d:", timetable.ContestId),
	}
	if len(timetable.TourTimes) == 0 {
		lines = append(lines, "Туры не настроены.")
		return strings.Join(lines, "\n")
	}

	for i, tour := range timetable.TourTimes {
		lines = append(lines, formatTourConfig(i+1, tour))
	}
	return strings.Join(lines, "\n")
}

func formatTourConfig(index int, tour regattapkg.TourConfig) string {
	return fmt.Sprintf("%d. %s", index, formatTour(tour))
}

func formatTour(tour regattapkg.TourConfig) string {
	return fmt.Sprintf("старт=%dс длительность=%dс запущен=%t", tour.StartTime, tour.Duration, tour.Started)
}

func formatTimetableError(err error) string {
	switch {
	case errors.Is(err, regattasvc.ErrInvalidTimetable):
		return "Invalid timetable: " + err.Error()
	case errors.Is(err, regattasvc.ErrTimetableAlreadyExists):
		return "Timetable already exists. Use /timetable update to replace it."
	case errors.Is(err, regattasvc.ErrTimetableNotFound):
		return "Timetable not found."
	case errors.Is(err, regattasvc.ErrTourNotFound):
		return "Tour not found."
	case errors.Is(err, regattasvc.ErrTourAlreadyStarted):
		return "Tour already started."
	default:
		return "Timetable command failed: " + err.Error()
	}
}

func timetableUsage() string {
	return strings.Join([]string{
		"Использование:",
		"/timetable create <contest_id> <start:duration[:started]> ...",
		"/timetable get <contest_id>",
		"/timetable update <contest_id> <start:duration[:started]> ...",
		"/timetable delete <contest_id>",
		"/timetable move <contest_id> <tour_number> <new_start_time>",
		"/timetable next <contest_id>",
		"",
		"Пример:",
		"/timetable create 613423 0:1800 2400:1800:false 4800:1800",
	}, "\n")
}

func (h *TimetableHandle) isAdmin(userID int64) bool {
	for _, adminID := range h.cfg.Telegram.Admins {
		if int64(adminID) == userID {
			return true
		}
	}
	return false
}

func (h *TimetableHandle) sendText(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		logger.Errorf(ctx, "error sending timetable message: %v", err)
	}
}
