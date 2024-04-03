package controllers

import (
	"github.com/beego/i18n"
	"marketing/models"
)

type SocialPlatfromController struct {
	BaseController
	i18n.Locale
}

func (c *SocialPlatfromController) ChildPrepare() {

}
//list social platform
func (c *SocialPlatfromController) Listplatform() {
	socialPlModel := models.SocialPlatform{}
	socialPls, err := socialPlModel.Listsocialplatform()
	if err != nil {
		c.ErrorJson(202304041002135, "get social platform error", nil)
	}
	c.SuccessJson(socialPls)
}
