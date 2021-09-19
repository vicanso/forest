package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type SubData struct {
	SubTitle string `json:"subTitle"`
}

type Data struct {
	Title   string   `json:"title"`
	SubData *SubData `json:"subData"`
}

type MergeData struct {
	SubData
	Title string `json:"title"`
	Name  string `json:"name" default:"test"`
}

func TestValidateQuery(t *testing.T) {
	assert := assert.New(t)
	md := MergeData{}
	data := map[string]string{
		"subTitle": "s",
		"title":    "t",
	}
	err := Query(&md, data)
	assert.Nil(err)
	assert.Equal("s", md.SubData.SubTitle)
	assert.Equal("s", md.SubTitle)
	assert.Equal("t", md.Title)
	assert.Equal("test", md.Name)
}

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	md := MergeData{}
	data := map[string]string{
		"subTitle": "s",
		"title":    "t",
		"name":     "123",
	}
	err := Do(&md, data)
	assert.Nil(err)
	assert.Equal("s", md.SubData.SubTitle)
	assert.Equal("s", md.SubTitle)
	assert.Equal("t", md.Title)
	assert.Equal("123", md.Name)

	d := Data{}
	err = Do(&d, []byte(`{
		"title": "t",
		"subData": {
			"subTitle": "s"
		}
	}`))
	assert.Nil(err)
	assert.Equal("s", d.SubData.SubTitle)
	assert.Equal("t", d.Title)
}
