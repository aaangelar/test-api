package user

import (
	"time"
	"golang.org/x/net/context"
	"github.com/go-kit/kit/log"
)


type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetExportTypeTemplate(ctx context.Context, user_id string) (userDetails []User{}, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetUsers", "user_id", user_id,"took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetUsers(ctx, user_id)
}
