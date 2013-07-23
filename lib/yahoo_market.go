// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	`bytes`
	`io/ioutil`
	`net/http`
	`regexp`
	`strings`
)

type Market struct {
	Open      bool
	Dow       map[string]string
	Nasdaq    map[string]string
	Sp500     map[string]string
	Advances  map[string]string
	Declines  map[string]string
	Unchanged map[string]string
	Highs     map[string]string
	Lows      map[string]string
}

//-----------------------------------------------------------------------------
func (self *Market) Initialize() *Market {
	self.Open       = true
	self.Dow        = make(map[string]string)
	self.Nasdaq     = make(map[string]string)
	self.Sp500      = make(map[string]string)
	self.Advances   = make(map[string]string)
	self.Declines   = make(map[string]string)
	self.Unchanged  = make(map[string]string)
	self.Highs      = make(map[string]string)
	self.Lows       = make(map[string]string)

	return self
}

//-----------------------------------------------------------------------------
func (self *Market) Fetch() *Market {
	response, err := http.Get(`http://finance.yahoo.com/marketupdate/overview`)
	if err != nil {
		panic(err)
	}

	// Fetch response and get its body.
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	body = self.check_if_market_open(body)
	return self.extract(self.trim(body))
}

//-----------------------------------------------------------------------------
func (self *Market) Format() string {
	return new(Formatter).Format(self)
}

// private
//-----------------------------------------------------------------------------
func (self *Market) check_if_market_open(body []byte) []byte {
	start := bytes.Index(body, []byte(`id="yfs_market_time"`))
	finish := start + bytes.Index(body[start:], []byte(`</span>`))
	self.Open = !bytes.Contains(body[start:finish], []byte(`closed`))

	return body[finish:]
}

//-----------------------------------------------------------------------------
func (self *Market) trim(body []byte) []byte {
	start := bytes.Index(body, []byte(`<table id="yfimktsumm"`))
	finish := bytes.LastIndex(body, []byte(`<table id="yfimktsumm"`))
	snippet := bytes.Replace(body[start:finish], []byte{'\n'}, []byte{}, -1)
	snippet = bytes.Replace(snippet, []byte(`&amp;`), []byte{'&'}, -1)

	return snippet
}

//-----------------------------------------------------------------------------
func (self *Market) extract(snippet []byte) *Market {
	const any = `\s*<.+?>`
	const some = `<.+?`
	const space = `\s*`
	const color = `#([08c]{6});">\s*`
	const price = `([\d\.,]+)`
	const percent = `\(([\d\.,%]+)\)`

	regex := []string{
		`(Dow)`, any, price, some, color, price, some, percent, any,
		`(Nasdaq)`, any, price, some, color, price, some, percent, any,
		`(S&P 500)`, any, price, some, color, price, some, percent, any,
		`(Advances)`, any, price, space, percent, any, price, space, percent, any,
		`(Declines)`, any, price, space, percent, any, price, space, percent, any,
		`(Unchanged)`, any, price, space, percent, any, price, space, percent, any,
		`(New Hi's)`, any, price, any, price, any,
		`(New Lo's)`, any, price, any, price, any,
	}

	re := regexp.MustCompile(strings.Join(regex, ``))
	matches := re.FindAllStringSubmatch(string(snippet), -1)

	// if len(matches) > 0 {
	//         fmt.Printf("%d matches\n", len(matches[0]))
	//         for i, str := range matches[0][1:] {
	//                 fmt.Printf("%d) [%s]\n", i, str)
	//         }
	// } else {
	//         println(`No matches`)
	// }


	self.Dow[`name`] = matches[0][1]
	self.Dow[`latest`] = matches[0][2]
	self.Dow[`change`] = matches[0][4]
	switch matches[0][3] {
	case `008800`:
		self.Dow[`change`] = `+` + matches[0][4]
		self.Dow[`percent`] = `+` + matches[0][5]
	case `cc0000`:
		self.Dow[`change`] = `-` + matches[0][4]
		self.Dow[`percent`] = `-` + matches[0][5]
	default:
		self.Dow[`change`] = matches[0][4]
		self.Dow[`percent`] = matches[0][5]
	}

	self.Nasdaq[`name`] = matches[0][6]
	self.Nasdaq[`latest`] = matches[0][7]
	switch matches[0][8] {
	case `008800`:
		self.Nasdaq[`change`] = `+` + matches[0][9]
		self.Nasdaq[`percent`] = `+` + matches[0][10]
	case `cc0000`:
		self.Nasdaq[`change`] = `-` + matches[0][9]
		self.Nasdaq[`percent`] = `-` + matches[0][10]
	default:
		self.Nasdaq[`change`] = matches[0][9]
		self.Nasdaq[`percent`] = matches[0][10]
	}

	self.Sp500[`name`] = matches[0][11]
	self.Sp500[`latest`] = matches[0][12]
	switch matches[0][13] {
	case `008800`:
		self.Sp500[`change`] = `+` + matches[0][14]
		self.Sp500[`percent`] = `+` + matches[0][15]
	case `cc0000`:
		self.Sp500[`change`] = `-` + matches[0][14]
		self.Sp500[`percent`] = `-` + matches[0][15]
	default:
		self.Sp500[`change`] = matches[0][14]
		self.Sp500[`percent`] = matches[0][15]
	}

	self.Advances[`name`] = matches[0][16]
	self.Advances[`nyse`] = matches[0][17]
	self.Advances[`nysep`] = matches[0][18]
	self.Advances[`nasdaq`] = matches[0][19]
	self.Advances[`nasdaqp`] = matches[0][20]

	self.Declines[`name`] = matches[0][21]
	self.Declines[`nyse`] = matches[0][22]
	self.Declines[`nysep`] = matches[0][23]
	self.Declines[`nasdaq`] = matches[0][24]
	self.Declines[`nasdaqp`] = matches[0][25]

	self.Unchanged[`name`] = matches[0][26]
	self.Unchanged[`nyse`] = matches[0][27]
	self.Unchanged[`nysep`] = matches[0][28]
	self.Unchanged[`nasdaq`] = matches[0][29]
	self.Unchanged[`nasdaqp`] = matches[0][30]

	self.Highs[`name`] = matches[0][31]
	self.Highs[`nyse`] = matches[0][32]
	self.Highs[`nasdaq`] = matches[0][33]
	self.Lows[`name`] = matches[0][34]
	self.Lows[`nyse`] = matches[0][35]
	self.Lows[`nasdaq`] = matches[0][36]

	return self
}
