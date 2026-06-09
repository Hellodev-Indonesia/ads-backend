package ad_account_test

import (
	"testing"

	"github.com/alex/ads_backend/internal/meta/ad_account"
	"github.com/alex/ads_backend/internal/meta/ad_account/dto"
	"github.com/stretchr/testify/assert"
)

func setupService(t *testing.T) (*ad_account.MockRepository, ad_account.Service) {
	mockRepo := ad_account.NewMockRepository(t)
	// We can pass nil for meta_client.Client since we are not testing SyncAdAccounts
	svc := ad_account.NewService(nil, mockRepo)
	return mockRepo, svc
}

func TestService_GetAdAccounts(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := ad_account.AdAccountFilter{
		Page:  1,
		Limit: 10,
	}

	accounts := []ad_account.MetaAdAccount{
		{ID: "act_1", Name: "Account 1", AccountStatus: 1},
	}

	mockRepo.On("FindAll", filter).Return(accounts, int64(1), nil)

	resp, meta, err := svc.GetAdAccounts(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "act_1", resp[0].ID)
	assert.NotNil(t, meta)
	assert.Equal(t, int64(1), meta.Total)
}

func TestService_GetUnassigned(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := ad_account.AdAccountFilter{
		Page:  1,
		Limit: 10,
	}

	accounts := []ad_account.MetaAdAccount{
		{ID: "act_2", Name: "Unassigned Account", AccountStatus: 1},
	}

	mockRepo.On("FindUnassigned", filter).Return(accounts, int64(1), nil)

	resp, meta, err := svc.GetUnassigned(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "act_2", resp[0].ID)
	assert.NotNil(t, meta)
	assert.Equal(t, int64(1), meta.Total)
}

func TestService_BulkAssignBrand(t *testing.T) {
	mockRepo, svc := setupService(t)

	ids := []string{"act_1", "act_2"}
	brandID := uint64(1)

	req := dto.AssignBrandRequest{
		AdAccountIDs: ids,
		BrandID:      &brandID,
	}

	mockRepo.On("UpdateBrandIDBatch", ids, &brandID).Return(nil)

	err := svc.BulkAssignBrand(req)

	assert.NoError(t, err)
}

func TestService_BulkAssignBrandByBusiness(t *testing.T) {
	mockRepo, svc := setupService(t)

	brandID := uint64(1)
	businessID := "bus_1"

	req := dto.AssignBrandRequest{
		BusinessID:   &businessID,
		BrandID:      &brandID,
	}

	mockRepo.On("UpdateBrandIDByBusiness", businessID, &brandID).Return(nil)

	err := svc.BulkAssignBrand(req)

	assert.NoError(t, err)
}
