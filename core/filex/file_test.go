package filex

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/fs"
)

const (
	longLine      = `Quid securi etiam tamquam eu fugiat nulla pariatur. Nec dubitamus multa iter quae et nos invenerat. Non equidem invideo, miror magis posuere velit aliquet. Integer legentibus erat a ante historiarum dapibus. Prima luce, cum quibus mons aliud consensu ab eo.Quid securi etiam tamquam eu fugiat nulla pariatur. Nec dubitamus multa iter quae et nos invenerat. Non equidem invideo, miror magis posuere velit aliquet. Integer legentibus erat a ante historiarum dapibus. Prima luce, cum quibus mons aliud consensu ab eo.Quid securi etiam tamquam eu fugiat nulla pariatur. Nec dubitamus multa iter quae et nos invenerat. Non equidem invideo, miror magis posuere velit aliquet. Integer legentibus erat a ante historiarum dapibus. Prima luce, cum quibus mons aliud consensu ab eo.Quid securi etiam tamquam eu fugiat nulla pariatur. Nec dubitamus multa iter quae et nos invenerat. Non equidem invideo, miror magis posuere velit aliquet. Integer legentibus erat a ante historiarum dapibus. Prima luce, cum quibus mons aliud consensu ab eo.Quid securi etiam tamquam eu fugiat nulla pariatur. Nec dubitamus multa iter quae et nos invenerat. Non equidem invideo, miror magis posuere velit aliquet. Integer legentibus erat a ante historiarum dapibus. Prima luce, cum quibus mons aliud consensu ab eo.`
	longFirstLine = longLine + "\n" + text
	text          = `first line
Cum sociis natoque penatibus et magnis dis parturient. Phasellus laoreet lorem vel dolor tempus vehicula. Vivamus sagittis lacus vel augue laoreet rutrum faucibus. Integer legentibus erat a ante historiarum dapibus.
Quisque ut dolor gravida, placerat libero vel, euismod. Quam temere in vitiis, legem sancimus haerentia. Qui ipsorum lingua Celtae, nostra Galli appellantur. Quis aute iure reprehenderit in voluptate velit esse. Fabio vel iudice vincam, sunt in culpa qui officia. Cras mattis iudicium purus sit amet fermentum.
Quo usque tandem abutere, Catilina, patientia nostra? Gallia est omnis divisa in partes tres, quarum. Quam diu etiam furor iste tuus nos eludet? Quid securi etiam tamquam eu fugiat nulla pariatur. Curabitur blandit tempus ardua ridiculous sed magna.
Magna pars studiorum, prodita quaerimus. Cum ceteris in veneratione tui montes, nascetur mus. Morbi odio eros, volutpat ut pharetra vitae, lobortis sed nibh. Plura mihi bona sunt, inclinet, amari petere vellent. Idque Caesaris facere voluntate liceret: sese habere. Tu quoque, Brute, fili mi, nihil timor populi, nihil!
Tityre, tu patulae recubans sub tegmine fagi dolor. Inmensae subtilitatis, obscuris et malesuada fames. Quae vero auctorem tractata ab fiducia dicuntur.
Cum sociis natoque penatibus et magnis dis parturient. Phasellus laoreet lorem vel dolor tempus vehicula. Vivamus sagittis lacus vel augue laoreet rutrum faucibus. Integer legentibus erat a ante historiarum dapibus.
Quisque ut dolor gravida, placerat libero vel, euismod. Quam temere in vitiis, legem sancimus haerentia. Qui ipsorum lingua Celtae, nostra Galli appellantur. Quis aute iure reprehenderit in voluptate velit esse. Fabio vel iudice vincam, sunt in culpa qui officia. Cras mattis iudicium purus sit amet fermentum.
Quo usque tandem abutere, Catilina, patientia nostra? Gallia est omnis divisa in partes tres, quarum. Quam diu etiam furor iste tuus nos eludet? Quid securi etiam tamquam eu fugiat nulla pariatur. Curabitur blandit tempus ardua ridiculous sed magna.
Magna pars studiorum, prodita quaerimus. Cum ceteris in veneratione tui montes, nascetur mus. Morbi odio eros, volutpat ut pharetra vitae, lobortis sed nibh. Plura mihi bona sunt, inclinet, amari petere vellent. Idque Caesaris facere voluntate liceret: sese habere. Tu quoque, Brute, fili mi, nihil timor populi, nihil!
Tityre, tu patulae recubans sub tegmine fagi dolor. Inmensae subtilitatis, obscuris et malesuada fames. Quae vero auctorem tractata ab fiducia dicuntur.
Cum sociis natoque penatibus et magnis dis parturient. Phasellus laoreet lorem vel dolor tempus vehicula. Vivamus sagittis lacus vel augue laoreet rutrum faucibus. Integer legentibus erat a ante historiarum dapibus.
Quisque ut dolor gravida, placerat libero vel, euismod. Quam temere in vitiis, legem sancimus haerentia. Qui ipsorum lingua Celtae, nostra Galli appellantur. Quis aute iure reprehenderit in voluptate velit esse. Fabio vel iudice vincam, sunt in culpa qui officia. Cras mattis iudicium purus sit amet fermentum.
Quo usque tandem abutere, Catilina, patientia nostra? Gallia est omnis divisa in partes tres, quarum. Quam diu etiam furor iste tuus nos eludet? Quid securi etiam tamquam eu fugiat nulla pariatur. Curabitur blandit tempus ardua ridiculous sed magna.
Magna pars studiorum, prodita quaerimus. Cum ceteris in veneratione tui montes, nascetur mus. Morbi odio eros, volutpat ut pharetra vitae, lobortis sed nibh. Plura mihi bona sunt, inclinet, amari petere vellent. Idque Caesaris facere voluntate liceret: sese habere. Tu quoque, Brute, fili mi, nihil timor populi, nihil!
Tityre, tu patulae recubans sub tegmine fagi dolor. Inmensae subtilitatis, obscuris et malesuada fames. Quae vero auctorem tractata ab fiducia dicuntur.
` + longLine
	textWithLastNewline = `first line
Cum sociis natoque penatibus et magnis dis parturient. Phasellus laoreet lorem vel dolor tempus vehicula. Vivamus sagittis lacus vel augue laoreet rutrum faucibus. Integer legentibus erat a ante historiarum dapibus.
Quisque ut dolor gravida, placerat libero vel, euismod. Quam temere in vitiis, legem sancimus haerentia. Qui ipsorum lingua Celtae, nostra Galli appellantur. Quis aute iure reprehenderit in voluptate velit esse. Fabio vel iudice vincam, sunt in culpa qui officia. Cras mattis iudicium purus sit amet fermentum.
Quo usque tandem abutere, Catilina, patientia nostra? Gallia est omnis divisa in partes tres, quarum. Quam diu etiam furor iste tuus nos eludet? Quid securi etiam tamquam eu fugiat nulla pariatur. Curabitur blandit tempus ardua ridiculous sed magna.
Magna pars studiorum, prodita quaerimus. Cum ceteris in veneratione tui montes, nascetur mus. Morbi odio eros, volutpat ut pharetra vitae, lobortis sed nibh. Plura mihi bona sunt, inclinet, amari petere vellent. Idque Caesaris facere voluntate liceret: sese habere. Tu quoque, Brute, fili mi, nihil timor populi, nihil!
Tityre, tu patulae recubans sub tegmine fagi dolor. Inmensae subtilitatis, obscuris et malesuada fames. Quae vero auctorem tractata ab fiducia dicuntur.
Cum sociis natoque penatibus et magnis dis parturient. Phasellus laoreet lorem vel dolor tempus vehicula. Vivamus sagittis lacus vel augue laoreet rutrum faucibus. Integer legentibus erat a ante historiarum dapibus.
Quisque ut dolor gravida, placerat libero vel, euismod. Quam temere in vitiis, legem sancimus haerentia. Qui ipsorum lingua Celtae, nostra Galli appellantur. Quis aute iure reprehenderit in voluptate velit esse. Fabio vel iudice vincam, sunt in culpa qui officia. Cras mattis iudicium purus sit amet fermentum.
Quo usque tandem abutere, Catilina, patientia nostra? Gallia est omnis divisa in partes tres, quarum. Quam diu etiam furor iste tuus nos eludet? Quid securi etiam tamquam eu fugiat nulla pariatur. Curabitur blandit tempus ardua ridiculous sed magna.
Magna pars studiorum, prodita quaerimus. Cum ceteris in veneratione tui montes, nascetur mus. Morbi odio eros, volutpat ut pharetra vitae, lobortis sed nibh. Plura mihi bona sunt, inclinet, amari petere vellent. Idque Caesaris facere voluntate liceret: sese habere. Tu quoque, Brute, fili mi, nihil timor populi, nihil!
Tityre, tu patulae recubans sub tegmine fagi dolor. Inmensae subtilitatis, obscuris et malesuada fames. Quae vero auctorem tractata ab fiducia dicuntur.
Cum sociis natoque penatibus et magnis dis parturient. Phasellus laoreet lorem vel dolor tempus vehicula. Vivamus sagittis lacus vel augue laoreet rutrum faucibus. Integer legentibus erat a ante historiarum dapibus.
Quisque ut dolor gravida, placerat libero vel, euismod. Quam temere in vitiis, legem sancimus haerentia. Qui ipsorum lingua Celtae, nostra Galli appellantur. Quis aute iure reprehenderit in voluptate velit esse. Fabio vel iudice vincam, sunt in culpa qui officia. Cras mattis iudicium purus sit amet fermentum.
Quo usque tandem abutere, Catilina, patientia nostra? Gallia est omnis divisa in partes tres, quarum. Quam diu etiam furor iste tuus nos eludet? Quid securi etiam tamquam eu fugiat nulla pariatur. Curabitur blandit tempus ardua ridiculous sed magna.
Magna pars studiorum, prodita quaerimus. Cum ceteris in veneratione tui montes, nascetur mus. Morbi odio eros, volutpat ut pharetra vitae, lobortis sed nibh. Plura mihi bona sunt, inclinet, amari petere vellent. Idque Caesaris facere voluntate liceret: sese habere. Tu quoque, Brute, fili mi, nihil timor populi, nihil!
Tityre, tu patulae recubans sub tegmine fagi dolor. Inmensae subtilitatis, obscuris et malesuada fames. Quae vero auctorem tractata ab fiducia dicuntur.
` + longLine + "\n"
	shortText = `first line
second line
last line`
	shortTextWithLastNewline = `first line
second line
last line
`
)

func TestFirstLine(t *testing.T) {
	filename, err := fs.TempFilenameWithText(longFirstLine)
	assert.Nil(t, err)
	defer os.Remove(filename)

	val, err := FirstLine(filename)
	assert.Nil(t, err)
	assert.Equal(t, longLine, val)
}

func TestFirstLineShort(t *testing.T) {
	filename, err := fs.TempFilenameWithText(shortText)
	assert.Nil(t, err)
	defer os.Remove(filename)

	val, err := FirstLine(filename)
	assert.Nil(t, err)
	assert.Equal(t, "first line", val)
}

func TestLastLine(t *testing.T) {
	filename, err := fs.TempFilenameWithText(text)
	assert.Nil(t, err)
	defer os.Remove(filename)

	val, err := LastLine(filename)
	assert.Nil(t, err)
	assert.Equal(t, longLine, val)
}

func TestLastLineWithLastNewline(t *testing.T) {
	filename, err := fs.TempFilenameWithText(textWithLastNewline)
	assert.Nil(t, err)
	defer os.Remove(filename)

	val, err := LastLine(filename)
	assert.Nil(t, err)
	assert.Equal(t, longLine, val)
}

func TestLastLineShort(t *testing.T) {
	filename, err := fs.TempFilenameWithText(shortText)
	assert.Nil(t, err)
	defer os.Remove(filename)

	val, err := LastLine(filename)
	assert.Nil(t, err)
	assert.Equal(t, "last line", val)
}

func TestLastLineWithLastNewlineShort(t *testing.T) {
	filename, err := fs.TempFilenameWithText(shortTextWithLastNewline)
	assert.Nil(t, err)
	defer os.Remove(filename)

	val, err := LastLine(filename)
	assert.Nil(t, err)
	assert.Equal(t, "last line", val)
}
