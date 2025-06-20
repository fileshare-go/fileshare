package util

func MergeList(list1, list2 []int32) []int32 {
	result := []int32{}
	i, j := 0, 0

	for i < len(list1) && j < len(list2) {
		val := int32(0)
		if list1[i] < list2[j] {
			val = list1[i]
			i++
		} else if list1[i] > list2[j] {
			val = list2[j]
			j++
		} else {
			val = list1[i]
			i++
			j++
		}
		if len(result) == 0 || result[len(result)-1] != val {
			result = append(result, val)
		}
	}

	for i < len(list1) {
		if len(result) == 0 || result[len(result)-1] != list1[i] {
			result = append(result, list1[i])
		}
		i++
	}

	for j < len(list2) {
		if len(result) == 0 || result[len(result)-1] != list2[j] {
			result = append(result, list2[j])
		}
		j++
	}

	return result
}

func MissingElementsInSortedList(total []int32, subList []int32) []int32 {
	if len(subList) == len(total) || len(total) == 0 {
		return []int32{}
	}

	i := 0
	j := 0
	result := []int32{}

	for i < len(total) && j < len(subList) {
		if total[i] == subList[j] {
			i++
			j++
		} else if total[i] < subList[j] {
			result = append(result, total[i])
			i++
		} else {
			return []int32{}
		}
	}

	for i < len(total) {
		result = append(result, total[i])
		i++
	}

	return result
}
