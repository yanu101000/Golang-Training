package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"solution1/assignment-3-short-url/entity"
	pb "solution1/assignment-3-short-url/proto/shorturl_service/v1"

	"github.com/redis/go-redis/v9"
	"github.com/teris-io/shortid"
)

type IUrlService interface {
	GetLong(ctx context.Context, req *pb.GetLongReq) (*pb.GetLongRes, error)
	ShortUrl(ctx context.Context, req *pb.ShortUrlReq) (*pb.ShortUrlRes, error)
	Redirect(ctx context.Context, req *pb.RedirectReq) (*pb.RedirectRes, error)
}

type IUrlRepository interface {
	CreateUrl(ctx context.Context, url *entity.Url) (entity.Url, error)
	GetUrlByShortUrl(ctx context.Context, shortUrl string) (entity.Url, error)
}

type urlService struct {
	urlRepo IUrlRepository
	rdb     *redis.Client
}

func NewUrlService(urlRepo IUrlRepository, rdb *redis.Client) IUrlService {
	return &urlService{urlRepo: urlRepo, rdb: rdb}
}

const redisUrlKey = "url:%s"

func (u *urlService) GetLong(ctx context.Context, req *pb.GetLongReq) (*pb.GetLongRes, error) {
	var url entity.Url

	val, err := u.rdb.Get(ctx, fmt.Sprintf(redisUrlKey, req.Short)).Result()
	if err == nil {
		return &pb.GetLongRes{
			LongUrl: val,
		}, nil
	}

	url, err = u.urlRepo.GetUrlByShortUrl(ctx, req.Short)
	if err != nil {
		log.Printf("error get long url: %v", err)
		return nil, err
	}

	err = u.rdb.Set(ctx, fmt.Sprintf(redisUrlKey, req.Short), url.OriginalUrl, 60*time.Second).Err()
	if err != nil {
		log.Println("error set data in redis")
		return nil, err
	}

	return &pb.GetLongRes{
		LongUrl: url.OriginalUrl,
	}, nil
}

func (u *urlService) ShortUrl(ctx context.Context, req *pb.ShortUrlReq) (*pb.ShortUrlRes, error) {
	randomShortId, _ := shortid.Generate()
	url := entity.Url{
		ShortUrl:    randomShortId,
		OriginalUrl: req.Url,
	}

	createdUrl, err := u.urlRepo.CreateUrl(ctx, &url)
	if err != nil {
		log.Printf("error create short url: %v", err)
		return nil, err
	}

	return &pb.ShortUrlRes{
		Id:       int32(createdUrl.ID),
		ShortUrl: createdUrl.ShortUrl,
	}, nil
}

func (u *urlService) Redirect(ctx context.Context, req *pb.RedirectReq) (*pb.RedirectRes, error) {
	var url entity.Url

	val, err := u.rdb.Get(ctx, fmt.Sprintf(redisUrlKey, req.ShortUrl)).Result()
	if err == nil {
		return &pb.RedirectRes{
			Url: val,
		}, nil
	}

	url, err = u.urlRepo.GetUrlByShortUrl(ctx, req.ShortUrl)
	if err != nil {
		log.Printf("error get long url: %v", err)
		return nil, err
	}

	err = u.rdb.Set(ctx, fmt.Sprintf(redisUrlKey, req.ShortUrl), url.OriginalUrl, 60*time.Second).Err()
	if err != nil {
		log.Println("error set data in redis")
		return nil, err
	}

	return &pb.RedirectRes{
		Url: url.OriginalUrl,
	}, nil
}
