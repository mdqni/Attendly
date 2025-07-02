package grpc

import (
	"context"
	qrv1 "github.com/mdqni/Attendly/proto/gen/go/qr/v1"
	"github.com/mdqni/Attendly/services/qr/internal/service"

	"google.golang.org/grpc"
)

type qrServer struct {
	qrv1.UnimplementedQRServiceServer
	service service.QrService
}

func (q qrServer) GenerateQR(ctx context.Context, request *qrv1.GenerateQRRequest) (*qrv1.GenerateQRResponse, error) {

	qr, err := q.service.GenerateQR(ctx, request.GetLessonId(), request.GetTeacherId(), request.GetExpiresUnix())
	if err != nil {
		return nil, err
	}
	return &qrv1.GenerateQRResponse{QrCode: qr}, nil
}

func (q qrServer) ValidateQR(ctx context.Context, request *qrv1.ValidateQRRequest) (*qrv1.ValidateQRResponse, error) {
	lessonID, err := q.service.ValidateQR(ctx, request.GetQrCode())
	if err != nil {
		return nil, err
	}
	return &qrv1.ValidateQRResponse{
		Valid:    true,
		LessonId: lessonID,
	}, nil
}

func Register(gRPCServer *grpc.Server, svc service.QrService) {
	qrv1.RegisterQRServiceServer(gRPCServer, &qrServer{service: svc})
}
