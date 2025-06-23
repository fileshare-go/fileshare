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

	handler := newLinkHandler(sharelinkRequest, s.PeerOs(ctx), s.PeerAddress(ctx))

	if resp, ok := handler.preCheck(); ok {
		return resp, nil
	}

	linkCode := util.GenerateCode(config.Cfg().ShareCodeLength)

	handler.PersistShareLink(linkCode)

	handler.PersistRecord()

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
type linkHandler struct {
	OsInfo    string
	PeerAddr  string
	FileInfo  *model.FileInfo
	ShareLink *model.ShareLink
	Request   *pb.ShareLinkRequest
}

func newLinkHandler(shareLinkRequest *pb.ShareLinkRequest, osInfo, peerAddr string) *linkHandler {
	return &linkHandler{
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

// if the second value is true, it indicates no need to do more actions
func (h *linkHandler) preCheck() (*pb.ShareLinkResponse, bool) {
	if service.Orm().Find(h.FileInfo).RowsAffected == 0 {
		return &pb.ShareLinkResponse{
			Status:   pb.Status_ERROR,
			Message:  "File not found",
			LinkCode: "",
		}, true
	}

	if service.Orm().Find(h.ShareLink).RowsAffected == 1 {
		logrus.Debugf("Existing sharelink for %s is %s", util.Render(h.Request.Meta.Sha256[:8]), util.Render(h.ShareLink.LinkCode))
		return &pb.ShareLinkResponse{
			Status:   pb.Status_OK,
			Message:  "Found existing sharelink code!",
			LinkCode: h.ShareLink.LinkCode,
		}, true
	}

	return nil, false
}

// persist changes in database
func (h *linkHandler) PersistShareLink(linkCode string) {
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

func (h *linkHandler) PersistRecord() {
	record := &model.Record{
		Os:             h.OsInfo,
		Sha256:         h.FileInfo.Sha256,
		InteractAction: core.LinkGenAction,
		ClientIp:       h.PeerAddr,
		Time:           time.Now(),
	}
	service.Orm().Save(record)
}
