package setting

import "time"

const UserIdPenetrationKey = "like:userid_penetration"

var (
	UserIdKeyPattern   = "like:userid:%d"
	UserIdValuePattern = "videoid:%d,authorid:%d"
	UserIdExpire       = time.Duration(24 * time.Hour)

	VideoIdKeyPattern = "like:videoid:%d"
	VideoIdExpire     = time.Duration(24 * time.Hour)
)
