package dropbox

type accountInfo struct {
	Country     string `json:"country"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	QuotaInfo   struct {
		Datastores float64 `json:"datastores"`
		Normal     float64 `json:"normal"`
		Quota      float64 `json:"quota"`
		Shared     float64 `json:"shared"`
	} `json:"quota_info"`
	ReferralLink string      `json:"referral_link"`
	Team         interface{} `json:"team"`
	Uid          float64     `json:"uid"`
}
