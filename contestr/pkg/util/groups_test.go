package util

import (
	"reflect"
	"testing"
)

func TestFormGroupsWithSwapProbability_EmptyInput(t *testing.T) {
	got := FormGroupsWithSwapProbability([]string{}, 3, 0)

	if len(got) != 0 {
		t.Fatalf("expected empty result, got %#v", got)
	}
}

func TestFormGroupsWithSwapProbability_GroupSizeLessThanTwoUsesTwo(t *testing.T) {
	input := []string{"a", "b", "c"}

	got := FormGroupsWithSwapProbability(input, 1, 0)

	want := [][]string{
		{"a", "b"},
		{"c"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_NoSwapExactGroups(t *testing.T) {
	input := []string{"a", "b", "c", "d"}

	got := FormGroupsWithSwapProbability(input, 2, 0)

	want := [][]string{
		{"a", "b"},
		{"c", "d"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_NoSwapUnevenGroupsWithoutRebalance(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e"}

	got := FormGroupsWithSwapProbability(input, 3, 0)

	want := [][]string{
		{"a", "b", "c"},
		{"d", "e"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_NoSwapUnevenGroupsWithRebalance(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e"}

	got := FormGroupsWithSwapProbability(input, 4, 0)

	want := [][]string{
		{"a", "b", "c"},
		{"d", "e"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_Size3_Balanced(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e", "f"}

	got := FormGroupsWithSwapProbability(input, 3, 0)

	want := [][]string{
		{"a", "b", "c"},
		{"d", "e", "f"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_Size3_NBalanced1(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e", "f", "g"}

	got := FormGroupsWithSwapProbability(input, 3, 0)

	want := [][]string{
		{"a", "b", "c"},
		{"d", "e"},
		{"f", "g"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_Size3_NBalanced2(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e", "f", "g", "h"}

	got := FormGroupsWithSwapProbability(input, 3, 0)

	want := [][]string{
		{"a", "b", "c"},
		{"d", "e", "f"},
		{"g", "h"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_SeveralGroupsRebalanceOnlyLastTwo(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}

	got := FormGroupsWithSwapProbability(input, 4, 0)

	want := [][]string{
		{"a", "b", "c", "d"},
		{"e", "f", "g"},
		{"h", "i"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_SingleGroupDoesNotPanic(t *testing.T) {
	input := []string{"a", "b"}

	got := FormGroupsWithSwapProbability(input, 10, 0)

	want := [][]string{
		{"a", "b"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_DoesNotMutateInput(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e"}
	original := append([]string(nil), input...)

	_ = FormGroupsWithSwapProbability(input, 3, 0)

	if !reflect.DeepEqual(input, original) {
		t.Fatalf("input was mutated:\nwant %#v\ngot  %#v", original, input)
	}
}

func TestFormGroupsWithSwapProbability_ContainsSameParticipants(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e", "f", "g"}

	got := FormGroupsWithSwapProbability(input, 3, 0.5)

	if !sameMultiset(input, flatten(got)) {
		t.Fatalf("participants mismatch:\ninput %#v\ngot   %#v", input, got)
	}
}

func TestFormGroupsWithSwapProbability_NegativeProbabilityClampedToZero(t *testing.T) {
	input := []string{"a", "b", "c", "d"}

	got := FormGroupsWithSwapProbability(input, 2, -10)

	want := [][]string{
		{"a", "b"},
		{"c", "d"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func TestFormGroupsWithSwapProbability_ProbabilityGreaterThanOneClampedToOne(t *testing.T) {
	input := []string{"a", "b", "c", "d"}

	got := FormGroupsWithSwapProbability(input, 2, 10)

	want := [][]string{
		{"b", "c"},
		{"d", "a"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected groups:\nwant %#v\ngot  %#v", want, got)
	}
}

func flatten(groups [][]string) []string {
	res := make([]string, 0)

	for _, group := range groups {
		res = append(res, group...)
	}

	return res
}

func sameMultiset(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	cnt := make(map[string]int, len(a))

	for _, x := range a {
		cnt[x]++
	}

	for _, x := range b {
		cnt[x]--
		if cnt[x] < 0 {
			return false
		}
	}

	return true
}
