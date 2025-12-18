package entity

type PostLikesCount struct {
	PostID    int64  `gorm:"column:post_id;primaryKey"`
	LikeCount uint64 `gorm:"column:like_count"`
}
