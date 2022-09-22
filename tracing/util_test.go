package tracing

import "testing"

func TestBuildShareUrl(t *testing.T) {

	tests := []struct {
		name    string
		url     string
		account string
		want    string
	}{
		{
			"test1",
			"",
			"xxx",
			"",
		},
		{
			"test2",
			"ddd",
			"W00000000514",
			"ddd?_acc_=465b7f1355597d9a32889191ee64a4fb",
		},
		{
			"test2.1",
			"ddd",
			"W00000000512",
			"ddd?_acc_=57999f24c5055fc8791730eedf36ad78",
		},
		{
			"test3",
			"ddd?x=a",
			"W00000000514",
			"ddd?x=a&_acc_=465b7f1355597d9a32889191ee64a4fb",
		},
		{
			"test4",
			"ddd?x=a#444",
			"W00000000514",
			"ddd?x=a&_acc_=465b7f1355597d9a32889191ee64a4fb#444",
		},
		{
			"test5",
			"ddd?x=a#",
			"W00000000514",
			"ddd?x=a&_acc_=465b7f1355597d9a32889191ee64a4fb#",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildShareUrl(tt.url, tt.account); got != tt.want {
				t.Errorf("BuildShareUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
