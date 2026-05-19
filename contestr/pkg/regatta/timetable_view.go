package regatta

import "time"

type SaveTimetableRequest struct {
	PendingSlots     []ScheduleSlot `json:"pending_slots"`
	AutoStartEnabled bool           `json:"auto_start_enabled"`
}

type UpdateActiveTourDurationRequest struct {
	Duration int `json:"duration"`
}

type TimelineSegment struct {
	Sequence     *int   `json:"sequence,omitempty"`
	PendingIndex *int   `json:"pending_index,omitempty"`
	Kind         string `json:"kind"`
	Round        *int   `json:"round,omitempty"`
	Duration     int    `json:"duration"`
	StartTime    int    `json:"start_time"`
	Status       string `json:"status"`
	Editable     bool   `json:"editable"`
}

type TimetableView struct {
	ContestId           int               `json:"contest_id"`
	ContestStartTime    *time.Time        `json:"contest_start_time,omitempty"`
	ServerNow           time.Time         `json:"server_now"`
	ElapsedSeconds      int               `json:"elapsed_seconds"`
	TimelineSegments    []TimelineSegment `json:"timeline_segments"`
	PendingSlots        []ScheduleSlot    `json:"pending_slots"`
	NextTourNumber      *int              `json:"next_tour_number"`
	AutoStartEnabled    bool              `json:"auto_start_enabled"`
	AutoStartAvailable  bool              `json:"auto_start_available"`
}
