package functions

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestToInt64(t *testing.T) {
	// String test
	i := ToInt64(1)
	logrus.WithFields(logrus.Fields{
		"Input":  "int: 1",
		"Output": i,
	}).Info()
	i = ToInt64("2")
	logrus.WithFields(logrus.Fields{
		"Input":  `String: 2`,
		"Output": i,
	}).Info()
	i = ToInt64(3.0)
	logrus.WithFields(logrus.Fields{
		"Input":  "Float: 3.0",
		"Output": i,
	}).Info()
	i = ToInt64("4.0")
	logrus.WithFields(logrus.Fields{
		"Input":  `Float (string): 4.0`,
		"Output": i,
	}).Info()
	i = ToInt64(5.1)
	logrus.WithFields(logrus.Fields{
		"Input":  "Float: 5.1",
		"Output": i,
	}).Info()
	i = ToInt64("6.1")
	logrus.WithFields(logrus.Fields{
		"Input":  `Float: (string) 6.1`,
		"Output": i,
	}).Info()
}
