package identifiers_test

import (
	"github.com/ldsec/medco-loader/loader/identifiers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlleleMaping(t *testing.T) {
	res, err := identifiers.AlleleMaping("A")
	assert.Equal(t, res, int64(0))
	assert.Nil(t, err)

	res, err = identifiers.AlleleMaping("T")
	assert.Equal(t, res, int64(1))
	assert.Nil(t, err)

	res, err = identifiers.AlleleMaping("G")
	assert.Equal(t, res, int64(2))
	assert.Nil(t, err)

	res, err = identifiers.AlleleMaping("C")
	assert.Equal(t, res, int64(3))
	assert.Nil(t, err)

	res, err = identifiers.AlleleMaping("")
	assert.Equal(t, res, int64(-1))
	assert.NotNil(t, err)

	res, err = identifiers.AlleleMaping("test")
	assert.Equal(t, res, int64(-1))
	assert.NotNil(t, err)
}

func TestGetMask(t *testing.T) {
	assert.Equal(t, identifiers.GetMask(1), int64(1))
	assert.Equal(t, identifiers.GetMask(4), int64(15))
	assert.Equal(t, identifiers.GetMask(10), int64(1023))
}

func TestPushBitsFromRight(t *testing.T) {
	assert.Equal(t, identifiers.PushBitsFromRight(int64(0), 2, int64(1)), int64(1))
	assert.Equal(t, identifiers.PushBitsFromRight(int64(0), 2, int64(7)), int64(3))
	assert.Equal(t, identifiers.PushBitsFromRight(int64(0), 3, int64(7)), int64(7))
	assert.Equal(t, identifiers.PushBitsFromRight(int64(0), 4, int64(7)), int64(7))
}

func TestEncodeAlleles(t *testing.T) {
	assert.Equal(t, identifiers.EncodeAlleles("A"), int64(0))
	assert.Equal(t, identifiers.EncodeAlleles("T"), int64(1024))
	assert.Equal(t, identifiers.EncodeAlleles("G"), int64(2048))
	assert.Equal(t, identifiers.EncodeAlleles("C"), int64(3072))

	assert.Equal(t, identifiers.EncodeAlleles("AA"), int64(0))
	assert.Equal(t, identifiers.EncodeAlleles("ATCG"), int64(480))
	assert.Equal(t, identifiers.EncodeAlleles("GGTTCA"), int64(2652))
	assert.Equal(t, identifiers.EncodeAlleles("TGACTA"), int64(1588))

	assert.Equal(t, identifiers.EncodeAlleles("TGACTAT"), int64(0)) //strange!!
}

func TestGetVariantID(t *testing.T) {
	res, err := identifiers.GetVariantID("1", int64(6), "AC", "ATTT")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(-8935141653966995120))

	res, err = identifiers.GetVariantID("10", int64(2300), "C", "T")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(-6341065805496577024))

	res, err = identifiers.GetVariantID("1", int64(999999), "TAAAC", "G")
	assert.Nil(t, err)
	assert.Equal(t, res, int64(-8934067919247763456))

}
