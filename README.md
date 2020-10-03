# rdiabv

具有双向验证能力的远程数据完整性审计(Remote Data Integrity Auditing with Bidirectional Verification)。

[English follows Chinese.](#English)

## 背景

此项目是我的[毕设论文](https://discretetom.github.io/academic/Thesis/research-on-remote-data-integrity-auditing-with-bidirectional-verification-in-cloud/)的一个POC，使用Golang实现，用来解决云存储上反向远程数据完整性验证的问题。

举个例子，我们可以把文件保存到百度网盘上，然后删除本地的文件，但是我们怎么知道文件被正确保存而没有被篡改呢？

一个方法是，把数据分块，并基于非对称加密，使用私钥加密每个数据块，得到数据块的【标签】，然后把数据和标签都存放在云上，并公开公钥。由于没有私钥，其他人无法伪造数据块的标签，但是所有人（包括数据持有者）都可以验证云上的数据块和标签是否匹配，从而验证云上的数据是否被篡改。正向数据完整性验证：用户可以验证云上保存的数据是否正确。

不过，云存储提供商也需要验证用户上传的数据块和标签是否匹配，防止用户诬陷自己“没有正确保存数据”。所以云存储提供商需要一个高效的验证手段来对数据块、标签进行验证。反向数据完整性验证：云存储提供商可以验证用户上传的数据是否正确。

此项目提供了两个方案：高次随机合并(HTRM)方案和动态高维防御(DHDD)方案。只要原本的非对称加密算法是同态的（比如RSA），这两个方案就可以通过合并数据块、标签来大幅度减少验证时间，并能以极高的准确率检测出合并数据块时可能出现的补偿攻击。

## 安装

```
go get github.com/DiscreteTom/rdiabv
```

## 使用方式

> 可以参考`examples`文件夹里面的示例程序进行学习

定义数据块。需要实现`Block`接口：

```go
type MyBlock struct {
	...
}

// 拷贝数据块。用来初始化合并缓冲区
func (block *MyBlock) Copy() rdiabv.Block {
	...
}

// 判断数据块是否合法
func (block *MyBlock) Validate() (ret bool) {
	...
}

// 定义如何合并两个数据块
func (block *MyBlock) Merge(x rdiabv.Block, y rdiabv.Block) rdiabv.Block {
	...
}
```

主程序工作流程：

```go
func main(){
	// 初始化算法
	dhdd := rdiabv.NewDHDD(blockCount, time.Now().UnixNano()).
		InitBuffers(&MyBlock{}) // 合并缓冲区会复制此数据块作为初始值

	// 合并所有数据块
	for i := 0; i < blockCount; i++ {
		block := MyBlock{}
		... // 构建数据块
		dhdd.MergeBlock(i, &block)
	}

	// 验证所有缓冲区中的数据域和标签域是否匹配
	dhdd.CheckAllBuffers()
}
```

## 性能

使用教科书RSA算法进行测试（源码位于`examples/dhdd-raw-rsa`），RSA密钥长度4096 bit，数据块大小256 Byte，使用`go run`命令进行测试，性能如下：

| 数据块数量 | 逐个验证耗时 | DHDD耗时 | DHDD维度 |
| --- | --- | --- | --- |
| 1024 | 44.6 s | 906 ms | 7 |
| 10000 | 7 m 56 s | 6.4 s | 9 |
| 59049 | 44 m 12 s | 20 s | 10 |

> 令DHDD维度为`x`，则DHDD防御失败的概率为`(5/27)**x`。维度为10时，DHDD防御失败的概率约等于中国国内中奖概率最低的大乐透中奖（约2142万分之一）。

# English

## Background

This project is a POC implementation of my [bachelor thesis](https://discretetom.github.io/academic/Thesis/research-on-remote-data-integrity-auditing-with-bidirectional-verification-in-cloud/). It's written with Golang, and used to solve the problem about reverse remote data integrity auditing.

For example, if we store our files on Google Drive and delete the local copy, how can we know whether our files were modified on the cloud?

One way to solve this problem is to split the files into many data blocks, and use asymmetric encryption methods to encrypt the data blocks to get the "tags" of each data block, then store both the data and tags on the cloud, and publish the public key. Anyone else can not generate tags for blocks since they do not have the private key, bu anyone(including the data owner) can validate whether the tags match the data, and know whether the files are modified. Forward data integrity auditing: users can validate whether the data on the cloud is correctly stored.

But the cloud storage provider also need to validate whether the data and tags are matched when users upload their files to prevent the users cheat them and upload mismatch data and tags. The cloud storage provider needs an efficient way to validate all the data and tags. Reverse data integrity auditing: the cloud storage provider can validate whether the data and tags are matched when users upload their files.

This project provides two methods to accelerate reverse data integrity auditing: High-Times Random Merging(HTRM) and Dynamic High-Dimensional Defense(DHDD). As long as the original encryption algorithm is homomorphic(for example, RSA), these two methods can reduce the validation time by merging data and tags, and detect complementation attack with a very high accuracy.

## Installation

```
go get github.com/DiscreteTom/rdiabv
```

## Usage

> You can always reference to the demo projects in the folder `examples`.

First we need to define our own block structure and implement the interface `Block`:

```go
type MyBlock struct {
	...
}

// Copy the current block. This is used to init DHDD buffers.
func (block *MyBlock) Copy() rdiabv.Block {
	...
}

// Check whether the data and tags of a block are matched.
func (block *MyBlock) Validate() (ret bool) {
	...
}

// Define how to merge two blocks.
func (block *MyBlock) Merge(x rdiabv.Block, y rdiabv.Block) rdiabv.Block {
	...
}
```

The workflow of the main process:

```go
func main(){
	// Initialize DHDD
	dhdd := rdiabv.NewDHDD(blockCount, time.Now().UnixNano()).
		InitBuffers(&MyBlock{}) // Buffers will copy the value of the parameter block.

	// Merge all blocks.
	for i := 0; i < blockCount; i++ {
		block := MyBlock{}
		... // construct your block
		dhdd.MergeBlock(i, &block)
	}

	// Check all buffers whether the data and tags are matched.
	dhdd.CheckAllBuffers()
}
```

## Performance

Using textbook RSA as the test algorithm(the source code is available in `examples/dhdd-raw-rsa`), with the length of the RSA key pair 4096, block data length 256 Bytes, using `go run` command to test. The results are as follows:

| Block Count | Time Consumption of Check One by One | Time Consumption of DHDD | Dimension of DHDD |
| --- | --- | --- | --- |
| 1024 | 44.6 s | 906 ms | 7 |
| 10000 | 7 m 56 s | 6.4 s | 9 |
| 59049 | 44 m 12 s | 20 s | 10 |

> Assume that the dimension of DHDD is `x`, then the probability of DHDD failed to detect the complementation attack is `(5/27)**x`. If `x == 10`, this probability are like 1/21420000.
