package semver

import (
  "fmt"
)

type Comparator struct {
  Value    string
  Operator string
  Version  string
  Semver   *Semver
}

func NewComparator(comp string) *Comparator {
  Debug("comparator %s\n", comp)

  comparator := Comparator{}
  comparator.Parse(comp)
  comparator.Value = comparator.Operator + comparator.Semver.Version

  return &comparator
}

func (comparator *Comparator) Parse(comp string) error {
  matches := RE_COMPARATOR.FindStringSubmatch(comp)

  if matches == nil {
    Debug("Invalid comparator: %s\n", comp)
    return fmt.Errorf("Invalid comparator: %s", comp)
  }

  var err error

  comparator.Operator = matches[1]
  comparator.Semver, err = NewSemver(matches[2])
  if err != nil {
    return err
  }

  // <1.2.3-rc DOES allow 1.2.3-beta (has prerelease)
  // >=1.2.3 DOES NOT allow 1.2.3-beta
  // <=1.2.3 DOES allow 1.2.3-beta
  // However, <1.2.3 does NOT allow 1.2.3-beta,
  // even though `1.2.3-beta < 1.2.3`
  // The assumption is that the 1.2.3 version has something you
  // *don't* want, so we push the prerelease down to the minimum.
  if comparator.Operator == "<" && !comparator.Semver.PreRelease.Empty {
    comparator.Semver.PreRelease = NewPreRelease("0")
    comparator.Semver.Format()
  }

  return nil
}

func (comparator *Comparator) Inspect() string {
  return fmt.Sprintf("%v", comparator)
}

func (comparator *Comparator) String() string {
  return comparator.Value
}

func (comparator *Comparator) Test(version string) (bool, error) {
  Debug("Comparator.test %s", version)

  return Cmp(version, comparator.Operator, comparator.Semver.Version)
}
