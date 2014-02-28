package semver

import (
  "fmt"
  "regexp"
  "sort"
)

const NUMERICIDENTIFIER = "0|[1-9]\\d*"
const NONNUMERICIDENTIFIER = "\\d*[a-zA-Z-][a-zA-Z0-9-]*"
const MAINVERSION = "(" + NUMERICIDENTIFIER + ")\\." +
  "(" + NUMERICIDENTIFIER + ")\\." +
  "(" + NUMERICIDENTIFIER + ")"
const PRERELEASEIDENTIFIER = "(?:" + NUMERICIDENTIFIER +
  "|" + NONNUMERICIDENTIFIER + ")"
const PRERELEASE = "(?:-(" + PRERELEASEIDENTIFIER +
  "(?:\\." + PRERELEASEIDENTIFIER + ")*))"
const BUILDIDENTIFIER = "[0-9A-Za-z-]+"
const BUILD = "(?:\\+(" + BUILDIDENTIFIER + "(?:\\." + BUILDIDENTIFIER + ")*))"
const FULLPLAIN = "v?" + MAINVERSION + PRERELEASE + "?" + BUILD + "?"
const FULL = "^" + FULLPLAIN + "$"
const GTLT = "((?:<|>)?=?)"
const XRANGEIDENTIFIER = NUMERICIDENTIFIER + "|x|X|\\*"
const XRANGEPLAIN = "[v=\\s]*(" + XRANGEIDENTIFIER + ")" +
  "(?:\\.(" + XRANGEIDENTIFIER + ")" +
  "(?:\\.(" + XRANGEIDENTIFIER + ")" +
  "(?:(" + PRERELEASE + ")" +
  ")?)?)?"
const XRANGE = "^" + GTLT + "\\s*" + XRANGEPLAIN + "$"
const LONETILDE = "(?:~>?)"
const TILDETRIM = "(\\s*)" + LONETILDE + "\\s+"
const TILDE = "^" + LONETILDE + XRANGEPLAIN + "$"
const LONECARET = "(?:\\^)"
const CARETTRIM = "(\\s*)" + LONECARET + "\\s+"
const CARET = "^" + LONECARET + XRANGEPLAIN + "$"
const COMPARATOR = "^" + GTLT + "\\s*(" + FULLPLAIN + ")$|^$"
const COMPARATORTRIM = "(\\s*)" + GTLT +
  "\\s*(" + FULLPLAIN + "|" + XRANGEPLAIN + ")"
const HYPHENRANGE = "^\\s*(" + XRANGEPLAIN + ")" +
  "\\s+-\\s+" +
  "(" + XRANGEPLAIN + ")" +
  "\\s*$"
const STAR = "(<|>)?=?\\s*\\*"

var RE_NUMERICIDENTIFIER = regexp.MustCompile(NUMERICIDENTIFIER)
var RE_NONNUMERICIDENTIFIER = regexp.MustCompile(NONNUMERICIDENTIFIER)
var RE_MAINVERSION = regexp.MustCompile(MAINVERSION)
var RE_PRERELEASEIDENTIFIER = regexp.MustCompile(PRERELEASEIDENTIFIER)
var RE_PRERELEASE = regexp.MustCompile(PRERELEASE)
var RE_BUILDIDENTIFIER = regexp.MustCompile(BUILDIDENTIFIER)
var RE_BUILD = regexp.MustCompile(BUILD)
var RE_FULLPLAIN = regexp.MustCompile(FULLPLAIN)
var RE_FULL = regexp.MustCompile(FULL)
var RE_GTLT = regexp.MustCompile(GTLT)
var RE_XRANGEIDENTIFIER = regexp.MustCompile(XRANGEIDENTIFIER)
var RE_XRANGEPLAIN = regexp.MustCompile(XRANGEPLAIN)
var RE_XRANGE = regexp.MustCompile(XRANGE)
var RE_LONETILDE = regexp.MustCompile(LONETILDE)
var RE_TILDETRIM = regexp.MustCompile(TILDETRIM)
var RE_TILDE = regexp.MustCompile(TILDE)
var RE_LONECARET = regexp.MustCompile(LONECARET)
var RE_CARETTRIM = regexp.MustCompile(CARETTRIM)
var RE_CARET = regexp.MustCompile(CARET)
var RE_COMPARATOR = regexp.MustCompile(COMPARATOR)
var RE_COMPARATORTRIM = regexp.MustCompile(COMPARATORTRIM)
var RE_HYPHENRANGE = regexp.MustCompile(HYPHENRANGE)
var RE_STAR = regexp.MustCompile(STAR)

var debug = false

func Debug(format string, a ...interface{}) {
  if debug {
    fmt.Printf(format, a...)
  }
}

func Sort(versions []Semver) []Semver {
  sort.Sort(Semvers(versions))
  return versions
}

func Compare(a string, b string) (int, error) {
  aVersion, err := NewSemver(a)
  if err != nil {
    return 0, err
  }

  bVersion, err := NewSemver(b)
  if err != nil {
    return 0, err
  }

  return aVersion.Compare(bVersion), nil
}

func Cmp(a string, op string, b string) (bool, error) {
  switch op {
  case "===":
    return a == b, nil
  case "!==":
    return a != b, nil
  case "", "=", "==":
    return EQ(a, b), nil
  case "!=":
    return NEQ(a, b), nil
  case ">":
    return GT(a, b), nil
  case ">=":
    return GTE(a, b), nil
  case "<":
    return LT(a, b), nil
  case "<=":
    return LTE(a, b), nil
  default:
    return false, fmt.Errorf("Invalid operator: %s", op)
  }
}

func GT(a string, b string) bool {
  compare, _ := Compare(a, b)
  return compare > 0
}

func LT(a string, b string) bool {
  compare, _ := Compare(a, b)
  return compare < 0
}

func EQ(a string, b string) bool {
  compare, _ := Compare(a, b)
  return compare == 0
}

func NEQ(a string, b string) bool {
  compare, _ := Compare(a, b)
  return compare != 0
}

func GTE(a string, b string) bool {
  compare, _ := Compare(a, b)
  return compare >= 0
}

func LTE(a string, b string) bool {
  compare, _ := Compare(a, b)
  return compare <= 0
}
