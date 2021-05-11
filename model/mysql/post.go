package mysql

import (
	"time"
)

func init() {
	RegisterAutoMigrates(&Post{})
}

type PostListPage struct {
	PageNumber int `form:"page_number" binding:"gte=0"`
	PageSize   int `form:"page_size" binding:"gte=1"`
}

type PostCreateReq struct {
	Title       string `json:"title" binding:"required,min=2,max=200"`
	Content     string `json:"content" binding:"required,min=2"`
	CommunityID int64  `json:"community_id,string" binding:"required"`
}

type PostItem struct {
	PostID        int64  `json:"post_id"`
	AuthorName    string `json:"author_name"`
	CommunityName string `json:"community_name"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Status        int8   `json:"status"`
}

type Post struct {
	ID          uint      `gorm:"primary_key;column:id" json:"id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime(6)" json:"updated_at"`
	PostID      int64     `gorm:"unique;column:post_id;type:bigint(20)" json:"post_id"`
	AuthorID    int64     `gorm:"unique_index:author_title;column:author_id;type:bigint(20)" json:"author_id"`
	CommunityID int64     `gorm:"column:community_id;type:bigint(20)" json:"community_id"`
	Status      int8      `gorm:"column:status;type:tinyint(4);default:1" json:"status"`
	Title       string    `gorm:"unique_index:author_title;column:title;type:varchar(200)" json:"title"`
	Content     string    `gorm:"column:content;type:longText" json:"content"`
}

func (p *Post) TableName() string {
	return "post"
}

func (p *Post) Insert() error {

	d := db.Save(p)
	return d.Error
}

func QueryPostList(page *PostListPage) ([]Post, error) {

	items := make([]Post, 0)
	d := db.Model(&Post{}).
		Limit(page.PageSize).
		Offset(page.PageNumber * page.PageSize).
		Order("created_at desc").
		Find(&items)
	if d.RecordNotFound() {
		return nil, d.Error
	}

	return items, d.Error
}

func QueryPost(id int64) (*Post, error) {

	post := new(Post)
	d := db.Where(Post{PostID: id}).Find(&post)
	if d.RecordNotFound() {
		return nil, nil
	}
	return post, d.Error
}
