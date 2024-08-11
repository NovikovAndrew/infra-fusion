package sentry

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/getsentry/sentry-go"
)

type config interface {
	GetAppName() string
	GetSentryDNS() string
	GetEnvironment() string
}

func InitSentryProvide(cfg config) error {
	dns := cfg.GetSentryDNS()
	if dns == "" {
		return fmt.Errorf("sentry dns is empty")
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn:                dns,
		Debug:              cfg.GetEnvironment() != "production",
		EnableTracing:      true,
		TracesSampleRate:   1,
		ProfilesSampleRate: 1,
		ServerName:         cfg.GetAppName(),
		Environment:        cfg.GetEnvironment(),
	})
}

func UnaryServerTracerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		md, _ := metadata.FromIncomingContext(ctx)
		options := []sentry.SpanOption{
			sentry.WithOpName("grpc.server"),
			continueFromGrpcMetadata(md),
			sentry.WithTransactionSource(sentry.SourceURL),
		}

		transaction := sentry.StartTransaction(ctx,
			fmt.Sprint(info.FullMethod),
			options...,
		)
		defer transaction.Finish()

		resp, err := handler(transaction.Context(), req)
		if err != nil {
			hub.CaptureException(err)
		}

		transaction.Status = toSpanStatus(status.Code(err))

		return resp, err
	}
}

func continueFromGrpcMetadata(md metadata.MD) sentry.SpanOption {
	traces := md[sentry.SentryTraceHeader]
	bgs := md[sentry.SentryBaggageHeader]

	trace := ""
	if len(traces) == 1 {
		trace = traces[0]
	}

	baggage := ""
	if len(bgs) == 1 {
		baggage = bgs[0]
	}

	return sentry.ContinueFromHeaders(trace, baggage)
}

func toSpanStatus(code codes.Code) sentry.SpanStatus {
	switch code {
	case codes.OK:
		return sentry.SpanStatusOK
	case codes.Canceled:
		return sentry.SpanStatusCanceled
	case codes.Unknown:
		return sentry.SpanStatusUnknown
	case codes.InvalidArgument:
		return sentry.SpanStatusInvalidArgument
	case codes.DeadlineExceeded:
		return sentry.SpanStatusDeadlineExceeded
	case codes.NotFound:
		return sentry.SpanStatusNotFound
	case codes.AlreadyExists:
		return sentry.SpanStatusAlreadyExists
	case codes.PermissionDenied:
		return sentry.SpanStatusPermissionDenied
	case codes.ResourceExhausted:
		return sentry.SpanStatusResourceExhausted
	case codes.FailedPrecondition:
		return sentry.SpanStatusFailedPrecondition
	case codes.Aborted:
		return sentry.SpanStatusAborted
	case codes.OutOfRange:
		return sentry.SpanStatusOutOfRange
	case codes.Unimplemented:
		return sentry.SpanStatusUnimplemented
	case codes.Internal:
		return sentry.SpanStatusInternalError
	case codes.Unavailable:
		return sentry.SpanStatusUnavailable
	case codes.DataLoss:
		return sentry.SpanStatusDataLoss
	case codes.Unauthenticated:
		return sentry.SpanStatusUnauthenticated
	default:
		return sentry.SpanStatusUndefined
	}
}
