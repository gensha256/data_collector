package common

func DoesArrayContainString(arr []string, str string) bool {

	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return true
		}
	}

	return false
}

func SortArray(arr []CmcEntity) []CmcEntity {

	if len(arr) == 1 {
		return arr
	}

	moved := true
	pointer := 0

	for moved == true {

		moved = false

		for pointer < len(arr)-1 {

			first := arr[pointer]
			second := arr[pointer+1]

			if first.LastUpdated.Unix() > second.LastUpdated.Unix() {

				arr[pointer] = second
				arr[pointer+1] = first

				moved = true
			}

			pointer++
		}

		pointer = 0
	}

	return arr
}
