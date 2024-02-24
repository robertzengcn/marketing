package controllers
import (
	"github.com/beego/i18n"
	"marketing/models"
	"marketing/dto"
)

type TagController struct {
	BaseController
	i18n.Locale
}

//get tag list by account id
func (c *TagController) ListTag() {
	// accountId, err := c.GetInt64("account_id")
	
	// if err != nil {
	// 	c.ErrorJson(20240107153815, err.Error(), nil)
	// }
	uid := c.GetSession("uid")
	accountId:=uid.(int64)
	tagModel := models.Tag{}
	tagList, err := tagModel.ListTagByAccountId(accountId)
	if err != nil {
		c.ErrorJson(20240107153921, err.Error(), nil)
	}
	var tagListDto []dto.TagDto
	//loop tag list 
	for _, tag := range tagList {
		tagListDto=append(tagListDto,dto.TagDto{
			Id: tag.TagId,
			Name: tag.TagName,
		} )
	}
	c.SuccessJson(tagListDto)
}
//get tag by keywords
func (c *TagController) GetTagKeyword() {
	
}