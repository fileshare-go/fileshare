package web

import (
	"net/http"

	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebService struct {
	*gin.Engine
}

func NewWebService(DB *gorm.DB) *WebService {
	service := &WebService{
		Engine: gin.Default(),
	}
	service.RegisterRoutes()
	return service
}

func (s *WebService) RegisterRoutes() {
	s.GET("/fileinfo", s.GetFileInfos)
	s.GET("/sharelink", s.GetShareLinks)
	s.GET("/record", s.GetRecords)
}

func (s *WebService) GetFileInfos(c *gin.Context) {
	var fileInfos []model.FileInfo
	service.Orm().
		Model(&model.FileInfo{}).
		Preload("Record").
		Preload("Link").
		Find(&fileInfos)

	c.JSON(http.StatusOK, gin.H{"data": fileInfos})
}

func (s *WebService) GetShareLinks(c *gin.Context) {
	shareLinks, err := service.Mgr().GetShareLink()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shareLinks})
}

func (s *WebService) GetRecords(c *gin.Context) {
	records, err := service.Mgr().GetRecord()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}
