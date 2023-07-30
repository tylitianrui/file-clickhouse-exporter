package xfile

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	txta     = "hello\n"
	txtb     = "world\n"
	txtc     = "Delicate, lustrous, and soft to the touch. The ethereal fabric that we now call silk has threaded its way throughout China's history.One cannot be certain of its origin, but the humble ancient Chinese people credited their own wisdom to Leizu, wife of the Yellow Emperor (the legendary ancestor of all Chinese people), as the inventor of sericulture.\n"
	txtd     = "Like almost all genres of art on the vast land of China, the style and texture of silk are also multifarious. Hangluo satin from Hangzhou, Zhejiang province, is renowned for its airy and sheer texture, while Yunjin brocade from Nanjing, Jiangsu province, a luxurious fabric often used for royal garments, represents China's silk weaving technique at its prime.Yunjin brocade is best made by hand on giant looms, in a complex procedure that comprises more than a hundred steps. Even the most skilled artisans can only weave a few centimeters a day. Time, patience and deftness all play imperative roles to its heavenly beauty, or as its name suggests, its cloud-like splendor.In the Western Han Dynasty (206 BC-AD 24), with diplomat and explorer Zhang Qian opening up the routes to the western regions, silk graced countries in Central Asia, later extending its reach to other parts of Eurasia and beyond. Fittingly, its name marked China's major international trade routes, the ancient Silk Road and Maritime Silk Road.In the hands of Chinese artists, the thinnest threads can weave pictures of immense possibilities, and the softest material can traverse thousands of years. As one of the many marvels of ancient China, silk is not merely a type of textile. It is a cultural icon, and an embodiment of elegance and grace.\n"
	caselist = []string{txta, txtb, txtc, txtd, txtc, txtd, txtc, txtd, txtc, txtd, txtc, txtd}
)

func TestFileReader_ReadLine(t *testing.T) {
	a := assert.New(t)
	filename := "../../test/test_read_line"
	reader, err := NewFileReader(filename)
	a.NoError(err)
	var idx = 0
	for {
		b, err := reader.ReadLine()
		if err != nil {
			a.EqualError(err, io.EOF.Error())
			break
		}
		a.NoError(err)
		a.Equal(caselist[idx], string(b))
		idx++

	}

}
