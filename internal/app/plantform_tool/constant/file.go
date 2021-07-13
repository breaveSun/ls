package constant

const(
	/*任务内容数据*/
	TransType    = "TransType" //上传
	UploadKey    = "upload"    //上传
	DownLoadKey  = "download"  //下载
	LocalPathKey = "LocalPathKey" //本地路径
	FileSizeKey  = "FileSize"  //文件大小
	CallBackUrlKey = "CallBackUrl" //回调地址
	CallBackDataKey = "CallBackData"//回调数据
	TaskStatusKey = "Status"//任务状态
	TaskStatusOver = 4 //完成
	MaxTryKey = "MaxTry"
	CallBackSurplusCountKey = "CallBackSurplusCount" //剩余回调次数

	NeedZipKey = "NeedZip"
	NeedZip = 1 //需要压缩
	NeedUnZip = 2 //需要解压
	ZipTargetKey = "ZipTargetKey" //解压|压缩之后存储位置

	TaskIdKey = "TaskId" //任务id
	CreateTimeKey = "CreateTime" //任务创建时间

	/*redis key*/
	TaskIdRDKey = "TaskId" //自增idkey
	TaskPollingListRDKey = "TaskPollingList"//未完成任务列表key
)