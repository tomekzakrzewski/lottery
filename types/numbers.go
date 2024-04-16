package types

const (
	numLen = 6
	minNum = 1
	maxNum = 99
)

type UserNumbers struct {
	Numbers []int `json:"numbers"`
}

func (input *UserNumbers) ValidateNumbers() bool {
	if len(input.Numbers) != numLen {
		return false
	}
	for _, num := range input.Numbers {
		if num < minNum || num > maxNum {
			return false
		}
	}
	return input.unique()
}
func (input *UserNumbers) unique() bool {
	unique := make(map[int]bool)
	for _, num := range input.Numbers {
		unique[num] = true
	}
	return len(unique) == len(input.Numbers)
}
