package dto

type BrandDashboardFilter struct {
	Search    string
	Page      int
	Limit     int
	BrandIDs  []uint64
	DateStart string
	DateStop  string
}
