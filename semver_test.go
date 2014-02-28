package semver

import (
  "errors"
  "fmt"
  "testing"
)

func CompareExpect(a string, b string, expect int) error {
  compare, err := Compare(a, b)
  if err != nil {
    return err
  }

  lang := [...]string{"less than", "equal to", "greater than"}

  if compare == expect {
    return nil
  }

  return errors.New(fmt.Sprintf("Expected %s to be %s %s but it was %s", a, lang[expect+1], b, lang[compare+1]))
}

func TestInvalid(t *testing.T) {
  if _, err := NewSemver(""); err == nil {
    t.Error("Expected a blank string to be invalid")
  }

  if _, err := NewSemver("0.0"); err == nil {
    t.Error("Expected a missing patch to be invalid")
  }

  if _, err := NewSemver("0.0.01"); err == nil {
    t.Error("Expected a leadng zero to be invalid")
  }
}

func TestMajor(t *testing.T) {
  if err := CompareExpect("0.0.0", "0.0.0", 0); err != nil {
    t.Error(err)
  }

  if err := CompareExpect("1.0.0", "0.0.0", 1); err != nil {
    t.Error(err)
  }

  if err := CompareExpect("0.0.0", "1.0.0", -1); err != nil {
    t.Error(err)
  }
}

func TestMinor(t *testing.T) {
  if err := CompareExpect("0.1.0", "0.0.0", 1); err != nil {
    t.Error(err)
  }

  if err := CompareExpect("0.0.0", "0.1.0", -1); err != nil {
    t.Error(err)
  }

  if err := CompareExpect("0.11.0", "0.10.0", 1); err != nil {
    t.Error(err)
  }
}

func TestPreRelease(t *testing.T) {
  if err := CompareExpect("0.0.0-0.0.0", "0.0.0-0", 1); err != nil {
    t.Error(err)
  }

  if err := CompareExpect("0.0.0-0.0.0", "0.0.0-0.0", 1); err != nil {
    t.Error(err)
  }

  if err := CompareExpect("0.0.0-0.0.0", "0.0.0-0.0.alpha", -1); err != nil {
    t.Error(err)
  }
}

func TestIncrementVersion(t *testing.T) {
  base, _ := NewSemver("0.0.0")

  if err := base.Increment("major"); err != nil {
    t.Error(err)
  }

  major, _ := NewSemver("1.0.0")
  if base.Compare(major) != 0 {
    t.Errorf("Expected %s to equal %s", base, major)
  }

  if err := base.Increment("minor"); err != nil {
    t.Error(err)
  }

  minor, _ := NewSemver("1.1.0")
  if base.Compare(minor) != 0 {
    t.Errorf("Expected %s to equal %s", base, minor)
  }

  if err := base.Increment("patch"); err != nil {
    t.Error(err)
  }

  patch, _ := NewSemver("1.1.1")
  if base.Compare(patch) != 0 {
    t.Errorf("Expected %s to equal %s", base, patch)
  }
}

func TestIncrementPreRelease(t *testing.T) {
  base, _ := NewSemver("0.0.0")

  if err := base.Increment("prerelease"); err != nil {
    t.Error(err)
  }

  a, _ := NewSemver("0.0.0-0")
  if base.Compare(a) != 0 {
    t.Errorf("Expected %s to equal %s", base, a)
  }

  base, _ = NewSemver("0.0.0-0")

  if err := base.Increment("prerelease"); err != nil {
    t.Error(err)
  }

  b, _ := NewSemver("0.0.0-1")
  if base.Compare(b) != 0 {
    t.Errorf("Expected %s to equal %s", base, b)
  }

  base, _ = NewSemver("0.0.0-alpha")

  if err := base.Increment("prerelease"); err != nil {
    t.Error(err)
  }

  c, _ := NewSemver("0.0.0-alpha.0")
  if base.Compare(c) != 0 {
    t.Errorf("Expected %s to equal %s", base, c)
  }
}

func TestSort(t *testing.T) {
  versions := []string{
    "1.0.0",
    "0.1.0",
    "0.0.1",
    "0.0.0",
  }

  semvers := Semvers{}

  for _, version := range versions {
    semver, _ := NewSemver(version)
    semvers = append(semvers, *semver)
  }

  sorted := Sort(semvers)

  for x := 1; x < len(sorted); x++ {
    if sorted[x-1].Compare(&sorted[x]) == 1 {
      t.Errorf("Expected %s to be less than or equal to %s", &sorted[x-1], &sorted[x])
    }
  }
}

func TestComparator(t *testing.T) {
  _ = NewComparator(">0.0.0")
}

func TestRange(t *testing.T) {
  _ = NewRange("1.2.3 - 1.2.4")
}
