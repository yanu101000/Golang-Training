package grpc

import (
	"context"
	"log"

	pb "Docker/proto/shorturl_service/v1"
	"Docker/service"
)

type UrlHandler struct {
	pb.UnimplementedUrlServiceServer
	urlService service.IUrlService
}

func NewUrlHandler(urlService service.IUrlService) *UrlHandler {
	return &UrlHandler{urlService: urlService}
}

func (h *UrlHandler) GetLong(ctx context.Context, req *pb.GetLongReq) (*pb.GetLongRes, error) {
	res, err := h.urlService.GetLong(ctx, req)
	if err != nil {
		log.Printf("Error getting long URL: %v", err)
		return nil, err
	}
	return res, nil
}

func (h *UrlHandler) ShortUrl(ctx context.Context, req *pb.ShortUrlReq) (*pb.ShortUrlRes, error) {
	res, err := h.urlService.ShortUrl(ctx, req)
	if err != nil {
		log.Printf("Error shortening URL: %v", err)
		return nil, err
	}
	return res, nil
}

func (h *UrlHandler) Redirect(ctx context.Context, req *pb.RedirectReq) (*pb.RedirectRes, error) {
	res, err := h.urlService.Redirect(ctx, req)
	if err != nil {
		log.Printf("Error redirecting URL: %v", err)
		return nil, err
	}
	return res, nil
}
