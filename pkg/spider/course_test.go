package spider

import (
	"NeuBot/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCourse(t *testing.T) {
	user := &model.User{
		StdNumber: "20206759",
		Password:  "xgs583719992",
	}
	couses, err := GetCourse(user)
	for i := range couses {
		fmt.Println(couses[i])
	}
	assert.NoError(t, err)
}
