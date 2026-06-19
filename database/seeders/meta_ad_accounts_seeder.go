package seeders

import (
	"github.com/alex/ads_backend/internal/meta/ad_account"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedMetaAdAccounts(db *gorm.DB) {
	accounts := []ad_account.MetaAdAccount{
		{ID: "act_1001304461010256", Name: "EVPO - CPAS - Stanley HKU - 12", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("754141426281556"), BusinessName: ptrString("EVPO External Imers 11"), IsActive: true},
		{ID: "act_1238860926711356", Name: "TZL - STANLEY DR O - 14", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("606063277489753"), BusinessName: ptrString("BM - SMB 1"), IsActive: true},
		{ID: "act_1333208877401665", Name: "MKH - STANLEY DR O (CPAS) - 53", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("103721901891307"), BusinessName: ptrString("BM - SMB 2"), IsActive: true},
		{ID: "act_1377261246549191", Name: "EVPO - CPAS - Stanley - 3", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("398909192302631"), BusinessName: ptrString("EVPO External Imers 8"), IsActive: true},
		{ID: "act_1402109817316002", Name: "MKH - STANLEY DR O (CPAS) - 48", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("103721901891307"), BusinessName: ptrString("BM - SMB 2"), IsActive: true},
		{ID: "act_1432455930649795", Name: "EVPO - CPAS - Stanley HKU - 11", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("6010201902403484"), BusinessName: ptrString("EVPO External Imers 10"), IsActive: true},
		{ID: "act_1696282174202730", Name: "MKH - STANLEY DR O (CPAS) - 50", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("103721901891307"), BusinessName: ptrString("BM - SMB 2"), IsActive: true},
		{ID: "act_181460128108745", Name: "TZL - STANLEY DR O - 07", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("226184613165542"), BusinessName: ptrString("BM - SMB 4"), IsActive: true},
		{ID: "act_1928676774142204", Name: "TZL - STANLEY DR O - 12", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("606063277489753"), BusinessName: ptrString("BM - SMB 1"), IsActive: true},
		{ID: "act_316869327657518", Name: "EVPO - CPAS - Stanley - 2", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("398909192302631"), BusinessName: ptrString("EVPO External Imers 8"), IsActive: true},
		{ID: "act_3507779102795118", Name: "EVPO - CPAS - Stanley HKU - 8", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("6010201902403484"), BusinessName: ptrString("EVPO External Imers 10"), IsActive: true},
		{ID: "act_365127726027187", Name: "EVPO - BM - Stanley HKU - 55", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("389290890382541"), BusinessName: ptrString("SM External Imers 5"), IsActive: true},
		{ID: "act_492332263095393", Name: "TZL - STANLEY DR O - 11", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("606063277489753"), BusinessName: ptrString("BM - SMB 1"), IsActive: true},
		{ID: "act_541050504790549", Name: "EVPO - BM - Stanley HKU - 48", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("754141426281556"), BusinessName: ptrString("EVPO External Imers 11"), IsActive: true},
		{ID: "act_6898666400179420", Name: "MKH - STANLEY DR O (CPAS) - 47", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("103721901891307"), BusinessName: ptrString("BM - SMB 2"), IsActive: true},
		{ID: "act_770786938140607", Name: "EVPO - BM - Stanley HKU - 42", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("3127037640939943"), BusinessName: ptrString("EVPO External Imers 7"), IsActive: true},
		{ID: "act_828773662282541", Name: "EVPO - CPAS - Stanley - 1", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("398909192302631"), BusinessName: ptrString("EVPO External Imers 8"), IsActive: true},
		{ID: "act_847465772995635", Name: "TZL - STANLEY DR O - 20", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("606063277489753"), BusinessName: ptrString("BM - SMB 1"), IsActive: true},
		{ID: "act_902301467650788", Name: "TZL - STANLEY DR O - 03", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("226184613165542"), BusinessName: ptrString("BM - SMB 4"), IsActive: true},
		{ID: "act_9080264445377165", Name: "TZL - STANLEY DR O - 04", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("226184613165542"), BusinessName: ptrString("BM - SMB 4"), IsActive: true},
		{ID: "act_999522741134376", Name: "EVPO - BM - Stanley HKU - 58", AccountStatus: 1, Currency: ptrString("IDR"), TimezoneName: ptrString("Asia/Jakarta"), BusinessID: ptrString("389290890382541"), BusinessName: ptrString("SM External Imers 5"), IsActive: true},
	}

	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "account_status", "currency", "timezone_name", "business_id", "business_name", "is_active"}),
	}).Create(&accounts)
}

func ptrString(s string) *string {
	return &s
}
