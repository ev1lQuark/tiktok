package pattern

import "fmt"

var (
	LikeMapDataKey          = "LIKE:MAP:LIKE_DATA"
	LikeMapUserIdCountKey   = "LIKE:MAP:USERID_COUNT"
	LikeMapVideoIdCountKey  = "LIKE:MAP:VIDEOID_COUNT"
	LikeMapAuthorIdCountKey = "LIKE:MAP:AUTHORID_COUNT"
	LikeSetUserIdKey        = "LIKE:SET:USERID:%d"
)

func GetLikeMapDataKey(userId, videoId int64) string {
	return fmt.Sprintf("%d:%d", userId, videoId)
}

func ParseLikeMapDataKey(key string) (userId, videoId int64) {
	fmt.Sscanf(key, "%d:%d", &userId, &videoId)
	return
}

func GetLikeSetUserIdKey(userId int64) string {
	return fmt.Sprintf(LikeSetUserIdKey, userId)
}
