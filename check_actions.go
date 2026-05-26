package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MetaInsight struct {
	CampaignID        string
	Actions           []byte
	CostPerActionType []byte
}

func main() {
	dsn := "root:@tcp(localhost:3306)/ads_backend?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var insights []MetaInsight
	db.Table("meta_insights").Select("campaign_id, actions, cost_per_action_type").Where("actions IS NOT NULL").Limit(5).Find(&insights)

	for _, i := range insights {
		fmt.Printf("CampaignID: %s\n", i.CampaignID)
		fmt.Printf("Actions: %s\n", string(i.Actions))
		fmt.Printf("Cost: %s\n", string(i.CostPerActionType))
		fmt.Println("-----------------")
	}
}
