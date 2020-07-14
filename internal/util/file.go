package util

import "os"

// FileExists stops the program if the file does not exist
func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return true, err
	}

	return true, nil
}

// // EnsureFolderExists stops the program if the file does not exist
// func EnsureFolderExists(name string) {
// 	stat, err := os.Stat(name)

// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return false, err
// 		}
// 		return false, err
// 	}

// 	return stat.IsDir(), nil
// }
