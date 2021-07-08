package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"time"
)

//只能输出结构化日志，但是性能要高于 SugaredLogger
var Logger *zap.Logger

func InitLog(){
	config := zapcore.EncoderConfig{
		MessageKey:   "msg",  //结构化（json）输出：msg的key
		LevelKey:     "level",//结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "ts",   //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file", //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder, //采用短文件路径编码输出（test/main.go:14	）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},//输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},//输出时间戳格式
	}
	//自定义日志级别：自定义Info级别
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	//自定义日志级别：自定义Info级别
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})

	//自定义日志级别：自定义Warn级别
	errLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	infoWriter := getWriter("tmp/info")
	errWriter := getWriter("tmp/err")
	// 实现多个输出
	core := zapcore.NewTee(
		//将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), infoLevel),
		//warn及以上写入errPath，NewConsoleEncoder 是非结构化输出
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(errWriter), errLevel),
		//同时将日志输出到控制台，NewJSONEncoder 是结构化输出
		zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), debugLevel),
	)
	Logger = zap.New(core, zap.AddCaller())
	Logger.Info("日志初始化成功")
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 filename.YYmmddHH
	// filename是指向最新日志的链接
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotatelogs.WithRotationTime(time.Hour*24), //切割频率 24小时
	)
	if err != nil {
		fmt.Println("日志启动异常")
		log.Println("日志启动异常")
		panic(err)
	}
	return hook
}
