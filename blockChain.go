package main

import (
	"./bolt"
	"fmt"
	"log"
	"os"
) //手动导入

const genesisInfo = "2009年1月3日，财政大臣正处于实施第二轮银行紧急援助的边缘"
const BlockChainDb = "blockChain.db"
const blockBucket = "blockBucket"
const lastHashKey = "lastHashKey"

//4.创建区块链结构体
type BlockChain struct {
	//Blocks []*Block
	//操作数捷库的句柄
	db *bolt.DB
	//尾巴,用来存储最后一个区块的哈希
	tail []byte
}

//5.定义生成区块链的函数,返回一个对象的实例
func isDbExist() bool {
	if _, e := os.Stat(BlockChainDb); os.IsNotExist(e) {
		return false
	}
	return true

}

//创建一个新的区块链
func CreateBlockChain() *BlockChain {
	if isDbExist() {
		fmt.Println("CreateBlockChain区块链已经存在!\n")
		os.Exit(1)
	}

	/*1.打开数据库
		2.找到抽屉
		a.找到了
		b.没找到
			①创建抽屉
			②添加创世块
			③ 更新"last"的key

	*/
	var lastHash []byte
	db, err := bolt.Open(BlockChainDb, 0600, nil)
	if err != nil {
		fmt.Println("bolt.open failed!", err)
		os.Exit(1)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		//2. 找到我们的桶，通过桶的名字
		bucket := tx.Bucket([]byte(blockBucket))
		//如果没找到
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)

			}
			//写数据
			genesisBlock := NewBlock(genesisInfo, []byte{})
			err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)

			}
			//最后记得更新lastHashKey到数据库中
			err = bucket.Put([]byte(lastHashKey), genesisBlock.Hash)
			//更新内存中的最后的区块hash
			lastHash = genesisBlock.Hash
		}
		return nil
	})

	return &BlockChain{db, lastHash} //直接初始化
}

//实例化一个已有的区块链对象,得到它的实体(instance)
func NewBlockChain() *BlockChain {
	var blockChain *BlockChain

	if !isDbExist() {
		//fmt.Println("NewBlockChain请先创建区块链!")
		blockChain = CreateBlockChain()
	}
	var lastHash []byte
	db, err := bolt.Open(BlockChainDb, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.View(func(tx *bolt.Tx) error {
		//如果bucket不存在,那么就返回一个空
		bucket := tx.Bucket([]byte(blockBucket))
		//如果没有找到,就去创建bucket
		if bucket == nil {

			bucket, err := tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)
			}

			//3.写数据
			//在创建区块链的时候,添加一个创世块
			genesisBlock := NewBlock(genesisInfo, []byte{})
			err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化成字节流*/)
			if err != nil {
				log.Panic(err)
			}
			//一定要更新"lashHash"这个key------最后一个区块的哈希
			err = bucket.Put([]byte(lastHashKey), genesisBlock.Hash)
			//更新内存中的最后的区块哈希
			lastHash = genesisBlock.Hash
		} else {
			lastHash = bucket.Get([]byte(lastHash))
		}
		return nil
	})
	//return &BlockChain{db, lastHash}
	return blockChain
}

//6.为区块链对象构造一个方法,用来添加区块
func (blockChain *BlockChain) AddBlock(data string) {
	//通过切片索引,获取前一区块的hash值
	/*lastBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
	//将前一个区块的hash值赋值给新区块的PrevHash
	prevHash:=lastBlock.Hash
	newBlock:=NewBlock(data,prevHash)
	//为区块链对象的属性添加切片新成员
	blockChain.Blocks=append(blockChain.Blocks, newBlock)*/
	//最后一个区块的哈希值,也就是新区块的前哈希值
	preBlockHash := blockChain.tail

	// 更新数据库
	//1. 找到bucket
	//2. 判断有没有，
	//   有，写入数据
	//更新区块数据
	//更新lastHashKey对应的值
	//   没有， 直接报错退出
	//tx 为事务
	blockChain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Println("AddBlock添加区块时,bucket不能为空,请检查!")
			os.Exit(1)
		}
		newBlock := NewBlock(data, preBlockHash)
		//更新数据库
		bucket.Put(newBlock.Hash, newBlock.Serialize())
		bucket.Put([]byte(lastHashKey), newBlock.Hash)
		//更新内存
		blockChain.tail = newBlock.Hash
		return nil
	})

}

//自己封装一个迭代器,因为bolt自带的迭代器是按照key的字节大小的顺序进行迭代的,不方便,不是按照插入的顺序
/*func (bc *BlockChain)PrintChain(){

}*/
