package moderation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_Contains_SensitiveWord(t *testing.T) {
	f := NewFilter(DefaultChineseWords())

	assert.True(t, f.Contains("这是一个包含赌博的帖子"))
	assert.True(t, f.Contains("涉及色情内容"))
	assert.False(t, f.Contains("这是一篇温暖的帖子"))
	assert.False(t, f.Contains("今天天气真好"))
}

func TestFilter_FindAll(t *testing.T) {
	f := NewFilter(DefaultChineseWords())

	words := f.FindAll("这个帖子涉及赌博和诈骗")
	assert.Len(t, words, 2)
	assert.Contains(t, words, "赌博")
	assert.Contains(t, words, "诈骗")
}

func TestFilter_Empty(t *testing.T) {
	f := NewFilter(nil)

	assert.False(t, f.Contains("任何内容"))
	assert.Empty(t, f.FindAll("任何内容"))
}

func TestFilter_CustomWords(t *testing.T) {
	words := []string{"测试敏感词", "badword"}
	f := NewFilter(words)

	assert.True(t, f.Contains("这是测试敏感词哦"))
	assert.True(t, f.Contains("This is a BADWORD here"))
	assert.False(t, f.Contains("正常内容"))
}

func TestFilter_NoFalsePositive(t *testing.T) {
	f := NewFilter(DefaultChineseWords())

	// "博" alone should not trigger "赌博"
	assert.False(t, f.Contains("博物馆"))
	assert.False(t, f.Contains("彩色"))
	assert.False(t, f.Contains("力量"))
}
