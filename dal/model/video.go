package model

type Video struct {
	ID            int64  `gorm:"primarykey"`
	AuthorID      int64  `gorm:"index:idx_authorid;not null" json:"author_id,omitempty"`
	PlayUrl       string `gorm:"type:varchar(255);not null" json:"play_url,omitempty"`
	CoverUrl      string `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount int64  `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	CommentCount  int64  `gorm:"default:0;not null" json:"comment_count,omitempty"`
	Title         string `gorm:"type:varchar(255);not null" json:"title,omitempty"`
	PublishTime   int64  `gorm:"not null;index:idx_publishtime" json:"publish_time,omitempty"`
}
