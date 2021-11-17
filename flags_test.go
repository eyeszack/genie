package geenee

import (
	"flag"
	"testing"
)

func TestIntSlice_String(t *testing.T) {
	t.Run("validate the expected string value", func(t *testing.T) {
		want := "3"
		subject := IntSlice{3}
		got := subject.String()

		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate the expected string values", func(t *testing.T) {
		want := "3, 6, 10, 9"
		subject := IntSlice{3, 6, 10, 9}
		got := subject.String()

		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}

func TestIntSlice_Set(t *testing.T) {
	t.Run("validate the expected value is set", func(t *testing.T) {
		want := 3
		subject := IntSlice{}
		err := subject.Set("3")
		if err != nil {
			t.Fatalf("unexpected err, %s", err)
		}

		got := subject[0]
		if got != want {
			t.Errorf("want %d, got %d", want, got)
		}
	})

	t.Run("validate the expected values are set", func(t *testing.T) {
		want := []int{3, 6, 9, 10, 11}
		subject := IntSlice{}
		err := subject.Set("3")
		if err != nil {
			t.Fatalf("unexpected err, %s", err)
		}
		err = subject.Set("6, 9,10")
		if err != nil {
			t.Fatalf("unexpected err, %s", err)
		}
		err = subject.Set("11")
		if err != nil {
			t.Fatalf("unexpected err, %s", err)
		}

		if len(subject) != 5 {
			t.Errorf("want 5, got %d", len(subject))
		}

		for i, got := range subject {
			if got != want[i] {
				t.Errorf("want %d, got %d", want[i], got)
			}
		}
	})

	t.Run("validate error is returned when error encountered", func(t *testing.T) {
		want := []int{3, 6}
		subject := IntSlice{}
		err := subject.Set("3")
		if err != nil {
			t.Fatalf("unexpected err, %s", err)
		}
		err = subject.Set("6")
		if err != nil {
			t.Fatalf("unexpected err, %s", err)
		}

		err = subject.Set("heyo")
		if err == nil {
			t.Error("expected an err, got nil")
		}

		err = subject.Set("x,y,z")
		if err == nil {
			t.Error("expected an err, got nil")
		}

		if len(subject) != 2 {
			t.Errorf("want 5, got %d", len(subject))
		}

		for i, got := range subject {
			if got != want[i] {
				t.Errorf("want %d, got %d", want[i], got)
			}
		}
	})
}

func TestIntSlice_FlagParse(t *testing.T) {
	t.Run("validate the expected values are set", func(t *testing.T) {
		want := []int{3, 6, 9, 10, 11}
		subject := IntSlice{}

		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&subject, "subject", "the test subject")
		err := fs.Parse([]string{"--subject", "3", "--subject", "6,9,10", "--subject", "11"})
		if err != nil {
			t.Fatalf("unexpected err, %s", err)
		}

		if len(subject) != 5 {
			t.Errorf("want 5, got %d", len(subject))
		}

		for i, got := range subject {
			if got != want[i] {
				t.Errorf("want %d, got %d", want[i], got)
			}
		}
	})

	t.Run("validate error is returned when value is invalid", func(t *testing.T) {
		want := []int{3, 6}
		subject := IntSlice{}

		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&subject, "subject", "the test subject")
		err := fs.Parse([]string{"--subject", "3", "--subject", "6", "--subject", "dang", "--subject", "x,y,z"})
		if err == nil {
			t.Error("expected an err, got nil")
		}

		if len(subject) != 2 {
			t.Errorf("want 5, got %d", len(subject))
		}

		for i, got := range subject {
			if got != want[i] {
				t.Errorf("want %d, got %d", want[i], got)
			}
		}
	})
}
