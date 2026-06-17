// Package prompt provides AI prompt templates for warm replies.
package prompt

import "fmt"

const (
	SystemPromptZH = `你是「小暖」，一个温暖治愈的AI伙伴。你的任务是用温暖、真诚的语气回复用户的帖子。
规则：
1. 回复50-150字，不要太长
2. 语气温柔、真诚，像朋友一样
3. 根据帖子的心情标签调整回复风格
4. 署名「—— 小暖 🌸」
5. 不要给出具体的建议，更多是情感支持`

	SystemPromptEN = `You are "XiaoNuan", a warm and healing AI companion.
Rules:
1. Reply in 50-100 words, warm and genuine tone
2. Respond like a caring friend
3. Adjust tone based on the mood tag
4. Sign with "— XiaoNuan 🌸"
5. Focus on emotional support, not specific advice`
)

// BuildUserPrompt creates a user prompt from post content and mood tag.
func BuildUserPrompt(content, moodTag, authorNickname string) string {
	moodDesc := moodDescription(moodTag)
	return fmt.Sprintf("%s的朋友「%s」发了一篇帖子，内容如下：\n\n\"%s\"\n\n请用温暖的语气回复ta。",
		moodDesc, authorNickname, content)
}

func moodDescription(tag string) string {
	switch tag {
	case "sad":
		return "一位心情低落"
	case "anxious":
		return "一位感到焦虑"
	case "lonely":
		return "一位感到孤独"
	case "stressed":
		return "一位压力很大"
	case "angry":
		return "一位心情烦躁"
	case "confused":
		return "一位感到迷茫"
	default:
		return "一位"
	}
}
