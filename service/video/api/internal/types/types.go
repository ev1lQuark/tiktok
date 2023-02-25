// Code generated by goctl. DO NOT EDIT.
package types

type FeedReq struct {
	LatestTime string `form:"latest_time,optional"`
	Token      string `form:"token,optional"`
}

type FeedReply struct {
	StatusCode int         `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	NextTime   int         `json:"next_time"`
	VideoList  []VideoList `json:"video_list"`
}

type Author struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	FollowCount     int    `json:"follow_count"`
	FollowerCount   int    `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  string `json:"total_favorited"`
	WorkCount       int    `json:"work_count"`
	FavoriteCount   int    `json:"favorite_count"`
}

type VideoList struct {
	ID            int    `json:"id"`
	Author        Author `json:"author"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type PublishActionReq struct {
	Token string `form:"token"`
	Title string `form:"title"`
}

type PublishActionReply struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type PublishListReq struct {
	Token  string `form:"token"`
	UserID string `form:"user_id"`
}

type PublishListReply struct {
	StatusCode int         `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  []VideoList `json:"video_list"`
}
