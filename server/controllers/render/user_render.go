package render

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	"strconv"
	"strings"

	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/date"
)

func BuildUserInfoDefaultIfNull(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.Id = id
		user.Username = simple.SqlNullString(strconv.FormatInt(id, 10))
		user.Nickname = "匿名用户" + strconv.FormatInt(id, 10)
		user.CreateTime = date.NowTimestamp()
	}
	return BuildUserInfo(user)
}

func BuildUserInfo(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}
	ret := &model.UserInfo{
		Id:           user.Id,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		SmallAvatar:  HandleOssImageStyleAvatar(user.Avatar),
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		FansCount:    user.FansCount,
		FollowCount:  user.FollowCount,
		Score:        user.Score,
		Description:  user.Description,
		CreateTime:   user.CreateTime,
	}
	if len(ret.Description) == 0 {
		ret.Description = "这家伙很懒，什么都没留下"
	}
	if user.Status == constants.StatusDeleted {
		ret.Nickname = "黑名单用户"
		ret.Description = ""
	}
	return ret
}

func BuildUserDetail(user *model.User) *model.UserDetail {
	if user == nil {
		return nil
	}
	ret := &model.UserDetail{
		UserInfo:             *BuildUserInfo(user),
		Username:             user.Username.String,
		BackgroundImage:      user.BackgroundImage,
		SmallBackgroundImage: HandleOssImageStyleSmall(user.BackgroundImage),
		Type:                 user.Type,
		HomePage:             user.HomePage,
		Forbidden:            user.IsForbidden(),
		Status:               user.Status,
	}
	if user.Status == constants.StatusDeleted {
		ret.Username = "blacklist"
		ret.HomePage = ""
		ret.Score = 0
		ret.Forbidden = true
	}
	return ret
}

func BuildUserProfile(user *model.User) *model.UserProfile {
	if user == nil {
		return nil
	}
	roles := strings.Split(user.Roles, ",")
	ret := &model.UserProfile{
		UserDetail:    *BuildUserDetail(user),
		Roles:         roles,
		Email:         user.Email.String,
		EmailVerified: user.EmailVerified,
		PasswordSet:   len(user.Password) > 0,
	}
	return ret
}
