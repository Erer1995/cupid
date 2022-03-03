package log

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var default_rlogger *Rlogger
var logOnce sync.Once

type ctxTraceId string

const (
	CTX_TRACE_ID ctxTraceId = "TRACE_ID"
)

//定义一个Rlogger类型，屏蔽底层实现
type Rlogger struct {
	logger zerolog.Logger //底层用zerolog实现
	//logger logrus.Logger //底层用logrus实现
}

type Revent struct {
	event *zerolog.Event //底层用zerolog实现
	//event *logrus.Entry //底层用logrus实现
}

//初始化
func InitRlogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	//尝试从环境变量取得全局LEVEL, 没设置默认Info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	default_rlogger = &Rlogger{
		logger: log.With().Logger(),
	}
	log.Logger = default_rlogger.logger
}

//设置全局log level
func SetGlobalLevel(level string) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

//设置是否开启日志中带上调用信息
func EnableCaller(enable bool) {
	if enable {
		default_rlogger.logger = default_rlogger.logger.With().Caller().Logger()
		log.Logger = default_rlogger.logger
	}
}

//添加全局自定义字段，必须是string类型
func AddGlobalFields(global_fields map[string]string) {
	for k, v := range global_fields {
		default_rlogger.logger = default_rlogger.logger.With().Str(k, v).Logger()
	}
	log.Logger = default_rlogger.logger
}

//返回一个自定义的Rlogger
func NewRlogger(enable_caller bool, global_fields map[string]string) *Rlogger {
	logger := zerolog.New(os.Stdout).With().Logger()
	if enable_caller {
		logger = logger.With().Caller().Logger()
	}
	for k, v := range global_fields {
		logger = logger.With().Str(k, v).Logger()
	}
	return &Rlogger{
		logger: logger,
	}
}

func newRevent(level string, rlogger *Rlogger) *Revent {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return &Revent{
			event: rlogger.logger.Debug(),
		}
	case "INFO":
		return &Revent{
			event: rlogger.logger.Info(),
		}
	case "WARN":
		return &Revent{
			event: rlogger.logger.Warn(),
		}
	case "ERROR":
		return &Revent{
			event: rlogger.logger.Error(),
		}
	case "FATAL":
		return &Revent{
			event: rlogger.logger.Fatal(),
		}
	case "PANIC":
		return &Revent{
			event: rlogger.logger.Panic(),
		}
	default:
		return &Revent{
			event: default_rlogger.logger.Info(),
		}
	}
}

//默认接口DEBUG级别的日志
func D(ctx context.Context) *Revent {
	return injectTraceId(newRevent("DEBUG", default_rlogger), ctx)
}

//默认接口INFO级别的日志
func I(ctx context.Context) *Revent {
	return injectTraceId(newRevent("INFO", default_rlogger), ctx)
}

//默认接口WARN级别的日志
func W(ctx context.Context) *Revent {
	return injectTraceId(newRevent("WARN", default_rlogger), ctx)
}

//默认接口ERROR级别的日志
func E(ctx context.Context) *Revent {
	return injectTraceId(newRevent("ERROR", default_rlogger), ctx)
}

//默认接口FATAL级别的日志
func F(ctx context.Context) *Revent {
	return injectTraceId(newRevent("FATAL", default_rlogger), ctx)
}

//默认接口PANIC级别的日志
func P(ctx context.Context) *Revent {
	return injectTraceId(newRevent("PANIC", default_rlogger), ctx)
}

//logger实例DEBUG级别的日志
func (rl *Rlogger) D(ctx context.Context) *Revent {
	return injectTraceId(newRevent("DEBUG", rl), ctx)
}

//logger实例INFO级别的日志
func (rl *Rlogger) I(ctx context.Context) *Revent {
	return injectTraceId(newRevent("INFO", rl), ctx)
}

//logger实例WARN级别的日志
func (rl *Rlogger) W(ctx context.Context) *Revent {
	return injectTraceId(newRevent("WARN", rl), ctx)
}

//logger实例ERROR级别的日志
func (rl *Rlogger) E(ctx context.Context) *Revent {
	return injectTraceId(newRevent("ERROR", rl), ctx)
}

//logger实例FATAL级别的日志
func (rl *Rlogger) F(ctx context.Context) *Revent {
	return injectTraceId(newRevent("FATAL", rl), ctx)
}

//logger实例PANIC级别的日志
func (rl *Rlogger) P(ctx context.Context) *Revent {
	return injectTraceId(newRevent("PANIC", rl), ctx)
}

//输出日志
func (re *Revent) Message(content string) {
	re.event.Msg(content)
}

//单条日志带上自定义字段
func (re *Revent) Field(key, value string) *Revent {
	re.event = re.event.Str(key, value)
	return re
}

//日志带上具体的error信息
func (re *Revent) Err(err error) *Revent {
	re.event = re.event.Str("error", err.Error())
	return re
}

//入参：ctx, trace_id
//返回：如果trace_id不为空，返回一个带trace_id的ctx
func WithLogTraceId(ctx context.Context, trace_id string) context.Context {
	if ctx == nil || trace_id == "" {
		return ctx
	}
	return context.WithValue(ctx, CTX_TRACE_ID, trace_id)
}

//从ctx中取出trace_id注入到revent中
func injectTraceId(revent *Revent, ctx context.Context) *Revent {
	var trace_id string
	raw := ctx.Value(CTX_TRACE_ID)
	if raw == nil {
		return revent
	} else {
		trace_id = raw.(string)
		revent.event = revent.event.Str("trace_id", trace_id)
	}
	return revent
}

//init
func init() {
	logOnce.Do(func() {
		InitRlogger()
	})
}
