package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const difficult = 20

//1.定义一个工作量证明的结构体对象
type ProofOfWork struct {
	// 数据来源
	block Block
	//目标难度值(整数)设定属性,可以使用big.Int这样一个内置的数据类型,方便后面调用它的一些比较的方法
	target *big.Int
}

//2.创建pow 函数,来实例化pow对象
func NewProofOfWork(block Block) *ProofOfWork {
	//传一个block的参数
	pow := ProofOfWork{
		block: block, //指定某些属性进行初始化 操作,可以不按顺序,无需全部都初始化,没有初始化的取默认值.
	}
	//初始化目标难度值,可一个得到1前导(difficult/4)-1个0
	//		a.初始化数值,末尾为1的int64
	targetInt := big.NewInt(1)
	//		b.二进制移位操作,总共256位,得到一个256/4=64 ,位的16进制数,因为2^4=16,4个二进制位=1个16进制位
	targetInt.Lsh(targetInt, 256-difficult)
	pow.target = targetInt
	return &pow
}

//封装一个方法,为hash计算做准备
func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	block := pow.block
	//定义并初始化二维切片,用于下一步拼接
	byteSlice := [][]byte{
		Uint2byte(block.Version),
		block.PrevBlockHash,
		block.MerkelRoot,
		Uint2byte(block.TimeStamp),
		Uint2byte(block.Difficulty),
		Uint2byte(nonce), //一个bug,如果不传过来,永远也挖不到矿
		block.Data,
	}
	info := bytes.Join(byteSlice, []byte{})
	return info
}

//3.为pow对象构建一个计算hash值的方法,返回两个参数,一个是当前的区块,一个是迭代运算得到的随机值nonce
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//TODO ,是一个关键字,用于增量式开发的一种提示,先预留,到一定时候再来实现
	fmt.Printf("pow run ...\n")
	//1.定义当前的hash
	var currentHash [32]byte
	var nonce uint64 //此处挖矿的随机数的初始值从0开始
	for {
		//获取挖矿的随机值
		//2.调用pow的一个封装的方法,返回字节流,从而下一步代入hash计算
		info := pow.prepareData(nonce)
		//计算hash值
		currentHash = sha256.Sum256(info)
		//比较hash值,判断是否满足条件
		var currentHashInt big.Int
		//为了调用Int的内置方法,方便比较,将获取的[]byte转换成big.Int类型
		currentHashInt.SetBytes(currentHash[:])
		//调用Int的内建函数,进行比较
		//fmt.Printf("currentHash:%v\n", currentHash)
		if currentHashInt.Cmp(pow.target) == -1 {
			// currentHashInt比目标值pow.target小,可以返回了
			fmt.Printf("找到哈希值了:%x,%d\n", currentHash, nonce)
			break
		} else {
			//比目标值大,继续挖矿,nonce++
			nonce++
		}

	}

	return currentHash[:], nonce
}

//4.验证找到的nonce是否正确
func (pow *ProofOfWork) IsValid() bool {
	//校验nonce
	info := pow.prepareData(pow.block.Nonce)
	//做hash运算
	hash := sha256.Sum256(info[:])
	//引入临时变量
	tmpInt := big.Int{}
	tmpInt.SetBytes(hash[:])
	//直接了当的使用return将比较结果返回

	return tmpInt.Cmp(pow.target) == -1
}
