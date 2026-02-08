package repository

import (
	"database/sql"
	"socialnet/internal/model"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(report *model.Report) (int64, error) {
	query := `INSERT INTO reports (reporter_id, target_type, target_id, reason) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, report.ReporterID, report.TargetType, report.TargetID, report.Reason)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ReportRepository) GetAll(status model.ReportStatus, limit int) ([]*model.Report, error) {
	query := `SELECT id, reporter_id, target_type, target_id, reason, status, created_at 
			  FROM reports WHERE status = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.Query(query, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

func (r *ReportRepository) GetByID(id int64) (*model.Report, error) {
	query := `SELECT id, reporter_id, target_type, target_id, reason, status, created_at FROM reports WHERE id = ?`
	report := &model.Report{}
	err := r.db.QueryRow(query, id).Scan(
		&report.ID, &report.ReporterID, &report.TargetType, &report.TargetID,
		&report.Reason, &report.Status, &report.CreatedAt,
	)
	return report, err
}

func (r *ReportRepository) UpdateStatus(id int64, status model.ReportStatus) error {
	query := `UPDATE reports SET status = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *ReportRepository) scanReports(rows *sql.Rows) ([]*model.Report, error) {
	var reports []*model.Report
	for rows.Next() {
		report := &model.Report{}
		err := rows.Scan(&report.ID, &report.ReporterID, &report.TargetType,
			&report.TargetID, &report.Reason, &report.Status, &report.CreatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, rows.Err()
}
