package main

func main() {
	bc := NewBlockChain()
	bc.AddBlock("hello world")
	bc.AddBlock("这是第一笔交易")

	/*for i:=range bc.Blocks{

		fmt.Printf("=====================区块高度%d===================\n",i)
		fmt.Printf("Version:%d\n",bc.Blocks[i].Version)
		fmt.Printf("prevHash:%x\n",bc.Blocks[i].PrevBlockHash)
		fmt.Printf("Hash:%x\n",bc.Blocks[i].Hash)
		time:=time.Unix(int64(bc.Blocks[i].TimeStamp),0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp:%s\n",time)
		fmt.Printf("data:%s\n",string(bc.Blocks[i].Data))
		pow := NewProofOfWork(*bc.Blocks[i])
		fmt.Printf("IsValid:%v\n",pow.IsValid())
	}*/
}
