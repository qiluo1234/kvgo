package fio

const DataFilePerm = 0644

// 抽象IO管理接口，可以接入不同的IO类型，目前支持标准文件IO
type IOManager interface {

	//Read 从文件指定位置读取对应的数据
	Read([]byte, int64) (int, error)

	//Write 写入字节数据到文件中
	Write([]byte) (int, error)

	//Sync持久化数据
	Sync() error

	//Close 关闭文件
	Close() error

	// Size 获取到文件大小
	Size() (int64, error)
}

// NewIOManager 初始化 IOManager，目前只支持标准 FileIO
func NewIOManager(fileName string) (IOManager, error) {
	return NewFileIOManager(fileName)
}
