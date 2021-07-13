package form

type UploadFileForm struct {
	//传输类型
	TransType string `json:"TransType"`
	//云端路径
	OssPath string `json:"OssPath"`
	//本地路径
	LocalPath string `json:"LocalPath"`
	//存储空间
	Bucket string `json:"Bucket"`
	//节点
	EndPoint string `json:"EndPoint"`
	//最大重试次数
	MaxTry string `json:"MaxTry"`
	//是否需要压缩包操作  0：无  1：需要压缩 2：需要解压
	NeedZip int `json:"NeedZip"`
	//NeedZip非0时，需要用到该字段
	ZipTarget string `json:"ZipTarget"`
	//私有数据
	Private string `json:"Private"`
	//任务类型
	TaskType string `json:"TaskType"`
	//属于某订单
	Order string `json:"Order"`
	//用户信息
	User string `json:"User"`
	//企业信息
	Company string `json:"Company"`
	//回调地址(回调地址合并后删除)
	CallBackUrl string `json:"CallBackUrl"`
	//回调数据
	CallBackData string `json:"CallBackData"`
}

// CheckExistsForm /*四、查询本地文件*/
type CheckExistsForm struct {
	Path string `json:"Path"`
}

type CheckExistsRBForm struct {
	Path string `json:"Path"`
	Exists bool `json:"Exists"`
}

// ReadFromFileForm /*十、读取文件*/
type ReadFromFileForm struct {
	Path string `json:"Path"`
}

type ReadFromFileRBForm struct {
	ReadFromFileForm
	Data string `json:"Data"`
}

// CompressForm /*十一、压缩*/
type CompressForm struct {
	Source string `json:"Source"`
	Dest string `json:"Dest"`
}
type CompressRBForm struct {
	Ret bool `json:"Ret"`
}

// DecompressForm /*十二、解压*/
type DecompressForm struct {
	Source string `json:"Source"`
	Dest string `json:"Dest"`
}
type DecompressRBForm struct {
	Ret bool `json:"Ret"`
}

// FileTransferRequestForm /*文件上传&下载回调*/
type FileTransferRequestForm struct {
	FileSize string `json:"FileSize"`//文件大小
	CallBackData string `json:"CallBackData"` //回调结果
}

// ServerResponseForm /*服务器返回结果*/
type ServerResponseForm struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}