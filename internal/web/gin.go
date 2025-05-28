package web

import (
	"net/http"

	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/dbmanager"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebService struct {
	*gin.Engine
	Manager *dbmanager.DBManager
}

func NewWebService(DB *gorm.DB) *WebService {
	service := &WebService{
		Manager: dbmanager.NewDBManager(DB),
		Engine:  gin.Default(),
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
	s.Manager.DB.Find(&fileInfos)

	c.JSON(http.StatusOK, gin.H{"data": fileInfos})
}

func (s *WebService) GetShareLinks(c *gin.Context) {
	var shareLinks []model.ShareLink
	s.Manager.DB.Find(&shareLinks)

	c.JSON(http.StatusOK, gin.H{"data": shareLinks})
}

func (s *WebService) GetRecords(c *gin.Context) {
	var records []model.Record
	s.Manager.DB.Find(&records)

	c.JSON(http.StatusOK, gin.H{"data": records})
}
