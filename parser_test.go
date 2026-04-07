// Copyright 2021-2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package budoux

import (
	"fmt"
	"slices"
	"testing"

	"github.com/soundkitchen/go-budoux/models"
)

// tests for parser with default japanese model.
func TestDefaultJapaneseParser(t *testing.T) {
	cases := []struct {
		Sentence string
		Expected []string
	}{
		{
			Sentence: "Google の使命は、世界中の情報を整理し、世界中の人がアクセスできて使えるようにすることです。",
			Expected: []string{
				"Google の",
				"使命は、",
				"世界中の",
				"情報を",
				"整理し、",
				"世界中の",
				"人が",
				"アクセスできて",
				"使えるように",
				"する",
				"ことです。",
			},
		},
	}
	p := NewDefaultJapaneseParser()
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actual := p.Parse(c.Sentence)
			if !slices.Equal(actual, c.Expected) {
				t.Errorf("Expected %v, but got %v", c.Expected, actual)
			}
		})
	}
}

func TestParserWithCustomModel(t *testing.T) {
	t.Run("split before a", func(t *testing.T) {
		p := New(models.Model{
			"UW4": {
				"a": 10000,
			},
		})
		actual := p.Parse("abcdeabcd")
		expected := []string{"abcde", "abcd"}
		if !slices.Equal(actual, expected) {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	})

	t.Run("split before b", func(t *testing.T) {
		p := New(models.Model{
			"UW4": {
				"b": 10000,
			},
		})
		actual := p.Parse("abcdeabcd")
		expected := []string{"a", "bcdea", "bcd"}
		if !slices.Equal(actual, expected) {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		p := New(models.Model{})
		actual := p.Parse("")
		expected := []string{}
		if !slices.Equal(actual, expected) {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	})

	t.Run("odd total score keeps the same boundary decision", func(t *testing.T) {
		p := New(models.Model{
			"UW4": {
				"b": 1,
			},
		})
		actual := p.Parse("ab")
		expected := []string{"a", "b"}
		if !slices.Equal(actual, expected) {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	})
}

// tests for parser with japanese KNBC base model.
func TestJapaneseKNBCParser(t *testing.T) {
	cases := []struct {
		Sentence string
		Expected []string
	}{
		{
			Sentence: "気に入っている本をもう一度読んだ。",
			Expected: []string{
				"気に",
				"入っている",
				"本を",
				"もう",
				"一度",
				"読んだ。",
			},
		},
	}
	p := NewJapaneseKNBCParser()
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actual := p.Parse(c.Sentence)
			if !slices.Equal(actual, c.Expected) {
				t.Errorf("Expected %v, but got %v", c.Expected, actual)
			}
		})
	}
}

// tests for parser with default thai model.
func TestDefaultThaiParser(t *testing.T) {
	cases := []struct {
		Sentence string
		Expected []string
	}{
		{
			Sentence: "วันนี้อากาศดี",
			Expected: []string{
				"วัน",
				"นี้",
				"อากาศ",
				"ดี",
			},
		},
		{
			Sentence: "ภารกิจของ Google คือการจัดระเบียบข้อมูลของโลก และทำให้ข้อมูลนั้นสามารถเข้าถึงและใช้งานได้สำหรับทุกคนทั่วโลก",
			Expected: []string{
				"ภาร",
				"กิจ",
				"ของ",
				" ",
				"Google",
				" ",
				"คือ",
				"การ",
				"จัดระเบียบ",
				"ข้อมูล",
				"ของ",
				"โลก",
				" ",
				"และ",
				"ทำ",
				"ให้",
				"ข้อมูลนั้น",
				"สามารถ",
				"เข้า",
				"ถึง",
				"และ",
				"ใช้",
				"งาน",
				"ได้",
				"สำหรับ",
				"ทุก",
				"คน",
				"ทั่ว",
				"โลก",
			},
		},
		{
			Sentence: "ฉันชอบอ่านหนังสือในตอนเช้า",
			Expected: []string{
				"ฉัน",
				"ชอบ",
				"อ่าน",
				"หนัง",
				"สือ",
				"ใน",
				"ตอน",
				"เช้า",
			},
		},
	}
	p := NewDefaultThaiParser()
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actual := p.Parse(c.Sentence)
			if !slices.Equal(actual, c.Expected) {
				t.Errorf("Expected %v, but got %v", c.Expected, actual)
			}
		})
	}
}

// tests for parser with default simplified chinese model.
func TestDefaultSimplifiedChineseParser(t *testing.T) {
	cases := []struct {
		Sentence string
		Expected []string
	}{
		{
			Sentence: "我们的使命是整合全球信息，供大众使用，让人人受益。",
			Expected: []string{
				"我们",
				"的",
				"使命",
				"是",
				"整合",
				"全球",
				"信息，",
				"供",
				"大众",
				"使用，",
				"让",
				"人",
				"人",
				"受益。",
			},
		},
	}
	p := NewDefaultSimplifiedChineseParser()
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actual := p.Parse(c.Sentence)
			if !slices.Equal(actual, c.Expected) {
				t.Errorf("Expected %v, but got %v", c.Expected, actual)
			}
		})
	}
}

// tests for parser with default traditional chinese model.
func TestDefaultTraditionalChineseParser(t *testing.T) {
	cases := []struct {
		Sentence string
		Expected []string
	}{
		{
			Sentence: "我們的使命是匯整全球資訊，供大眾使用，使人人受惠。",
			Expected: []string{
				"我們",
				"的",
				"使命",
				"是",
				"匯整",
				"全球",
				"資訊，",
				"供",
				"大眾",
				"使用，",
				"使",
				"人",
				"人",
				"受惠。",
			},
		},
	}
	p := NewDefaultTraditionalChineseParser()
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actual := p.Parse(c.Sentence)
			if !slices.Equal(actual, c.Expected) {
				t.Errorf("Expected %v, but got %v", c.Expected, actual)
			}
		})
	}
}

func BenchmarkParserParse(b *testing.B) {
	b.Helper()
	p := NewDefaultJapaneseParser()
	sentence := "Google の使命は、世界中の情報を整理し、世界中の人がアクセスできて使えるようにすることです。"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = p.Parse(sentence)
	}
}
