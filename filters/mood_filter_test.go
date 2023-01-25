package filters

import (
	"testing"
	"zundafilter/zunda_mecab"
)

func TestMoodFilter(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		expect string
	}{
		{
			name:   "断定-φ:僕は腹筋できる人",
			text:   "僕は腹筋できる人",
			expect: "僕は腹筋できる人なのだ",
		},
		{
			name:   "断定+だ-φ:僕は腹筋できる人だ",
			text:   "僕は腹筋できる人だ",
			expect: "僕は腹筋できる人なのだ",
		},
		{
			name:   "断定-φ:「僕は腹筋できる人」",
			text:   "「僕は腹筋できる人」",
			expect: "「僕は腹筋できる人なのだ」",
		},
		{
			name:   "断定+だ-φ:「僕は腹筋できる人だ」",
			text:   "「僕は腹筋できる人だ」",
			expect: "「僕は腹筋できる人なのだ」",
		},
		{
			name:   "意志-φ",
			text:   "僕は腹筋する",
			expect: "僕は腹筋するのだ",
		},
		{
			name:   "推量-らしい",
			text:   "僕は腹筋するらしい",
			expect: "僕は腹筋するらしいのだ",
		},
		{
			name:   "確信-はずだ",
			text:   "僕は腹筋するはずだ",
			expect: "僕は腹筋するはずなのだ",
		},
		{
			name:   "説明-のだ",
			text:   "僕は腹筋できるのだ",
			expect: "僕は腹筋できるのなのだ",
		},
		{
			name:   "非断定-と思う",
			text:   "僕は腹筋できると思う",
			expect: "僕は腹筋できると思うのだ",
		},
		{
			name:   "可能性-かもしれない",
			text:   "僕は腹筋できるかもしれない",
			expect: "僕は腹筋できるかもしれないのだ",
		},
		{
			name:   "意志-つもりだ",
			text:   "僕は腹筋するつもりだ",
			expect: "僕は腹筋するつもりなのだ",
		},
		{
			name:   "願望-したい",
			text:   "僕は腹筋したい",
			expect: "僕は腹筋したいのだ",
		},
		{
			name:   "勧誘-しましょう",
			text:   "一緒に腹筋しましょう",
			expect: "一緒に腹筋しましょうなのだ",
		},
		{
			name:   "命令-なさい",
			text:   "一緒に腹筋しなさい",
			expect: "一緒に腹筋しなさいなのだ",
		},
		{
			name:   "許可-てもよい",
			text:   "一緒に腹筋してもよい",
			expect: "一緒に腹筋してもよいのだ",
		},
		{
			name:   "禁止-てはいけない",
			text:   "一緒に腹筋してはいけない",
			expect: "一緒に腹筋してはいけないのだ",
		},
		{
			name:   "質問-か",
			text:   "一緒に腹筋したか",
			expect: "一緒に腹筋したのだ",
		},
		{
			name:   "確認-だろう",
			text:   "昨日、一緒に腹筋しただろう",
			expect: "昨日、一緒に腹筋しただろうなのだ",
		},
		{
			name:   "同意-ね",
			text:   "一緒に腹筋したいね",
			expect: "一緒に腹筋したいね",
		},
		{
			name:   "断定-φ+記号",
			text:   "僕は腹筋できる人だ。",
			expect: "僕は腹筋できる人なのだ。",
		},
		{
			name:   "意志-φ+記号",
			text:   "僕は腹筋する。",
			expect: "僕は腹筋するのだ。",
		},
		{
			name:   "推量-らしい+記号",
			text:   "僕は腹筋するらしい。",
			expect: "僕は腹筋するらしいのだ。",
		},
		{
			name:   "確信-はずだ+記号",
			text:   "僕は腹筋するはずだ。",
			expect: "僕は腹筋するはずなのだ。",
		},
		{
			name:   "説明-のだ+記号",
			text:   "僕は腹筋できるのだ。",
			expect: "僕は腹筋できるのなのだ。",
		},
		{
			name:   "非断定-と思う+記号",
			text:   "僕は腹筋できると思う。",
			expect: "僕は腹筋できると思うのだ。",
		},
		{
			name:   "可能性-かもしれない+記号",
			text:   "僕は腹筋できるかもしれない。",
			expect: "僕は腹筋できるかもしれないのだ。",
		},
		{
			name:   "意志-つもりだ+記号",
			text:   "僕は腹筋するつもりだ。",
			expect: "僕は腹筋するつもりなのだ。",
		},
		{
			name:   "願望-したい+記号",
			text:   "僕は腹筋したい。",
			expect: "僕は腹筋したいのだ。",
		},
		{
			name:   "勧誘-しましょう+記号",
			text:   "一緒に腹筋しましょう。",
			expect: "一緒に腹筋しましょうなのだ。",
		},
		{
			name:   "命令-なさい+記号",
			text:   "一緒に腹筋しなさい。",
			expect: "一緒に腹筋しなさいなのだ。",
		},
		{
			name:   "許可-てもよい+記号",
			text:   "一緒に腹筋してもよい。",
			expect: "一緒に腹筋してもよいのだ。",
		},
		{
			name:   "禁止-てはいけない+記号",
			text:   "一緒に腹筋してはいけない。",
			expect: "一緒に腹筋してはいけないのだ。",
		},
		{
			name:   "質問-か+記号",
			text:   "一緒に腹筋したか？",
			expect: "一緒に腹筋したのだ？",
		},
		{
			name:   "確認-だろう+記号",
			text:   "昨日、一緒に腹筋しただろう。",
			expect: "昨日、一緒に腹筋しただろうなのだ。",
		},
		{
			name:   "依頼-てください",
			text:   "一緒に腹筋してください",
			expect: "一緒に腹筋してくださいなのだ",
		},
		{
			name:   "同意-ね+記号",
			text:   "一緒に腹筋したいね。",
			expect: "一緒に腹筋したいね。",
		},
		{
			name:   "依頼-てください+記号",
			text:   "一緒に腹筋してください。",
			expect: "一緒に腹筋してくださいなのだ。",
		},
		{
			name:   "命令-体言止め-末尾記号あり",
			text:   "一緒に来い！！",
			expect: "一緒に来いなのだ！！",
		},
		{
			name:   "命令-体言止め-末尾記号無し",
			text:   "一緒に来い",
			expect: "一緒に来いなのだ",
		},
		{
			name:   "断定-口頭",
			text:   "一緒に来るんだ",
			expect: "一緒に来るのだ",
		},
		{
			name:   "断定-口頭+末尾記号あり",
			text:   "一緒に来るんだ？",
			expect: "一緒に来るのだ？",
		},
		{
			name:   "質問(意志)-のか",
			text:   "一緒に来るのか",
			expect: "一緒に来るのだ",
		},
		{
			name:   "質問(意志)-のか+末尾記号あり",
			text:   "一緒に来るのか？",
			expect: "一緒に来るのだ？",
		},
		{
			name:   "質問(意志)-の+末尾記号なし",
			text:   "一緒に来るの",
			expect: "一緒に来るのだ",
		},
		{
			name:   "質問(意志)-の+末尾記号あり",
			text:   "一緒に来るの？",
			expect: "一緒に来るのだ？",
		},
		{
			name:   "動詞分-過去-た",
			text:   "一緒に来た？",
			expect:   "一緒に来たのだ？",
		},
		{
			name:   "動詞分-過去-た+末尾記号あり",
			text:   "一緒に来た",
			expect:   "一緒に来たのだ",
		},
		// ナイ形容詞語幹
		{
			name:   "ナイ形容詞語幹",
			text:   "それはしょうがない",
			expect:   "それはしょうがないのだ",
		},
		{
			name:   "ナイ形容詞語幹+末尾記号あり",
			text:   "それはしょうがない！",
			expect:   "それはしょうがないのだ！",
		},
		{
			name:   "不安の「の」:それはいいの",
			text:   "それはいいの",
			expect:   "それはいいのだ",
		},
		{
			name:   "不安の「の」:「それはいいの」",
			text:   "「それはいいの」",
			expect:   "「それはいいのだ」",
		},
	}

	mecabWrapper := zunda_mecab.MecabWrapper{
		Logger: getTestLogger(),
	}
	filter := MoodFilter{
		MecabWrapper: &mecabWrapper,
		Logger:       getTestLogger(),
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			// t.Parallel()
			actual, _ := filter.Convert(testCase.text)
			if actual != testCase.expect {
				t.Fatalf("MoodFilter.Convert() = %v, expect %v", actual, testCase.expect)
			}
		})
	}
}
