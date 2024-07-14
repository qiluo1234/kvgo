package data

import (
	"bitcast-go/fio"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"path/filepath"
)

var (
	ErrInvalidCRC = errors.New("invalid crc value, log record maybe corrupted")
)

const DataFilesNameSuffix = ".data"

// DataFile 数据文件
type DataFile struct {
	FileId    uint32        //文件id
	WriteOff  int64         //文件写到了哪个位置
	IoManager fio.IOManager //io 读写管理
}

// OpenDataFile 打开新的数据文件
func OpenDataFile(dirPath string, fileID uint32) (*DataFile, error) {

	fileName := filepath.Join(dirPath, fmt.Sprintf("%09d", fileID)+DataFilesNameSuffix)
	ioManager, err := fio.NewIOManager(fileName)
	if err != nil {
		return nil, err
	}
	return &DataFile{
		FileId:    fileID,
		WriteOff:  0,
		IoManager: ioManager,
	}, nil
}

// 根据 offset 从数据文件中读取 LogRecord
func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, int64, error) {
	filesize, err := df.IoManager.Size()
	if err != nil {
		return nil, 0, err
	}

	//如果读取的最大的 header 长度已经超过了文件的长度，则只需要读取到文件的末尾即可
	var headerBytes int64 = maxLogRecordHeaderSize
	if offset+maxLogRecordHeaderSize > filesize {
		headerBytes = filesize - offset
	}

	//读取 Header信息
	headerBuf, err := df.readNBytes(headerBytes, offset)
	if err != nil {
		return nil, 0, err
	}

	header, headerSize := decodeLogRecordHeader(headerBuf)
	//下面的两个条件表示读取到了文件末尾，直接返回EOF错误
	if header == nil {
		return nil, 0, io.EOF
	}
	if header.crc == 0 && header.keySize == 0 && header.valueSize == 0 {
		return nil, 0, io.EOF
	}

	keySize, valueSize := int64(header.keySize), int64(header.valueSize)
	var recordSize = headerSize + keySize + valueSize

	logRecord := &LogRecord{Type: header.recordType}
	//开始读取用户实际存储的key/value 数据
	if keySize > 0 || valueSize > 0 {
		kvBuf, err := df.readNBytes(keySize+valueSize, offset+headerSize)
		if err != nil {
			return nil, 0, err
		}

		//解出 key和 value
		logRecord.Key = kvBuf[:keySize]
		logRecord.Value = kvBuf[keySize:]

	}

	//校验数据的有效性
	crc := getlogRecordCRC(logRecord, headerBuf[crc32.Size:headerSize])
	if crc != header.crc {
		return nil, 0, ErrInvalidCRC
	}
	return logRecord, recordSize, nil
}

// 去除对应的 key和 value 的长度

func (df *DataFile) Write(buf []byte) error {
	n, err := df.IoManager.Write(buf)
	if err != nil {
		return err
	}
	df.WriteOff += int64(n)
	return nil
}
func (df *DataFile) Sync() error {

	return df.IoManager.Sync()
}

func (df *DataFile) Close() error {
	return df.IoManager.Close()
}

func (df *DataFile) readNBytes(n int64, offset int64) (b []byte, err error) {

	b = make([]byte, n)
	_, err = df.IoManager.Read(b, offset)
	return
}
