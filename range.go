package semver

import (
  "fmt"
  "regexp"
  "strconv"
  "strings"
)

var RE_RANGE = regexp.MustCompile(`\s*\|\|\s*`)

type Range struct {
  Raw string
}

func NewRange(raw string) *Range {
  rng := &Range{
    Raw: raw,
  }

  rng.ParseRange(raw)

  return rng
}

func (this *Range) Format() {

}

// function Range(range, loose) {
//   if ((range instanceof Range) && range.loose === loose)
//     return range;

//   if (!(this instanceof Range))
//     return new Range(range, loose);

//   this.loose = loose;

//   // First, split based on boolean or ||
//   this.raw = range;
//   this.set = range.split(/\s*\|\|\s*/).map(function(range) {
//     return this.parseRange(range.trim());
//   }, this).filter(function(c) {
//     // throw out any that are not relevant for whatever reason
//     return c.length;
//   });

//   if (!this.set.length) {
//     throw new TypeError('Invalid SemVer Range: ' + range);
//   }

//   this.format();
// }

// Range.prototype.inspect = function() {
//   return '<SemVer Range "' + this.range + '">';
// };

// Range.prototype.format = function() {
//   this.range = this.set.map(function(comps) {
//     return comps.join(' ').trim();
//   }).join('||').trim();
//   return this.range;
// };

// Range.prototype.toString = function() {
//   return this.range;
// };

func (this *Range) ParseRange(rng string) {
  rng = strings.Trim(rng, " ")
  Debug("range %s\n", rng)
  matches := RE_HYPHENRANGE.FindAllStringSubmatch(rng, -1)
  Debug("matches: %v\n", len(matches[0][1:]))
  rng = hyphenReplace(matches[0])
  Debug("hyphen replace: %s\n", rng)
}

// Range.prototype.parseRange = function(range) {
//   var loose = this.loose;
//   range = range.trim();
//   debug('range', range, loose);
//   // `1.2.3 - 1.2.4` => `>=1.2.3 <=1.2.4`
//   var hr = loose ? re[HYPHENRANGELOOSE] : re[HYPHENRANGE];
//   range = range.replace(hr, hyphenReplace);
//   debug('hyphen replace', range);
//   // `> 1.2.3 < 1.2.5` => `>1.2.3 <1.2.5`
//   range = range.replace(re[COMPARATORTRIM], comparatorTrimReplace);
//   debug('comparator trim', range, re[COMPARATORTRIM]);

//   // `~ 1.2.3` => `~1.2.3`
//   range = range.replace(re[TILDETRIM], tildeTrimReplace);

//   // `^ 1.2.3` => `^1.2.3`
//   range = range.replace(re[CARETTRIM], caretTrimReplace);

//   // normalize spaces
//   range = range.split(/\s+/).join(' ');

//   // At this point, the range is completely trimmed and
//   // ready to be split into comparators.

//   var compRe = loose ? re[COMPARATORLOOSE] : re[COMPARATOR];
//   var set = range.split(' ').map(function(comp) {
//     return parseComparator(comp, loose);
//   }).join(' ').split(/\s+/);
//   if (this.loose) {
//     // in loose mode, throw out any that are not valid comparators
//     set = set.filter(function(comp) {
//       return !!comp.match(compRe);
//     });
//   }
//   set = set.map(function(comp) {
//     return new Comparator(comp, loose);
//   });

//   return set;
// };

func isX(id string) bool {
  return id == "" || strings.ToLower(id) == "x" || id == "*"
}

func hyphenReplace(m []string) string {
  from, fM, fm, fp, _, _, to, tM, tm, tp, tpr, _ := m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10], m[11]

  if isX(fM) {
    from = ""
  } else if isX(fm) {
    from = ">=" + fM + ".0.0-0"
  } else if isX(fp) {
    from = ">=" + fM + "." + fm + ".0-0"
  } else {
    from = ">=" + from
  }

  if isX(tM) {
    to = ""
  } else if isX(tm) {
    inc, _ := strconv.Atoi(tM)
    to = fmt.Sprintf("<%d.0.0-0", inc+1)
  } else if isX(tp) {
    inc, _ := strconv.Atoi(tm)
    to = fmt.Sprintf("<%s.%d.0-0", tM, inc+1)
  } else if tpr != "" {
    to = fmt.Sprintf("<=%s.%s.%s-%s", tM, tm, tp, tpr)
  } else {
    to = "<=" + to
  }

  return strings.Trim(from+" "+to, " ")
}
