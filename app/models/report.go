package models

type DashboardStatistics struct {
	TotalPrestasi int64            `json:"total_achievements"`
	ByStatus      map[string]int64 `json:"status_breakdown"`
}

type StudentReportProfile struct {
	StudentID   string `json:"student_id"`
	NIM         string `json:"nim"`
	FullName    string `json:"full_name"`
	TotalPoints int    `json:"total_points"`
	TotalItems  int    `json:"total_items"`
}

type StudentReportResponse struct {
	Profile      StudentReportProfile  `json:"profile"`
	Achievements []AchievementResponse `json:"achievements"`
}