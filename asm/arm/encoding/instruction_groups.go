package encoding

type GroupType uint8

const (
	GT_Reserved        GroupType = 0b0000
	GT_Unallocated1    GroupType = 0b0001
	GT_SVE             GroupType = 0b0010
	GT_Unallocated2    GroupType = 0b0011
	GT_Data_Immediate1 GroupType = 0b1000
	GT_Data_Immediate2 GroupType = 0b1001
	GT_Load_Store1     GroupType = 0b0100
	GT_Load_Store2     GroupType = 0b0110
	GT_Load_Store3     GroupType = 0b1100
	GT_Load_Store4     GroupType = 0b1110
	GT_Data1           GroupType = 0b0101
	GT_Data2           GroupType = 0b1101
	GT_Data_Scalar1    GroupType = 0b0111
	GT_Data_Scalar2    GroupType = 0b1111
)
