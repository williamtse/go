package convert

import (
	"fmt"
	"strconv"
)

func StrToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return 0, err
	}
	return num, nil
}
