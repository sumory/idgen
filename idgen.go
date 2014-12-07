package idgen

// `idgen` is actually the go implement of `snowflake`
// see https://github.com/sumory/uc/blob/master/src/com/sumory/uc/id/IdWorker.java for java implement

import (
	"errors"
	"fmt"
	"github.com/sumory/baseN4go"
	"sync"
	"time"
)

const (
	idgenEpoch         = int64(1417942588000) //2014-12-07 16:56:28
	workerIdBits       = uint(10)
	maxWorkerId        = -1 ^ (-1 << workerIdBits) //1023
	sequenceBits       = uint(12)
	workerIdShift      = sequenceBits                //12
	timestampLeftShift = sequenceBits + workerIdBits //22
	sequenceMask       = -1 ^ (-1 << sequenceBits)   //4095
)

/**
 * 42位的时间前缀+10位的节点标识+12位的sequence避免并发的数字（12位不够用时强制得到新的时间前缀）
 * 对系统时间的依赖性非常强，一旦开启使用，需要关闭ntp的时间同步功能，或者当检测到ntp时间调整后，拒绝分配id
 *
 * @author sumory.wu
 * @date 2014-12-07 16:56:28
 */
type IdWorker struct {
	sequence      int64
	lastTimestamp int64
	workerId      int64
	mutex         *sync.Mutex
	baseN4go      *baseN4go.BaseN //shorten the id
}

func NewIdWorker(workerId int64) (error, *IdWorker) {
	idWorker := &IdWorker{}
	if workerId > maxWorkerId || workerId < 0 {
		return errors.New(fmt.Sprintf("illegal worker id: %d", workerId)), nil
	}

	idWorker.workerId = workerId
	idWorker.lastTimestamp = -1
	idWorker.sequence = 0
	idWorker.mutex = &sync.Mutex{}
	err, baseN := baseN4go.NewBaseN(int8(62)) //默认radix为62，详见baseN4go
	if err != nil {
		return errors.New("can not initialize 'baseN4go'"), nil
	}
	idWorker.baseN4go = baseN
	return nil, idWorker
}

func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

// need synchronized
func (id *IdWorker) NextId() (error, int64) {
	id.mutex.Lock()
	defer id.mutex.Unlock()

	timestamp := timeGen()
	if timestamp < id.lastTimestamp {
		return errors.New(fmt.Sprintf("Clock moved backwards.Refusing to generate id for %d milliseconds", id.lastTimestamp-timestamp)), 0
	}
	if id.lastTimestamp == timestamp {
		id.sequence = (id.sequence + 1) & sequenceMask
		if id.sequence == 0 {
			timestamp = tilNextMillis(id.lastTimestamp)
		}
	} else {
		id.sequence = 0
	}
	id.lastTimestamp = timestamp
	return nil, ((timestamp - idgenEpoch) << timestampLeftShift) | (id.workerId << workerIdShift) | id.sequence
}

//func (id *IdWorker) Convert(source int64) string {
//	return strconv.FormatInt(source, 16)
//}
//
//func (id *IdWorker) ConvertWithRadix(source int64, radix int) string {
//	return strconv.FormatInt(source, radix)
//}

//重置用于shorten id的基数，参见baseN4go实现
func (id *IdWorker) RabaseShortRadix(radix int8) error {
	err, baseN := baseN4go.NewBaseN(radix)
	if err != nil {
		return err
	}
	id.baseN4go = baseN
	return nil
}

func (id *IdWorker) ShortId() (error, string) {
	err, newId := id.NextId()
	if err != nil {
		return err, ""
	}
	return id.baseN4go.Encode(newId)
}

func (id *IdWorker) ShortenId(genId int64) (error, string) {
	return id.baseN4go.Encode(genId)
}

//返回id是由哪个workerId生成的
func (id *IdWorker) WorkerId(genId int64) int64 {
	workerId := uint(uint(genId<<42)>>54)
	return int64(workerId)
}
