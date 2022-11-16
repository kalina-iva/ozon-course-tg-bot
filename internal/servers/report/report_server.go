package report

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/api/report"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"google.golang.org/grpc"
)

type server struct {
	report.UnimplementedReportServer
	msgModel *messages.Model
}

func NewReportServer(msgModel *messages.Model, serverAddress string) error {
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	s := grpc.NewServer()
	report.RegisterReportServer(s, &server{msgModel: msgModel})
	logger.Info("listening report messages")
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve")
	}
	return nil
}

func (s *server) SendReport(ctx context.Context, request *report.ReportRequest) (*report.ReportReply, error) {
	err := s.msgModel.SendReport(ctx, request.GetReport(), request.GetUserID())
	if err != nil {
		return nil, errors.Wrap(err, "cannot send report")
	}
	return &report.ReportReply{Message: "successfully sent"}, nil
}
