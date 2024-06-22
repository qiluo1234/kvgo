package bitcast_go

import "os"

type Options struct {
	//数据库数据目录
	DirPath string

	//活跃数据阈值
	DataFileSize int64

	//每次些数据选择是否持久化
	SyncWrites bool

	//索引类型
	IndexType IndexerType
}

// IteratorOptions 索引迭代器配置项
type IteratorOptions struct {
	// 遍历前缀为指定值的 Key，默认为空
	Prefix []byte
	// 是否反向遍历，默认 false 是正向的
	Reverse bool
}

type IndexerType = int8

const (
	//BTree 索引
	BTree IndexerType = iota + 1

	// ART Adpataive Radix Tree 自适应基数树索引
	ART
)

var DefaultOptions = Options{
	DirPath:      os.TempDir(),
	DataFileSize: 256 * 1024 * 1024, //256MB
	SyncWrites:   false,
	IndexType:    BTree,
}

var DefaultIteratorOptions = IteratorOptions{
	Prefix:  nil,
	Reverse: false,
}
