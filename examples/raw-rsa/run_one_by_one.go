package main

// check every block
func runOneByOne(fm *FileManager) bool {
	fm.StartSession()
	defer fm.EndSession()

	// check each block
	for i := 0; i < fm.blockCount; i++ {
		var data = fm.NextBlockData()
		var tag = fm.NextBlockTag()
		if data.Cmp(fm.rr.RawDecrypt(tag)) != 0 {
			return false
		}
	}
	return true
}
