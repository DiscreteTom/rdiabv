# rdiabv

具有双向验证能力的远程数据完整性审计(Remote Data Integrity Auditing with Bidirectional Verification)。

[English follows Chinese.](#English)

## 背景

此项目是我的[毕设论文](https://discretetom.github.io/academic/Thesis/research-on-remote-data-integrity-auditing-with-bidirectional-verification-in-cloud/)的一个POC，使用Golang实现，用来解决云存储上反向远程数据完整性验证的问题。

举个例子，我们可以把文件保存到百度网盘上，然后删除本地的文件，但是我们怎么知道文件被正确保存而没有被篡改呢？

一个方法是，把数据分块，并基于非对称加密，使用私钥加密每个数据块，得到数据块的【标签】，然后把数据和标签都存放在云上，并公开公钥。由于没有私钥，其他人无法伪造数据块的标签，但是所有人（包括数据持有者）都可以验证云上的数据块和标签是否匹配，从而验证云上的数据是否被篡改。正向数据完整性验证：用户可以验证云上保存的数据是否正确。

不过，云存储提供商也需要验证用户上传的数据块和标签是否匹配，防止用户诬陷自己“没有正确保存数据”。所以云存储提供商需要一个高效的验证手段来对数据块、标签进行验证。反向数据完整性验证：云存储提供商可以验证用户上传的数据是否正确。

此项目提供了两个方案：高次随机合并(HTRM)方案和动态高维防御(DHDD)方案。只要原本的非对称加密算法是同态的（比如RSA），这两个方案就可以通过合并数据块、标签来大幅度减少验证时间，并能以极高的准确率检测出合并数据块时可能出现的补偿攻击。论文中已经证明，这两个方案可以应用在基于RSA的S-PDP和基于BLS的POR上。

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
	copied := &MyBlock{}

	... // 从block复制值
	
	return copied
}

// 判断数据块是否合法
func (block *MyBlock) Validate() (ret bool) {
	...
}

// 定义如何合并两个数据块
func (block *MyBlock) Merge(x rdiabv.Block, y rdiabv.Block) rdiabv.Block {
	// 类型转换
	blockX, _ := x.(*MyBlock)
	blockY, _ := y.(*MyBlock)
	
	... // 合并blockX & blockY到block
	
	return block
}
```

主程序工作流程：

```go
func main(){
	// 使用DHDD
	validator := rdiabv.NewDHDD(blockCount, time.Now().UnixNano()).
		InitBuffers(&MyBlock{}) // 合并缓冲区会复制此数据块作为初始值
	
	// 使用HTRM
	// validator := rdiabv.NewHTRM(times).
		// InitBuffers(&MyBlock{}) // 合并缓冲区会复制此数据块作为初始值

	// 合并所有数据块
	for i := 0; i < blockCount; i++ {
		block := MyBlock{}
		... // 构建数据块
		validator.MergeBlock(i, &block) // 使用DHDD
		// validator.MergeBlock(&block) // 使用HTRM
	}

	// 验证所有缓冲区中的数据域和标签域是否匹配
	validator.CheckAllBuffers()
}
```

## 并发

DHDD和HTRM天然支持并发。`InitBuffers`和`CheckAllBuffers`显然是支持并发的，毕竟他们处理的是毫不相关的各个缓冲区。`MergeBlock`阶段，DHDD会操作位于不同维度的缓冲区，HTRM会操作不同次的缓冲区，所以一次`MergeBlock`调用中完全不涉及同时对一个相同缓冲区的写操作，因此单次`MergeBlock`在执行时，其内部是可以并行进行的。

因此，此项目在库的层面就支持并发（这也是此项目使用go语言进行编写的原因），使用`ParallelInitBuffers`/`ParallelMergeBlock`/`ParallelCheckAllBuffers`替换`InitBuffers`/`MergeBlock`/`CheckAllBuffers`即可使用并发操作。当然，并发与非并发API可以混合使用。

因为是库级别的并发，所以`ParallelXXX`内部会使用多个goroutine进行处理。而如果使用多个goroutine同时调用API（比如多个goroutine同时调用`MergeBlock`），才可能导致程序逻辑错误。当然，如果`MyBlock`被良好地设计，这些逻辑错误也可以被避免（比如`MyBlock`包含一个锁，`Merge`方法会检查这个锁）。

使用`ParallelXXX`时，goroutine的数量为DHDD的维度/HTRM的次数。

## 性能

使用教科书RSA算法进行测试（源码位于`examples/raw-rsa`），RSA密钥长度4096 bit，数据块大小256 Byte，使用`go run`命令进行测试，参考性能如下（由于随机生成数据不同，建议仅关注耗时之间的差距而不是具体数值）：

| 数据块数量 | 数据文件大小 | 逐个验证耗时 | DHDD耗时（非并发/并发） | DHDD维度 | HTRM耗时（非并发/并发） | HTRM次数 |
| --- |--- | --- | --- | --- | --- | --- |
| 1024 | 256 KiB | 30.30 s | 569.66 ms/224.0 ms | 7 | 1.54 s/420.84 ms | 15 |
| 10000 | 2500 KiB | 5 m 30 s | 3.16 s/1.46 s | 9 | 4.79 s/2.50 s | 15 |
| 59049 | 14.4 MiB | 27 m 11 s | 20.15 s/7.98 s | 10 | 22.44 s/9.72 s | 15 |

> 令DHDD维度为`x`，则DHDD防御失败的概率约为`(5/27)**x`。维度为10时，DHDD防御失败的概率约等于中国国内中奖概率最低的大乐透中奖（约1/21420000）。

> 令HTRM次数为`n`，则HTRM防御失败的概率为`1/3**n`。次数为15时，HTRM防御失败的概率约为1/14350000。

> 如果需要测试检错性能，只需要修改测试生成的`examples/raw-rsa/tag.txt`文件，然后修改`examples/raw-rsa/main.go`，使用已有文件而不是创建新文件即可。

# English

## Background

This project is a POC implementation of my [bachelor thesis](https://discretetom.github.io/academic/Thesis/research-on-remote-data-integrity-auditing-with-bidirectional-verification-in-cloud/). It's written with Golang, and used to solve the problem about reverse remote data integrity auditing.

For example, if we store our files on Google Drive and delete the local copy, how can we know whether our files were modified on the cloud?

One way to solve this problem is to split the files into many data blocks, and use asymmetric encryption methods to encrypt the data blocks to get the "tags" of each data block, then store both the data and tags on the cloud, and publish the public key. Anyone else can not generate tags for blocks since they do not have the private key, but anyone(including the data owner) can validate whether the tags match the data, and know whether the files are modified. Forward data integrity auditing: users can validate whether the data on the cloud is correctly stored.

But the cloud storage provider also need to validate whether the data and tags are matched when users upload their files to prevent the users cheat them and upload mismatch data and tags. The cloud storage provider needs an efficient way to validate all the data and tags. Reverse data integrity auditing: the cloud storage provider can validate whether the data and tags are matched when users upload their files.

This project provides two methods to accelerate reverse data integrity auditing: High-Times Random Merging(HTRM) and Dynamic High-Dimensional Defense(DHDD). As long as the original encryption algorithm is homomorphic(for example, RSA), these two methods can reduce the validation time by merging data and tags, and detect complementation attack with a very high accuracy. In the thesis paper we have proved that, these two methods can apply to RSA based S-PDP(Secure Provable Data Possession) and BLS based POR(Proofs of Retrievability).

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
	copied := &MyBlock{}
	
	... // copy fields from block

	return copied
}

// Check whether the data and tags of a block are matched.
func (block *MyBlock) Validate() (ret bool) {
	...
}

// Define how to merge two blocks.
func (block *MyBlock) Merge(x rdiabv.Block, y rdiabv.Block) rdiabv.Block {
	// type assertion
	blockX, _ := x.(*MyBlock)
	blockY, _ := y.(*MyBlock)

	... // merge blockX & blockY to block

	return block
}
```

The workflow of the main process:

```go
func main(){
	// Using DHDD
	validator := rdiabv.NewDHDD(blockCount, time.Now().UnixNano()).
		InitBuffers(&MyBlock{}) // Buffers will copy the value of the parameter block.
	
	// Using HTRM
	// validator := rdiabv.NewHTRM(times).
	// 	InitBuffers(&MyBlock{}) // Buffers will copy the value of the parameter block.

	// Merge all blocks.
	for i := 0; i < blockCount; i++ {
		block := MyBlock{}
		... // construct your block
		validator.MergeBlock(i, &block) // Using DHDD
		// validator.MergeBlock(&block) // Using HTRM
	}

	// Check all buffers whether the data and tags are matched.
	validator.CheckAllBuffers()
}
```

## Parallelism

DHDD & HTRM are born to support parallelism. `InitBuffers` & `CheckAllBuffers` are processing different buffers so these two methods can be implemented in parallel. In `MergeBlock` stage, DHDD will process buffers in different dimension, and HTRM will process buffers in different time, so there will be no possibility to write to a same buffer many times in one `MergeBlock` call. So `MergeBlock` can be implemented in parallel too.

Therefore, this project supports parallel processing(and that's why this project is written with Golang). You can just replace `InitBuffers`/`MergeBlock`/`CheckAllBuffers` with `ParallelInitBuffers`/`ParallelMergeBlock`/`ParallelCheckAllBuffers` to use parallel processing. You can also mix parallel APIs and non-parallel APIs.

There will be many goroutines in one `ParallelXXX` call. If you use many goroutines to call APIs in parallel(for example, use many goroutines to call `MergeBlock` simultaneously), there might be some logical problems. You can design your `MyBlock` properly to avoid these problems of course(for example, add a lock to the `MyBlock` structure and check this lock in `Merge` method).

The number of goroutines are the dimension of DHDD and the times of HTRM when using parallel APIs.

## Performance

Using textbook RSA as the test algorithm(the source code is available in `examples/raw-rsa`), with 4096 bits RSA key pair, each block contains 256 bytes data, using `go run` command to test. The results are as follows(since each test will generate different data, we should focus on the differences of time consumption between all tests instead of the value of time consumptions):

| Block Count | File Size | Time Consumption of Check One by One | Time Consumption of DHDD(Non-Parallel/Parallel) | Dimension of DHDD | Time Consumption of HTRM(Non-Parallel/Parallel) | HTRM Times |
| --- |--- | --- | --- | --- | --- | --- |
| 1024 | 256 KiB | 30.30 s | 569.66 ms/224.0 ms | 7 | 1.54 s/420.84 ms | 15 |
| 10000 | 2500 KiB | 5 m 30 s | 3.16 s/1.46 s | 9 | 4.79 s/2.50 s | 15 |
| 59049 | 14.4 MiB | 27 m 11 s | 20.15 s/7.98 s | 10 | 22.44 s/9.72 s | 15 |

> Assume that the dimension of DHDD is `x`, then the probability of DHDD failed to detect the complementation attack is about `(5/27)**x`. If `x == 10`, this probability is about `1/21420000`.

> Assume that the times of HTRM is `n`, then the probability of HTRM failed to detect the complementation attack is `1/3**n`. If `n == 15`, this probability is about `1/14350000`.

> If you want to test the ability of error detection, just modify the generated `examples/raw-rsa/tag.txt`, the modify `examples/raw-rsa/main.go` to use existing files instead creating new ones, then re-run the test.