package dal

import (
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

func CreateComment(tx *gorm.DB, VideoId int64, UserId int64, CommentText string) (*Comment, error) {
	comment := Comment{
		VideoID: VideoId,
		UserID:  UserId,
		Content: CommentText,
	}
	err := tx.Create(&comment).Error
	return &comment, err
}

func GetCommentsByVideoId(tx *gorm.DB, VideoId int64) ([]model.Comment, error) {
	var commentsDal []Comment
	err := tx.Model(&Comment{}).Order("create_time desc").Where("video_id=?", VideoId).Find(&commentsDal).Error
	comments := make([]model.Comment, len(commentsDal))
	for i, row := range commentsDal {
		user := FindUserById(mysql.DB, row.UserID)
		comments[i] = model.Comment{
			Id:         row.Id,
			User:       user,
			Content:    row.Content,
			CreateDate: row.CreateTime.Format("01-02"),
		}
	}
	return comments, err
}

func DeleteCommentByCommentId(tx *gorm.DB, CommentId int64) error {
	var comment Comment
	err := tx.Delete(Comment{Id: CommentId}).Delete(&comment).Error
	return err
}
