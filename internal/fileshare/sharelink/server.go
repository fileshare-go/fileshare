package sharelink

import (
	"context"
	"time"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/pkg/debugprint"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"gorm.io/gorm"
)

type ShareLinkServer struct {
	pb.UnimplementedShareLinkServiceServer
	Settings *config.Settings
	Manager  *dbmanager.DBManager
	RangGen  *RandStringGen
}

func NewShareLinkServer(settings *config.Settings, DB *gorm.DB) *ShareLinkServer {
	return &ShareLinkServer{
		Settings: settings,
		RangGen:  NewRandStringGen(),
		Manager:  dbmanager.NewDBManager(DB),
	}
}

func (s *ShareLinkServer) GenerateLink(ctx context.Context, sharelinkRequest *pb.ShareLinkRequest) (*pb.ShareLinkResponse, error) {
	logrus.Debugf("Generating sharelink for %s", debugprint.Render(sharelinkRequest.Meta.Sha256[:8]))

	handler := NewLinkHandler(sharelinkRequest, s.Settings, s.PeerOs(ctx), s.PeerAddress(ctx), s.Manager)

	if !handler.Manager.SelectFileInfo(handler.FileInfo) {
		return &pb.ShareLinkResponse{
			Status:   pb.Status_ERROR,
			Message:  "File not found",
			LinkCode: "",
		}, nil
	}

	if handler.Manager.SelectValidShareLink(handler.ShareLink) {
		logrus.Debugf("Existing sharelink for %s is %s", debugprint.Render(sharelinkRequest.Meta.Sha256[:8]), debugprint.Render(handler.ShareLink.LinkCode))
		return &pb.ShareLinkResponse{
			Status:   pb.Status_OK,
			Message:  "Found existing sharelink code!",
			LinkCode: handler.ShareLink.LinkCode,
		}, nil
	}

	linkCode := s.RangGen.generateCode(s.Settings.ShareCodeLength)

	handler.PersistShareLink(linkCode)
	handler.PersistRecords()

	logrus.Debugf("Generated sharelink for %s is %s", debugprint.Render(sharelinkRequest.Meta.Sha256[:8]), debugprint.Render(linkCode))

	return &pb.ShareLinkResponse{
		Status:   pb.Status_OK,
		Message:  "Generated code for sharelink",
		LinkCode: linkCode,
	}, nil
}

func (s *ShareLinkServer) PeerAddress(ctx context.Context) string {
	peer, ok := peer.FromContext(ctx)
	if ok {
		return peer.Addr.String()
	}
	return "unknown"
}

func (s *ShareLinkServer) PeerOs(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "unknown"
	}

	if osInfo, ok := md["os"]; ok && len(osInfo) != 0 {
		return osInfo[0]
	}
	return "unknown"
}

type LinkHandler struct {
	OsInfo    string
	PeerAddr  string
	Settings *config.Settings
	FileInfo  *model.FileInfo
	ShareLink *model.ShareLink
	Manager   *dbmanager.DBManager
	Request   *pb.ShareLinkRequest
}

func NewLinkHandler(shareLinkRequest *pb.ShareLinkRequest, settings *config.Settings, osInfo, peerAddr string, manager *dbmanager.DBManager) *LinkHandler {
	return &LinkHandler{
		Settings: settings,
		OsInfo:   osInfo,
		PeerAddr: peerAddr,
		FileInfo: &model.FileInfo{
			Sha256: shareLinkRequest.Meta.Sha256,
		},
		ShareLink: &model.ShareLink{
			Sha256: shareLinkRequest.Meta.Sha256,
		},
		Manager: manager,
		Request: shareLinkRequest,
	}
}

func (h *LinkHandler) PersistShareLink(linkCode string) {
	h.ShareLink.LinkCode = linkCode
	h.ShareLink.CreatedAt = time.Now()
	h.ShareLink.CreatedBy = h.PeerAddr

	if h.Request.ValidDays == 0 {
		logrus.Debug("[ShareLink] request days invalid, Using default valid days")
		h.ShareLink.OutdatedAt = time.Now().AddDate(0, 0, int(h.Settings.ValidDays))
	} else {
		h.ShareLink.OutdatedAt = time.Now().AddDate(0, 0, int(h.Request.ValidDays))
	}

	h.Manager.UpdateShareLink(h.ShareLink)
}


func (h *LinkHandler) PersistRecords() {
	h.Manager.CreateRecord(h.MakeRecord())
}

func (h *LinkHandler) MakeRecord() *model.Record {
	return &model.Record{
		Os:             h.OsInfo,
		Sha256:         h.FileInfo.Sha256,
		InteractAction: "linkgen",
		ClientIp:       h.PeerAddr,
		Time:           time.Now(),
	}
}
