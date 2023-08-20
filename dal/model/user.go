package model

type User struct {
	ID              int64  `gorm:"index:idx_userid,primarykey"`
	UserName        string `gorm:"index:idx_username,unique;type:varchar(32);not null" json:"name,omitempty"`
	Password        string `gorm:"type:varchar(32);not null" json:"password,omitempty"`
	FollowCount     int64  `gorm:"default:0;not null" json:"follow_count,omitempty"`                                                           // 关注总数
	FollowerCount   int64  `gorm:"default:0;not null" json:"follower_count,omitempty"`                                                         // 粉丝总数
	Avatar          string `gorm:"type:varchar(255)" json:"avatar,omitempty"`                                                                  // 用户头像
	BackgroundImage string `gorm:"column:background_image;type:varchar(256);default:default_background.jpg" json:"background_image,omitempty"` // 用户个人页顶部大图
	WorkCount       int64  `gorm:"default:0;not null" json:"work_count,omitempty"`                                                             // 作品数
	FavoriteCount   int64  `gorm:"default:0;not null" json:"favorite_count,omitempty"`                                                         // 喜欢数
	TotalFavorited  int64  `gorm:"default:0;not null" json:"total_favorited,omitempty"`                                                        // 获赞总量
	Signature       string `gorm:"type:varchar(255)" json:"signature,omitempty"`
}
