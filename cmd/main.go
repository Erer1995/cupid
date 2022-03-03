package main

import (
	"context"
	"errors"

	"github.com/zzerer/cupid/log"
)

func main() {

	//使用默认的全局logger
	ctx := context.Background()

	//默认info级别
	log.D(ctx).Message("debug level message")
	log.I(ctx).Message("info level message")
	//{"level":"info","time":"2022-03-03T11:29:13.085192+08:00","message":"info level message"}

	//设置全局为debug级别
	log.SetGlobalLevel("DEBUG")
	log.D(ctx).Message("debug level message")
	log.I(ctx).Message("info level message")
	log.W(ctx).Message("warn level message")
	log.E(ctx).Message("error level message")
	//log.F(ctx).Message("fatal level message") //执行完后会执行 os.Exit(1)
	//log.P(ctx).Message("panic level message") //执行完后会执行 panic()
	//{"level":"debug","time":"2022-03-03T14:40:16.257813+08:00","message":"debug level message"}
	//{"level":"info","time":"2022-03-03T14:40:16.257819+08:00","message":"info level message"}
	//{"level":"warn","time":"2022-03-03T14:40:16.257823+08:00","message":"warn level message"}
	//{"level":"error","time":"2022-03-03T14:40:16.257826+08:00","message":"error level message"}
	//{"level":"fatal","time":"2022-03-03T14:40:16.25783+08:00","message":"fatal level message"}
	//{"level":"panic","time":"2022-03-03T14:41:32.980277+08:00","message":"panic level message"}

	//启用带上调用信息，默认关闭
	log.EnableCaller(true)
	log.D(ctx).Message("debug level message with caller")
	//{"level":"debug","time":"2022-03-03T11:29:13.085394+08:00","caller":"/Users/zengyuzhao/gopath/src/github/cupid/log/log.go:194","message":"debug level message with caller"}

	//添加全局自定义字段
	fields := make(map[string]string)
	fields["key_1"] = "value_1"
	fields["key_2"] = "value_2"
	log.AddGlobalFields(fields)
	log.D(ctx).Message("debug level message with global fields")
	//{"level":"debug","key_1":"value_1","key_2":"value_2","time":"2022-03-03T11:41:14.784602+08:00","message":"debug level message with global fields"}

	//单条日志添加自定义字段
	log.D(ctx).Field("key_3", "value_3").Message("debug level message with global fields and extra field")
	//{"level":"debug","key_1":"value_1","key_2":"value_2","key_3":"value_3","time":"2022-03-03T11:45:52.739347+08:00","message":"debug level message with global fields and extra field"}

	//日志带上Error信息
	err := errors.New("new error")
	log.E(ctx).Err(err).Message("error message with detail")
	//{"level":"error","key_1":"value_1","key_2":"value_2","error":"new error","time":"2022-03-03T11:47:40.062588+08:00","message":"error message with detail"}

	//ctx中注入trace_id
	ctx = log.WithLogTraceId(ctx, "test_traceId")
	log.D(ctx).Message("message with ctx trace_id")

	//使用logger实例
	mylogger_fields := make(map[string]string)
	mylogger_fields["logger_key1"] = "logger_value1"
	mylogger := log.NewRlogger(true, mylogger_fields)
	mylogger.D(ctx).Message("mylogger debug level message")
	mylogger.D(ctx).Field("logger_key2", "logger_value2").Err(err).Message("mylogger debug level with full test")
}
