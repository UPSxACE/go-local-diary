package utils;

func BoolToInt(x bool) int{
	if(x){
		return 1;
	}
	return 0;
}

func IntToBool(x int) bool{
	return x != 0
}