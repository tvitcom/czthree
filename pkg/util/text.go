package util

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	// "strconv"
	"github.com/tvitcom/czthree/bkp/swearfilter"
)

func SwearDetector(message string) ([]string, error) {
	var swears = []string{`аеб`,`аёб`,`банамат`,`бля`,`6ля`,`6ляд`,`6лят`,`зъeб`,`cock`,`cunt`,`e6aл`,`ebal`,
	`eblan`,`eбaт`,`eбyч`,`eбёт`,`eблан`,`fuck`,`xyё`,`xyй`,`xyя`,`xуе`,`xую`,`zaeb`,`zaebal`,`zaebali`,
	`zaebat`,`бзд`,`вафел`,`вафлё`,`перд`,`высра`,`ссать`,`ебен`,`гавн`,`гамн`,`гандо`,`гнид`,`говен`,
	`говешк`,`говн`,`гондон`,`долбо`,`дрисн`,`дрист`,`дроче`,`дрочи`,`дрочк`,`дрочу`,`е6ал`,`е6ут`,
	`вашумать`,`ёбaн`,`ебaть`,`ебyч`,`ебал`,`ебан`,`ёбан`,`ебарь`,`ебат`,`ёбат`,`ебашить`,`ебё`,`ебет`,
	`ебец`,`ебик`,`ебин`,`ебись`,`ебиче`,`ебки`,`ебла`,`ебливый`,`еблище`,`ебло`,`еблыс`,`ебля`,`ёбн`,
	`ебнут`,`ебня`,`ебошить`,`ебск`,`ебтв`,`ебун`,`ебут`,`ебуч`,`ебущ`,`ебырь`,`елда`,`ебет`,`жоп`,
	`задра`,`дрот`,`зае6`,`заё6`,`залуп`,`засер`,`засир`,`засру`,`заябестая`,`ибонех`,`издев`,`ипать`,
	`ипацц`,`конча`,`курва`,`курвят`,` лох`,`лошар`,`лошок`,`лярва`,`малафь`,`манда`,`мандей`,`мандень`,
	`мандеть`,`мандища`,`мандой`,`манду`,`мандюк`,`минет`,`млять`,`щелка$`,`щёлка`,`мраз`,`мудak`,
	`мудaк`,`мудаг`,`мудак`,`мудач`,`мудазв`,`мудозв`,`муде`,`муди`,`мудн`,`мудо`,`наговнять`,`нахер`,
	`нехира`,`хера`,`нииб`,`ниипац`,`ниипет`,`никуя`,`обдрист`,`обоср`,`обосц`,`обсир`,`оеба`,`пезд`,
	`отмудох`,`отпороть`,`очко`,`очку`,`падон`,`педерас`,`педик`,`педри`,`перд`,`пёрд`,`ёбок`,`пернуть`,
	`пёрнуть`,`пи3д`,`пи3d`,`pizd`,`pi3d`,`пи3де`,`пи3ду`,`пиzд`,`пидар`,`пидор`,`пидрас`,`пизд`,`письк`,
	`писю`,`подонк`,`падонок`,`оебе`,`оёбы`,`скуда`,`срать`,`скух`,`потаскушка`,`похер`,`хрен`,`похрену`,
	`похую`,`придурок`,`приебаться`,`проблядь`,`своло`,`серун`,`серька`,`роеб`,`соси`,`срак`,`сран`,
	`срун`,`ссака`,`ссышь`,`стерва`,`сука`,`суки`,`суходрочка`,`сучар`,`сучий`,`сучка`,`сучко`,` сучоно`,
	`сучь`,`сцание`,`сцать`,`сцук`,`сцуль`,`сцыха`,`сцышь`,`сыкун`,`трах`,`ублюд`,`уеб`,`уёб`,`ёбки`,`ебок`,
	`урюк`,`ушлеп`,`х_у_`,`хyё`,`хyй`,`хамло`,`херо`,`херн`,`ыеб`,`хитрожопый`,`хуе`,`хуё`,
	`хуи`,`хули`,`хуля`,`хуя`,`целка`,`чмо`,`чмыр`,`шалав`,`шлюх`,`шлюшк`,`ъеб`,`ьеб`,`ябывает`}
	filter := swearfilter.NewSwearFilter(true, swears...)
	return filter.Check(message)
}

func IsBadWords(str string, bads string) bool {
	// bads := "dsa,dri,def,da,dass,sdd"
	for _, v := range strings.Split(bads, ",") {
		re := regexp.MustCompile(`.?`+ v + `.?`)
		res := re.Match([]byte(str))
		if res {
			return true
		}
	}
	return false
}

//TODO use that: regexp.QuoteMeta(`Escaping symbols like: .+*?()|[]{}^$`)
func Addcslashes(s string) string {
	var	slashedCharList string = "_*~'`<>#|\\"
	var result []rune
	for _, ch := range []rune(s) {
		for _, v := range []rune(slashedCharList) {
			if ch == v {
				result = append(result, '\\')
			}
		}
		result = append(result, ch)
	}
	return string(result)
}

func PhoneNormalisation(n string) string {
	//remove `()-+:\b` in phone number string
	r := strings.NewReplacer("-", "", "(", "", ")", "", "+", "", " ", "", ":", "")
	return r.Replace(n)
}

func ExtractTitle(str string, limitChars, limitWords int) string {
  // Validation and Triming
  prefix := "Продам"
  str = strings.TrimSpace(strings.ToValidUTF8(str,""))
  var fields []string
  fields = strings.Fields(str)
  if strings.EqualFold(fields[0], prefix) {
      fields = fields[1:]
  }
  if foldedLen := len(fields); foldedLen > limitWords {
    fields[0] = strings.Title(fields[0])
    fields = fields[:limitWords]
  }
  return Substr(strings.Join(fields, " "), 0, limitChars)
}

func Substr(input string, start int, length int) string {
    asRunes := []rune(input)
    if start >= len(asRunes) {
        return ""
    }
    if start+length > len(asRunes) {
        length = len(asRunes) - start
    }
    return string(asRunes[start : start+length])
}

func ExtractDigitsString(str string) string {
	var digits []rune
	for _, c := range str {
	    if unicode.IsDigit(c) {
	        digits = append(digits, c)
	    }
	}
	return string(digits)
}

func Stringer(raw interface{}) string {
    return fmt.Sprint(raw)
}
