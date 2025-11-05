package main

import (
	"reflect"
	"testing"
)

func TestUniq(t *testing.T) {
	tests := []struct {
		name   string
		lines  []string
		opts   Options
		expect []string
	}{
		{
			name:   "Без параметров",
			lines:  []string{"a", "a", "b", "b", "c"},
			opts:   Options{},
			expect: []string{"a", "b", "c"},
		},
		{
			name:   "Подсчет повторов",
			lines:  []string{"a", "a", "b"},
			opts:   Options{Count: true},
			expect: []string{"2 a", "1 b"},
		},
		{
			name:   "Только дубликаты",
			lines:  []string{"a", "a", "b", "c", "c"},
			opts:   Options{Duplicates: true},
			expect: []string{"a", "c"},
		},
		{
			name:   "Только уникальные",
			lines:  []string{"a", "a", "b", "c"},
			opts:   Options{Unique: true},
			expect: []string{"b", "c"},
		},
		{
			name:   "Игнор регистра",
			lines:  []string{"Hello", "hello", "World"},
			opts:   Options{IgnoreCase: true},
			expect: []string{"Hello", "World"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Uniq(tt.lines, tt.opts)
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("got %v, want %v", got, tt.expect)
			}
		})
	}
}
