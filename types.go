package semver

import (
  "errors"
  "fmt"
  "strconv"
  "strings"
)

type Indentifier struct {
  Val     string
  IntVal  int
  Numeric bool
}

type PreRelease struct {
  Raw     string
  Pieces  []*Indentifier
  Empty   bool
  Version string
}

type Build struct {
  Raw     string
  Pieces  []*Indentifier
  Empty   bool
  Version string
}

type Semver struct {
  Raw        string
  Major      int
  Minor      int
  Patch      int
  PreRelease *PreRelease
  Build      *Build
  Version    string
}

type Semvers []Semver

func (this Semvers) Len() int {
  return len(this)
}

func (this Semvers) Less(i, j int) bool {
  return this[i].Compare(&this[j]) == -1
}

func (this Semvers) Swap(i, j int) {
  this[i], this[j] = this[j], this[i]
}

func NewSemver(version string) (*Semver, error) {
  semver := new(Semver)

  semver.Raw = version

  matches := RE_FULL.FindStringSubmatch(version)
  if matches == nil {
    return nil, fmt.Errorf("Invalid Version: %s", version)
  }

  semver.Major, _ = strconv.Atoi(matches[1])
  semver.Minor, _ = strconv.Atoi(matches[2])
  semver.Patch, _ = strconv.Atoi(matches[3])
  semver.PreRelease = NewPreRelease(matches[4])
  semver.Build = NewBuild(matches[5])
  semver.Version = semver.Format()

  Debug("created: %s\n", semver.Version)

  return semver, nil
}

func (semver *Semver) Format() string {
  str := fmt.Sprintf("%d.%d.%d", semver.Major, semver.Minor, semver.Patch)

  if !semver.PreRelease.Empty {
    str = fmt.Sprintf("%s-%s", str, semver.PreRelease.String())
  }

  if !semver.Build.Empty {
    str = fmt.Sprintf("%s-%s", str, semver.Build.String())
  }

  return str
}

func (semver *Semver) Inspect() string {
  return fmt.Sprintf("%v", semver)
}

func (semver *Semver) String() string {
  return semver.Version
}

func (a *Semver) Compare(b *Semver) int {
  Debug("comparing %s to %s\n", a, b)

  if a.Major > b.Major {
    return 1
  } else if a.Major < b.Major {
    return -1
  }

  if a.Minor > b.Minor {
    return 1
  } else if a.Minor < b.Minor {
    return -1
  }

  if a.Patch > b.Patch {
    return 1
  } else if a.Patch < b.Patch {
    return -1
  }

  return a.PreRelease.Compare(b.PreRelease)
}

func (semver *Semver) Increment(release string) error {
  switch release {
  case "major":
    semver.Major += 1
    semver.Minor = 0
    semver.Patch = 0
  case "minor":
    semver.Minor += 1
    semver.Patch = 0
  case "patch":
    semver.Patch += 1
    semver.PreRelease = NewPreRelease("")
    break
  case "prerelease":
    if semver.PreRelease.Empty {
      semver.PreRelease = NewPreRelease("0")
    } else {
      i := len(semver.PreRelease.Pieces) - 1
      for ; i >= 0; i -= 1 {
        if semver.PreRelease.Pieces[i].Numeric {
          semver.PreRelease.Pieces[i].SetInt(semver.PreRelease.Pieces[i].IntVal + 1)
          i = -2
        }
      }
      if i == -1 {
        semver.PreRelease.Pieces = append(semver.PreRelease.Pieces, NewIdentifier("0"))
      }
    }
  default:
    return errors.New("invalid increment argument: " + release)
  }

  semver.Version = semver.Format()
  return nil
}

func NewPreRelease(version string) *PreRelease {
  pieces := strings.Split(version, ".")

  pre := PreRelease{
    Raw:    version,
    Pieces: make([]*Indentifier, len(pieces)),
  }

  if version == "" {
    pre.Empty = true
  } else {
    for x := 0; x < len(pieces); x++ {
      pre.Pieces[x] = NewIdentifier(pieces[x])
    }
  }

  return &pre
}

func (prerelease *PreRelease) String() string {
  if prerelease.Empty {
    return ""
  }

  str := prerelease.Pieces[0].Val
  for x := 1; x < len(prerelease.Pieces); x++ {
    str = fmt.Sprintf("%s.%s", str, prerelease.Pieces[x].Val)
  }
  return str
}

func (a *PreRelease) Compare(b *PreRelease) int {
  // NOT having a PreRelease is > having one
  if !a.Empty && b.Empty {
    return -1
  } else if a.Empty && !b.Empty {
    return 1
  } else if a.Empty && b.Empty {
    return 0
  }

  aLen := len(a.Pieces)
  bLen := len(b.Pieces)

  i := 0
  for {
    if aLen == i && bLen == i {
      return 0
    } else if bLen == i {
      return 1
    } else if aLen == i {
      return -1
    } else {
      compare := a.Pieces[i].Compare(b.Pieces[i])
      if compare == 0 {
        i = i + 1
        continue
      } else {
        return compare
      }
    }
  }
}

func NewBuild(version string) *Build {
  pieces := strings.Split(version, ".")

  build := Build{
    Raw:    version,
    Pieces: make([]*Indentifier, len(pieces)),
  }

  if version == "" {
    build.Empty = true
  } else {
    for x := 0; x < len(pieces); x++ {
      build.Pieces[x] = NewIdentifier(pieces[x])
    }
  }

  return &build
}

func (build *Build) String() string {
  if build.Empty {
    return ""
  }

  str := build.Pieces[0].Val
  for x := 1; x < len(build.Pieces); x++ {
    str = fmt.Sprintf("%s.%s", str, build.Pieces[x].Val)
  }
  return str
}

func (build *Build) Compare(other *Build) int {
  return 0
}

func NewIdentifier(id string) *Indentifier {
  identifier := new(Indentifier)
  identifier.Set(id)
  return identifier
}

func (id *Indentifier) String() string {
  return id.Val
}

func (id *Indentifier) Set(version string) {
  id.Val = version

  intVal, err := strconv.Atoi(version)
  id.Numeric = err == nil

  if id.Numeric {
    id.IntVal = intVal
  }
}

func (id *Indentifier) SetInt(version int) {
  id.Val = strconv.Itoa(version)
  id.Numeric = true
  id.IntVal = version
}

func (a *Indentifier) Compare(b *Indentifier) int {
  if a.Numeric && !b.Numeric {
    return -1
  } else if b.Numeric && !a.Numeric {
    return 1
  } else if a.Numeric && b.Numeric {
    if a.IntVal < b.IntVal {
      return -1
    } else if a.IntVal > b.IntVal {
      return 1
    } else {
      return 0
    }
  } else {
    if a.Val < b.Val {
      return -1
    } else if a.Val > b.Val {
      return 1
    } else {
      return 0
    }
  }
}
