package grpcreport

import (
	"github.com/pkg/errors"
	reportClient "gitlab.ozon.dev/mary.kalina/telegram-bot/internal/api/report"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/helper/grpcconn"
	"google.golang.org/grpc"
)

var conn *grpc.ClientConn

func NewReportClient(serverAddr string) (reportClient.ReportClient, error) {
	var err error
	conn, err = grpcconn.NewClientConn(serverAddr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create grpcreport client connection")
	}

	return reportClient.NewReportClient(conn), nil
}

func Close() error {
	return conn.Close()
}
