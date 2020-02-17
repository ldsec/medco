package survivalserver

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"errors"
)

type BatchIterator struct {
	timeCodes         []string
	length            int
	batchNumber       int
	batchSize         float64
	currentBatchIndex int
	//currentStateLower int
	//currentStateUpper int
	endReached bool
}

func NewBatchItertor(timePoints []string, batchNumber int) (batches *BatchIterator, err error) {
	length := len(timePoints)
	if length == 0 {
		err = errors.New("Input array must contain at least 1 time code")
		return
	}
	if batchNumber < 1 {
		err = errors.New("Number of batch should be at least 1")
		return
	}

	if batchNumber > length {
		logrus.Info(fmt.Sprintf("Batch number %d higher than lenght of time points array %d. Changing batch number to %d to avoid empty batch", batchNumber, length, length))
		batchNumber = length
	}

	batches = &BatchIterator{
		timeCodes:   timePoints,
		length:      length,
		batchNumber: batchNumber,
		batchSize:   float64(length) / float64(batchNumber),
	}
	return
}

func (batches *BatchIterator) Next() (res []string) {
	resLower := int(math.Floor(float64(batches.currentBatchIndex) * batches.batchSize))
	resUpper := int(math.Floor(float64(batches.currentBatchIndex+1) * batches.batchSize))
	res = batches.timeCodes[resLower:resUpper]
	if batches.currentBatchIndex < batches.batchNumber-1 {
		batches.currentBatchIndex++
	} else if batches.endReached == false {
		batches.endReached = true
	}

	return res
}

func (batches *BatchIterator) Done() bool {
	return batches.endReached
}