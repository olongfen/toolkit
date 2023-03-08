package xerror

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	// ProjectNotBelongToPlatform 1
	ProjectNotBelongToPlatform = 40007
	// GenPDFFailed 生成pdf失败
	GenPDFFailed = 40008
)

func TestSetDefault(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()
	DefaultErrorMul.Set(ProjectNotBelongToPlatform, "zh", "项目不属于该平台")
	DefaultErrorMul.Set(ProjectNotBelongToPlatform, "en", "item does not belong to this platform")
	DefaultErrorMul.Set(GenPDFFailed, "zh", "生成pdf失败")
	DefaultErrorMul.Set(GenPDFFailed, "en", "gen pdf failed")
	assert.Equal(t, NewError(GenPDFFailed, "en").Error(), "gen pdf failed")
}
