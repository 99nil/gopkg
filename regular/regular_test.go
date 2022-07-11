// Created by zc on 2022/7/11.

package regular

import "testing"

func TestCheckTime(t *testing.T) {
	type args struct {
		startHour     int
		startMinute   int
		endHour       int
		endMinute     int
		currentHour   int
		currentMinute int
	}
	tests := []struct {
		name      string
		args      args
		wantStart bool
		wantEnd   bool
	}{
		{
			name: "all day",
			args: args{
				startHour:     0,
				startMinute:   0,
				endHour:       0,
				endMinute:     0,
				currentHour:   0,
				currentMinute: 0,
			},
			wantStart: true,
			wantEnd:   false,
		},
		{
			name: "all day2",
			args: args{
				startHour:     10,
				startMinute:   10,
				endHour:       10,
				endMinute:     10,
				currentHour:   0,
				currentMinute: 0,
			},
			wantStart: true,
			wantEnd:   false,
		},
		{
			name: "start",
			args: args{
				startHour:     10,
				startMinute:   0,
				endHour:       11,
				endMinute:     0,
				currentHour:   10,
				currentMinute: 10,
			},
			wantStart: true,
			wantEnd:   false,
		},
		{
			name: "equal start",
			args: args{
				startHour:     10,
				startMinute:   0,
				endHour:       11,
				endMinute:     0,
				currentHour:   10,
				currentMinute: 0,
			},
			wantStart: true,
			wantEnd:   false,
		},
		{
			name: "end",
			args: args{
				startHour:     10,
				startMinute:   0,
				endHour:       11,
				endMinute:     0,
				currentHour:   11,
				currentMinute: 10,
			},
			wantStart: true,
			wantEnd:   true,
		},
		{
			name: "equal end",
			args: args{
				startHour:     10,
				startMinute:   0,
				endHour:       11,
				endMinute:     0,
				currentHour:   11,
				currentMinute: 0,
			},
			wantStart: true,
			wantEnd:   true,
		},
		{
			name: "no start",
			args: args{
				startHour:     10,
				startMinute:   0,
				endHour:       11,
				endMinute:     0,
				currentHour:   9,
				currentMinute: 10,
			},
			wantStart: false,
			wantEnd:   false,
		},
		{
			name: "start > end, start",
			args: args{
				startHour:     10,
				startMinute:   40,
				endHour:       10,
				endMinute:     30,
				currentHour:   10,
				currentMinute: 50,
			},
			wantStart: true,
			wantEnd:   false,
		},
		{
			name: "start > end, no start",
			args: args{
				startHour:     10,
				startMinute:   40,
				endHour:       10,
				endMinute:     30,
				currentHour:   10,
				currentMinute: 35,
			},
			wantStart: false,
			wantEnd:   true,
		},
		{
			name: "start > end, no end",
			args: args{
				startHour:     10,
				startMinute:   40,
				endHour:       10,
				endMinute:     30,
				currentHour:   10,
				currentMinute: 20,
			},
			wantStart: true,
			wantEnd:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := CheckTime(tt.args.startHour, tt.args.startMinute, tt.args.endHour, tt.args.endMinute, tt.args.currentHour, tt.args.currentMinute)
			if gotStart != tt.wantStart {
				t.Errorf("CheckTime() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("CheckTime() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}
