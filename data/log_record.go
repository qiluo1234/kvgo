package data

import (
	"encoding/binary"
	"hash/crc32"
)

// 写入到数据文件的记录
// 之所以叫日志，是因为数据文件中的数据追加写入的，类似日志的格式
type LogRecordType = byte

const (
	LogRecordNormal LogRecordType = iota
	LogRecordDeleted
	LogRecordTxnFinished
)

// crc type keySize valueSize
// 4+ 1 + 5 + 5 = 15
const maxLogRecordHeaderSize = binary.MaxVarintLen32*2 + 5

type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecordType
}

// LogRecord 的头部信息
type logRecordHeader struct {
	crc        uint32        //crc 校验值
	recordType LogRecordType //标识 LogRecord 的类型
	keySize    uint32        // key的长度
	valueSize  uint32        //value的长度
}

// LogRecordPos 数据内存索引，主要描述数据在磁盘的位置
type LogRecordPos struct {
	Fid    uint32 //文件 id 表示讲数据存储到了哪个文件当中
	Offset int64  //偏移，表示将数据存储到了数据文件哈总的哪个位置
}

// TransactionRecord 暂存的事务相关的数据
type TransactionRecord struct {
	Record *LogRecord
	Pos    *LogRecordPos
}

// EncodeLogRecord 对LogRecord 进行编码，返回字节数组以及长度
// +-----------+---------+-----------+-------------+-------+--------+
// | crc校验值  | type类型  | keysize  | value size  |  key  |  value |
// +-----------+---------+-----------+-------------+-------+--------+
//
//	4字节       1字节    变长（最大5）     变长（最大5）     变长       变长
func EncodeLogRecord(logRecord *LogRecord) ([]byte, int64) {
	//初始化一个 header 部分的字节数组
	header := make([]byte, maxLogRecordHeaderSize)

	// 第五个字节存储Type
	header[4] = logRecord.Type
	var index = 5
	// 5 字节之后，存储的是key和value的长度信息
	// 使用变长类型，节省空间
	index += binary.PutVarint(header[index:], int64(len(logRecord.Key)))
	index += binary.PutVarint(header[index:], int64(len(logRecord.Value)))

	var size = index + len(logRecord.Key) + len(logRecord.Value)
	encBytes := make([]byte, size)

	// 将 header 部分内容拷贝过来
	copy(encBytes[:index], header[:index])
	//将 key 和 value 数据拷贝到字节数组中
	copy(encBytes[index:], logRecord.Key)
	copy(encBytes[index+len(logRecord.Key):], logRecord.Value)

	//对整个 LogRecord的数据进行crc校验
	crc := crc32.ChecksumIEEE(encBytes[4:])
	binary.LittleEndian.PutUint32(encBytes[:4], crc)

	return encBytes, int64(size)
}

// EncodeLogRecordPos 对位置信息进行编码
func EncodeLogRecordPos(pos *LogRecordPos) []byte {
	buf := make([]byte, binary.MaxVarintLen32+binary.MaxVarintLen64)
	var index = 0
	index += binary.PutVarint(buf[index:], int64(pos.Fid))
	index += binary.PutVarint(buf[index:], pos.Offset)
	return buf[:index]
}

// DecodeLogRecordPos 解码 LogRecordPos
func DecodeLogRecordPos(buf []byte) *LogRecordPos {
	var index = 0
	fileId, n := binary.Varint(buf[index:])
	index += n
	offset, _ := binary.Varint(buf[index:])
	return &LogRecordPos{Fid: uint32(fileId), Offset: offset}
}

// 对字节数组中的Header 信息进行解码
func decodeLogRecordHeader(buf []byte) (*logRecordHeader, int64) {
	if len(buf) <= 4 {
		return nil, 0
	}

	header := &logRecordHeader{
		crc:        binary.LittleEndian.Uint32(buf[:4]),
		recordType: buf[4],
	}

	var index = 5
	// 取出实际的 key size
	keySize, n := binary.Varint(buf[index:])
	header.keySize = uint32(keySize)
	index += n

	// 取出实际的 value size
	valueSize, n := binary.Varint(buf[index:])
	header.valueSize = uint32(valueSize)
	index += n

	return header, int64(index)
}

func getlogRecordCRC(lr *LogRecord, header []byte) uint32 {
	if lr == nil {
		return 0
	}

	crc := crc32.ChecksumIEEE(header[:])
	crc = crc32.Update(crc, crc32.IEEETable, lr.Key)
	crc = crc32.Update(crc, crc32.IEEETable, lr.Value)

	return crc
}
