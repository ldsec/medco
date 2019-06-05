package identifiers_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlleleMaping(t *testing.T) {
	res, err := AlleleMaping("A")
	assert.Equal(t, res, int64(0))
	assert.Nil(t, err)

	res, err = AlleleMaping("T")
	assert.Equal(t, res, int64(1))
	assert.Nil(t, err)

	res, err = AlleleMaping("G")
	assert.Equal(t, res, int64(2))
	assert.Nil(t, err)

	res, err = AlleleMaping("C")
	assert.Equal(t, res, int64(3))
	assert.Nil(t, err)

	res, err = AlleleMaping("")
	assert.Equal(t, res, int64(-1))
	assert.NotNil(t, err)

	res, err = AlleleMaping("test")
	assert.Equal(t, res, int64(-1))
	assert.NotNil(t, err)
}

func TestGetMask(t *testing.T) {
	assert.Equal(t, GetMask(1), int64(1))
	assert.Equal(t, GetMask(4), int64(15))
	assert.Equal(t, GetMask(10), int64(1023))
}

func TestPushBitsFromRight(t *testing.T) {
	assert.Equal(t, PushBitsFromRight(int64(0), 2, int64(1)), int64(1))
	assert.Equal(t, PushBitsFromRight(int64(0), 2, int64(7)), int64(3))
	assert.Equal(t, PushBitsFromRight(int64(0), 3, int64(7)), int64(7))
	assert.Equal(t, PushBitsFromRight(int64(0), 4, int64(7)), int64(7))
}

func TestEncodeAlleles(t *testing.T) {
	assert.Equal(t, EncodeAlleles("A"), int64(0))
	assert.Equal(t, EncodeAlleles("T"), int64(1024))
	assert.Equal(t, EncodeAlleles("G"), int64(2048))
	assert.Equal(t, EncodeAlleles("C"), int64(3072))

	assert.Equal(t, EncodeAlleles("AA"), int64(0))
	assert.Equal(t, EncodeAlleles("ATCG"), int64(480))
	assert.Equal(t, EncodeAlleles("GGTTCA"), int64(2652))
	assert.Equal(t, EncodeAlleles("TGACTA"), int64(1588))

	assert.Equal(t, EncodeAlleles("TGACTAT"), int64(0)) //strange!!
}

func TestGetVariantID(t *testing.T) {
	res, err := GetVariantID("1", int64(6), "AC", "ATTT")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(-8935141653966995120))

	res, err = GetVariantID("10", int64(2300), "C", "T")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(-6341065805496577024))

	res, err = GetVariantID("1", int64(999999), "TAAAC", "G")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(-8934067919247763456))

}
