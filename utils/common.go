package utils

import "github.com/sirupsen/logrus"

func CheckErr(err error) bool {
	if err != nil {
		logrus.Error(err)
		return false
	} else {
		return true
	}
}
