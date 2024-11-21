package dto
import (
	"marketing/models"	
)
type SocialtaskDto struct {
	Id           int64  `json:"id"`
	TaskName	 string `json:"task_name"`
	CampaignId   int64  `json:"campaign_id"`
	// CampaignName string `json:"campaign_name"`
	Tags         []string `json:"tag"`
	TypeId         int64 `json:"type_id"`
	TypeName         string `json:"type_name"`
	Keywords     []string `json:"keywords"`
	ExtraTaskIfo models.TaskExtaInfo `json:"extra_task_info"`	
}