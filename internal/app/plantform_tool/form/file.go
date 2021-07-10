package form

type UploadFileForm struct {
	Id uint32 `json:"id"`
}
/*四、查询本地文件*/
type CheckExistsForm struct {
	Path string `json:"Path"`
}

type CheckExistsRBForm struct {
	Path string `json:"Path"`
	Exists bool `json:"Exists"`
}
/*十、读取文件*/
type ReadFromFileForm struct {
	Path string `json:"Path"`
}

type ReadFromFileRBForm struct {
	ReadFromFileForm
	Data string `json:"Data"`
}
/*十一、压缩*/
type CompressForm struct {
	Source string `json:"Source"`
	Dest string `json:"Dest"`
}
type CompressRBForm struct {
	Ret bool `json:"Ret"`
}
/*十二、解压*/
type DecompressForm struct {
	Source string `json:"Source"`
	Dest string `json:"Dest"`
}
type DecompressRBForm struct {
	Ret bool `json:"Ret"`
}
