package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telekom-mms/corp-net-indicator/internal/util"
)

func TestFormatDate(t *testing.T) {
	assert.Equal(t, time.Unix(1039717812, 0).Local().Format(util.DATE_TIME_FORMAT), util.FormatDate(1039717812))
	assert.Equal(t, "-", util.FormatDate(0))
	assert.Equal(t, time.Unix(1, 0).Local().Format(util.DATE_TIME_FORMAT), util.FormatDate(int64(1)))
}

func TestFormatValue(t *testing.T) {
	assert.Equal(t, "-", util.FormatValue(""))
	assert.Equal(t, "value", util.FormatValue("value"))
}
