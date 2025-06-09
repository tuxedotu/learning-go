package lasagna

const OvenTime = 20

func RemainingOvenTime(actual int) int {
	return OvenTime - actual
}

func PreperationTime(numberOfLayers int) int {
	return numberOfLayers * 2
}
