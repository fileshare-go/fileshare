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
	s.GET("/fileinfo", s.GetFileInfo)
}

func (s *WebService) GetFileInfo(c *gin.Context) {
	var fileInfos []model.FileInfo
	s.Manager.DB.Find(&fileInfos)

	c.JSON(http.StatusOK, gin.H{"data": fileInfos})
}
