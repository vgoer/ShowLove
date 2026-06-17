package service

import (
	"context"
	"fmt"
	"time"

	"showlove/services/quote-service/internal/model"
	"showlove/services/quote-service/internal/repository"
)

type QuoteService struct{ repo repository.QuoteRepository }

func NewQuoteService(repo repository.QuoteRepository) *QuoteService {
	return &QuoteService{repo: repo}
}

func (s *QuoteService) GetTodayQuote(ctx context.Context) (*model.DailyQuote, error) {
	today := time.Now().Format("2006-01-02")
	return s.repo.FindByDate(ctx, today)
}

func (s *QuoteService) CreateQuote(ctx context.Context, textZH, textEN, author, date string) (*model.DailyQuote, error) {
	q := &model.DailyQuote{
		TextZH:        textZH,
		TextEN:        textEN,
		Author:        author,
		ScheduledDate: date,
	}
	if err := s.repo.Create(ctx, q); err != nil {
		return nil, fmt.Errorf("quote service: create: %w", err)
	}
	return q, nil
}

// SeedQuotes returns 30 bilingual healing quotes.
func SeedQuotes() []model.DailyQuote {
	today := time.Now()
	quotes := []struct {
		zh, en, author string
	}{
		{"万物皆有裂痕，那是光照进来的地方。", "There is a crack in everything. That's how the light gets in.", "Leonard Cohen"},
		{"每一个不曾起舞的日子，都是对生命的辜负。", "Every day without dancing is a betrayal of life.", "尼采"},
		{"生活明朗，万物可爱。", "Life is bright, everything is lovable.", "丰子恺"},
		{"你比你想象的要勇敢，比你看起来要坚强。", "You are braver than you believe, stronger than you seem.", "A.A. Milne"},
		{"今天所有的努力，都是为了明天能有更多选择的权利。", "All of today's efforts are for more choices tomorrow.", "佚名"},
		{"星光不问赶路人，时光不负有心人。", "The stars don't ask the traveler, time rewards the persistent.", "佚名"},
		{"世上只有一种英雄主义，就是在认清生活真相之后依然热爱生活。", "There is only one heroism: to see the world as it is and to love it.", "罗曼·罗兰"},
		{"愿你被这个世界温柔以待。", "May you be treated gently by the world.", "佚名"},
		{"所有的失去，都会以另一种方式归来。", "Everything you lose comes back in another form.", "Rumi"},
		{"不必太纠结于当下，也不必太忧虑未来。", "Don't dwell too much on the present, nor worry too much about the future.", "村上春树"},
		{"慢慢来，你又不差。", "Take your time, you're not behind.", "佚名"},
		{"如果生活给你柠檬，就做成柠檬水吧。", "When life gives you lemons, make lemonade.", "佚名"},
		{"愿你成为自己的太阳，无需凭借谁的光。", "May you become your own sun, without needing anyone else's light.", "佚名"},
		{"你笑起来真好看，像春天的花一样。", "Your smile is beautiful, like spring flowers.", "佚名"},
		{"人生没有白走的路，每一步都算数。", "No step in life is wasted; every step counts.", "佚名"},
		{"最暗的夜，才能看见最亮的星。", "Only in the darkest night can you see the brightest stars.", "佚名"},
		{"你要相信，所有的美好都会如期而至。", "Believe that all good things will come in time.", "佚名"},
		{"别忘了，你也是某人的光。", "Don't forget, you are someone's light too.", "佚名"},
		{"心之所向，素履以往。", "Where the heart goes, go barefoot.", "《礼记》"},
		{"你已经很棒了，休息一下也没关系。", "You're already doing great; it's okay to rest.", "小暖"},
		{"哪怕只有一点点进步，也是值得庆祝的。", "Even a little progress is worth celebrating.", "佚名"},
		{"世界很大，有人在偷偷爱你。", "The world is big, and someone secretly loves you.", "佚名"},
		{"温柔地对待自己，就像对待最好的朋友。", "Be gentle with yourself, like you would with your best friend.", "佚名"},
		{"风会记得每一朵花的香。", "The wind remembers the fragrance of every flower.", "佚名"},
		{"别怕，我在这儿呢。", "Don't be afraid, I'm right here.", "小暖"},
		{"生活的美好，往往藏在细微之处。", "The beauty of life is often hidden in the details.", "佚名"},
		{"雨过天晴，彩虹就会出现。", "After the rain, a rainbow appears.", "佚名"},
		{"你值得这世间所有的美好。", "You deserve all the good things in this world.", "佚名"},
		{"今天很好，明天会更好。", "Today is good, tomorrow will be even better.", "佚名"},
		{"做你自己，因为别人都有人做了。", "Be yourself; everyone else is already taken.", "Oscar Wilde"},
	}

	var result []model.DailyQuote
	for i, q := range quotes {
		date := today.AddDate(0, 0, i-30).Format("2006-01-02")
		result = append(result, model.DailyQuote{
			TextZH:        q.zh,
			TextEN:        q.en,
			Author:        q.author,
			ScheduledDate: date,
		})
	}
	return result
}
