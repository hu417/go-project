package common

//数组去重
func RemoveDuplicates(nums []uint) []uint {
	m := make(map[uint]bool)
	uniqueNums := []uint{}

	for _, num := range nums {
		if !m[num] {
			m[num] = true
			uniqueNums = append(uniqueNums, num)
		}
	}

	return uniqueNums
}