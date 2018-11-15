package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

//1.定义区块结构体对象
type Block struct {
	//区块头
	Version       uint64
	PrevBlockHash []byte
	MerkelRoot    []byte
	TimeStamp     uint64
	Difficulty    uint64
	Nonce         uint64
	Hash          []byte
	//区块体
	Data []byte
}

//2.创建新区块,实例化对象
func NewBlock(data string, prevHash []byte) *Block {
	//初始化
	block := Block{
		Version:       00,
		PrevBlockHash: prevHash,
		MerkelRoot:    []byte{},
		TimeStamp:     uint64(time.Now().Unix()),
		Difficulty:   difficult,
		Nonce:         0,        //暂时不挖矿,先随便赋值
		Hash:          []byte{}, //具体值,专门通过 SetHash方法来实现
		Data:          []byte(data),
	}
	//version1自己实现一个计算hash的方法,
	//block.SetHash()
	//version2,pow 工作量证明
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	//赋值操作
	block.Hash=hash
	block.Nonce=nonce
	return &block
}

// 3.构造一个函数,专门来计算Hash值
func (block *Block) SetHash() {
	//将区块中除了hash之外的所有属性拼接起来[]byte{},通过计算得到当前区块的hash
	blockSlice := [][]byte{
		//调用函数生成[]byte
		Uint2byte(block.Version),
		block.PrevBlockHash,
		block.MerkelRoot,
		Uint2byte(block.TimeStamp),
		Uint2byte(block.Difficulty),
		Uint2byte(block.Nonce),
		block.Data,
	}
	//拼接
	info:=bytes.Join(blockSlice,[]byte{})
	hash:=sha256.Sum256(info)//返回值为数组[32]byte
	//截取
	block.Hash=hash[:]

}
//uint64数据类型转换为[]byte
func Uint2byte(num uint64) []byte {
	var bufer bytes.Buffer
	//这是一个序列化的过程,将num转成buffer字节流
	err := binary.Write(&bufer, binary.BigEndian, &num)
	if err != nil {
		log.Panic(err)
	}

	return bufer.Bytes()

}
//对区块进行序列化,存储进入数据库(当交易通过网络传输或在应用程序之间交换时，它们被序列化)
func(block *Block)Serialize()[]byte{
	var buffer bytes.Buffer
	//定义编码器
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	if err!=nil{
		log.Panic(err)
	}

	return buffer.Bytes()
}

//将接受到的字节流转换成目标结构:反序列化
func Deserialize(data []byte)*Block{
	decoder := gob.NewDecoder(bytes.NewBuffer(data))
	var block Block
	//解码
	err := decoder.Decode(&block)
	if err!=nil{
		log.Panic(err)
	}

	return &block
}

