package jobs

import (
	"log"
	"time"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/meta"
)

type MetaAdsSyncJob struct {
	metaService meta.Service
}

func NewMetaAdsSyncJob(metaService meta.Service) *MetaAdsSyncJob {
	return &MetaAdsSyncJob{metaService: metaService}
}

// Start launches a ticker-based job that runs every 15 minutes
func (j *MetaAdsSyncJob) Start() {
	log.Println("Meta Ads Sync Job scheduler initialized (Interval: 15 minutes)")
	
	// Run immediately on start in background
	go j.Run()

	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for range ticker.C {
			j.Run()
		}
	}()
}

// Run performs the actual synchronization logic
func (j *MetaAdsSyncJob) Run() {
	log.Println("Starting Meta Ads synchronization job...")
	
	adAccountID := config.MetaAdAccountID
	if adAccountID == "" {
		log.Println("Warning: META_AD_ACCOUNT_ID is empty, skipping background sync")
		return
	}

	campaigns, err := j.metaService.GetCampaigns(adAccountID)
	if err != nil {
		log.Printf("Error syncing Meta campaigns: %v", err)
		return
	}

	log.Printf("Successfully synchronized %d campaigns for ad account %s", len(campaigns), adAccountID)
}
