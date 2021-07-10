1、日志的使用方法：lib.Logger.Info("msg",zap.String("key","string"))

2、编译:go build -tags=jsoniter -o .\build\package\main.exe  -ldflags "-H=windowsgui" .\cmd\plantform_tool\main.go