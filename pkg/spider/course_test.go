package spider

import (
	"NeuBot/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCourse(t *testing.T) {
	user := &model.User{
		StdNumber: "",
		Password:  "",
	}
	couses, err := GetCourse(user)
	for i := range couses {
		fmt.Println(couses[i])
	}
	assert.NoError(t, err)
}
