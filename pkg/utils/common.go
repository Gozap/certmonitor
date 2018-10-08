package utils

import "github.com/Sirupsen/logrus"

func CheckErr(err error) bool {
	if err != nil {
		logrus.Error(err)
		return false
	} else {
		return true
	}
}
