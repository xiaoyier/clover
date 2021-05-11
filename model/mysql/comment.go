package mysql

import "time"

func init() {
	RegisterAutoMigrates(&Comment{})
}

type CommentCreateReq struct {
	PostID  int64  `json:"post_id,string" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type Comment struct {
	ID        uint      `gorm:"primary_key;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(6)" json:"created_at"`
	CommentID int64     `gorm:"unique_index;column:comment_id;type:bigint(20)" json:"comment_id"`
	PostID    int64     `gorm:"column:post_id;type:bigint(20)" json:"post_id"`
	AuthorID  int64     `gorm:"column:author_id;type:bigint(20)" json:"author_id"`
	Content   string    `gorm:"column:content;type:varchar(200)" json:"content"`
}

func (c *Comment) TableName() string {
	return "comment"
}

func (c *Comment) Insert() error {
	d := db.Save(c)
	return d.Error
}

func QueryCommentList(postId int64) ([]Comment, error) {
	comments := new([]Comment)
	d := db.Where(Comment{PostID: postId}).Limit(5).Order("created_at desc").Find(&comments)
	if d.RecordNotFound() {
		return nil, nil
	}

	return *comments, d.Error
}
