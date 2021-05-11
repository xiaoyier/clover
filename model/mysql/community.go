package mysql

import (
	"time"
)

func init() {
	RegisterAutoMigrates(&Community{})
}

type CommunityCreateReq struct {
	CommunityName string `json:"community_name" binding:"required,min=2,max=30"`
	Introduction  string `json:"introduction"`
}

type CommunityItem struct {
	CommunityID   int64  `json:"community_id"`
	CommunityName string `json:"community_name"`
}

type Community struct {
	ID            uint      `gorm:"primary_key;column:id" json:"id"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
	CommunityID   int64     `gorm:"unique_index;column:community_id;type:bigint(20)" json:"community_id"`
	CommunityName string    `gorm:"unique_index;column:community_name;type:varchar(30)" json:"community_name"`
	Introduction  string    `gorm:"column:introduction;type:longText" json:"introduction"`
}

func (community *Community) TableName() string {
	return "community"
}

func QueryCommunity(communityName string) (*Community, error) {
	com := new(Community)
	d := db.Where("community_name = ?", communityName).First(&com)
	if d.RecordNotFound() {
		return nil, nil
	}

	return com, d.Error
}

func QueryAllCommunities() ([]CommunityItem, error) {

	result := make([]CommunityItem, 0)
	d := db.Model(&Community{}).Scan(&result)
	if d.RecordNotFound() {
		return nil, nil
	}

	return result, d.Error
}

func QueryCommunityByID(id int64) (*Community, error) {
	community := new(Community)
	d := db.Where(Community{CommunityID: id}).First(&community)
	if d.RecordNotFound() {
		return nil, nil
	}

	return community, d.Error
}

func (community *Community) Insert() error {

	d := db.Save(community)
	return d.Error
}
