package sharelink

import (
	"context"
	"time"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/core"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/chanmaoganda/fileshare/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type ShareLinkServer struct {
	pb.UnimplementedShareLinkServiceServer
}

func NewShareLinkServer() *ShareLinkServer {
	return &ShareLinkServer{}
}

func (s *ShareLinkServer) GenerateLink(ctx context.Context, sharelinkRequest *pb.ShareLinkRequest) (*pb.ShareLinkResponse, error) {
	logrus.Debugf("Generating sharelink for %s", util.Render(sharelinkRequest.Meta.Sha256[:8]))

	handler := NewLinkHandler(sharelinkRequest, s.PeerOs(ctx), s.PeerAddress(ctx))

	if service.Orm().Find(handler.FileInfo).RowsAffected == 0 {
		return &pb.ShareLinkResponse{
			Status:   pb.Status_ERROR,
			Message:  "File not found",
			LinkCode: "",
		}, nil
	}

	if service.Orm().Find(handler.ShareLink).RowsAffected == 1 {
		logrus.Debugf("Existing sharelink for %s is %s", util.Render(sharelinkRequest.Meta.Sha256[:8]), util.Render(handler.ShareLink.LinkCode))
		return &pb.ShareLinkResponse{
			Status:   pb.Status_OK,
			Message:  "Found existing sharelink code!",
			LinkCode: handler.ShareLink.LinkCode,
		}, nil
	}

	linkCode := util.GenerateCode(config.Cfg().ShareCodeLength)

	handler.PersistShareLink(linkCode)
	handler.PersistRecords()

	logrus.Debugf("Generated sharelink for %s is %s", util.Render(sharelinkRequest.Meta.Sha256[:8]), util.Render(linkCode))

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

// handles share link
type LinkHandler struct {
	OsInfo    string
	PeerAddr  string
	FileInfo  *model.FileInfo
	ShareLink *model.ShareLink
	Request   *pb.ShareLinkRequest
}

func NewLinkHandler(shareLinkRequest *pb.ShareLinkRequest, osInfo, peerAddr string) *LinkHandler {
	return &LinkHandler{
		OsInfo:   osInfo,
		PeerAddr: peerAddr,
		FileInfo: &model.FileInfo{
			Sha256: shareLinkRequest.Meta.Sha256,
		},
		ShareLink: &model.ShareLink{
			Sha256: shareLinkRequest.Meta.Sha256,
		},
		Request: shareLinkRequest,
	}
}

// persist changes in database
func (h *LinkHandler) PersistShareLink(linkCode string) {
	h.ShareLink.LinkCode = linkCode
	h.ShareLink.CreatedAt = time.Now()
	h.ShareLink.CreatedBy = h.PeerAddr

	if h.Request.ValidDays == 0 {
		logrus.Debug("[ShareLink] request days invalid, Using default valid days")
		h.ShareLink.OutdatedAt = time.Now().AddDate(0, 0, int(config.Cfg().ValidDays))
	} else {
		h.ShareLink.OutdatedAt = time.Now().AddDate(0, 0, int(h.Request.ValidDays))
	}

	service.Orm().Save(h.ShareLink)
}

func (h *LinkHandler) PersistRecords() {
	service.Mgr().InsertRecord(h.MakeRecord())
}

func (h *LinkHandler) MakeRecord() *model.Record {
	return &model.Record{
		Os:             h.OsInfo,
		Sha256:         h.FileInfo.Sha256,
		InteractAction: core.LinkGenAction,
		ClientIp:       h.PeerAddr,
		Time:           time.Now(),
	}
}
