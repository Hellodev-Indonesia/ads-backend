package dashboard

import (
	"context"

	"github.com/alex/ads_backend/internal/reporting/dashboard/dto"
)

type DashboardService interface {
	GetSummary(ctx context.Context) (*dto.SummaryResponse, error)
	GetBrandsMonitoring(ctx context.Context) ([]dto.BrandMonitoringItem, error)
	GetRecentActivities(ctx context.Context) ([]dto.RecentActivityItem, error)
}

type dashboardService struct {
	repo DashboardRepository
}

func NewDashboardService(repo DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetSummary(ctx context.Context) (*dto.SummaryResponse, error) {
	totalSpend, totalOngoingAds, securityIssues, err := s.repo.GetSummaryMetrics(ctx)
	if err != nil {
		return nil, err
	}

	status := "No suspicious activity detected"
	if securityIssues > 0 {
		status = "Suspicious activity detected"
	}

	return &dto.SummaryResponse{
		TotalSpend:      totalSpend,
		TotalOngoingAds: totalOngoingAds,
		SecurityStatus:  status,
		SecurityIssues:  securityIssues,
	}, nil
}

func (s *dashboardService) GetBrandsMonitoring(ctx context.Context) ([]dto.BrandMonitoringItem, error) {
	rawItems, err := s.repo.GetBrandsMonitoring(ctx)
	if err != nil {
		return nil, err
	}

	var items []dto.BrandMonitoringItem
	for _, raw := range rawItems {
		items = append(items, dto.BrandMonitoringItem{
			BrandID:        raw.BrandID,
			BrandName:      raw.BrandName,
			AdAccountCount: raw.AdAccountCount,
			ActiveCampaign: raw.ActiveCampaign,
			TotalSpends:    raw.TotalSpends,
		})
	}

	return items, nil
}

func (s *dashboardService) GetRecentActivities(ctx context.Context) ([]dto.RecentActivityItem, error) {
	rawItems, err := s.repo.GetRecentActivities(ctx)
	if err != nil {
		return nil, err
	}

	var items []dto.RecentActivityItem
	for _, raw := range rawItems {
		dateStr := ""
		if raw.EventTime != nil {
			dateStr = raw.EventTime.Format("02 Jan at 15.04") // Format similar to design "03 April at 15.40"
		}

		items = append(items, dto.RecentActivityItem{
			ActivityID:      raw.ActivityID,
			BrandName:       raw.BrandName,
			AdAccountName:   raw.AdAccountName,
			Activity:        raw.Activity,
			ActivityDetails: raw.ActivityDetails,
			ItemChanged:     raw.ItemChanged,
			ChangeBy:        raw.ChangeBy,
			DateAndTime:     dateStr,
		})
	}

	return items, nil
}
