package douyin

import (
	"log"
	"testing"
)

func TestUserInfo(t *testing.T) {
	accessToken := "act.3.14FVDCPbj3NDTMpIjb1G-X8HmUHPSKIFE5tFmJgjfmMvjtEG00t9J5PhOwNTuk7iZMno9Y69hgTpdqmVa90E8VVyX_GgUhIHlGoouDBmhNnjkKiAzjW_nch0jPqAtk4C2WoNxCzXFnmb1FQNN6C1Xaj6oZ4TNwyCegtBqqRP7m7GUY0mbKpTzyTlvYg=_lq"
	openId := "_000J7VTtKDbE7OqIZogdj6dS_HPXK7o7RU3"
	client := NewClient(Conf{
		ClientKey:    "awezk63nfz38px7q",
		ClientSecret: "3b2a495c51bd9baf30b3851cf066f578",
		DirectURL:    "https://douyinshanghu.sis868.com/auth",
		Scopes:       "user_info,trial.whitelist,data.external.user",
	})
	fans, err := client.GetUserFans(openId, accessToken, 7)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fans)
}
