package unlynx

import (
	"github.com/ldsec/medco/connector/util/server"
	libunlynx "github.com/ldsec/unlynx/lib"
	"github.com/sirupsen/logrus"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/app"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// warning: needs medco deployment dev-3nodes-local running

var cothorityRoster *onet.Roster

func init() {
	utilserver.SetLogLevel("5")
	utilserver.MedCoNodeIdx = 0

	_, filename, _, _ := runtime.Caller(0)
	utilserver.UnlynxGroupFilePath = filepath.Dir(filename) + "/test_group.toml"

	groupFile, err := os.Open(utilserver.UnlynxGroupFilePath)
	if err != nil {
		logrus.Panic("unlynx error opening group file: ", err)
	}

	group, err := app.ReadGroupDescToml(groupFile)
	if err != nil || len(group.Roster.List) <= 0 {
		logrus.Panic("unlynx error parsing group file: ", err)
	}

	cothorityRoster = group.Roster
}

func TestGetQueryTermsDDT(t *testing.T) {
	encryptedInt0, err := libunlynx.EncryptInt(cothorityRoster.Aggregate, 366).Serialize()
	if err != nil {
		t.Fail()
	}

	encryptedInt1, err := libunlynx.EncryptInt(cothorityRoster.Aggregate, 2).Serialize()
	if err != nil {
		t.Fail()
	}

	t.Log(encryptedInt0, encryptedInt1)

	tags, _, err := DDTagValues("test query "+time.Now().Format(time.RFC3339Nano), []string{
		encryptedInt0,
		encryptedInt1,
	})

	if err != nil {
		t.Fail()
	}
	t.Log(tags)
}

// todo: need to submit to all 3 nodes
//func TestAggregateAndKeySwitchDummyFlags(t *testing.T) {
//	privKey, pubKey := libunlynx.GenKey()
//	pubKeySer, err := libunlynx.SerializePoint(pubKey)
//	if err != nil {
//		t.Fail()
//	}
//
//	value := make(chan int64)
//	for i := 0 ; i < 3 ; i++ {
//
//		//fmt.Println("id:", i)
//		//fmt.Println(cothorityRoster.List[i])
//		iCp := i
//		go func() {
//			util.UnlynxGroupFileIdx = iCp
//			fmt.Println("id:", iCp)
//			fmt.Println(cothorityRoster.List[iCp])
//			setupUnlynxClient()
//
//			agg, err := AggregateAndKeySwitchDummyFlags(
//				"test query " + time.Now().Format(time.RFC3339Nano),
//				[]string{
//					libunlynx.EncryptInt(cothorityRoster.Aggregate, 0).Serialize(),
//					libunlynx.EncryptInt(cothorityRoster.Aggregate, 1).Serialize(),
//				},
//				pubKeySer,
//			)
//
//			if err != nil {
//				t.Fail()
//			}
//
//			aggDes := &libunlynx.CipherText{}
//			err = aggDes.Deserialize(agg)
//			if err != nil {
//				t.Fail()
//			}
//
//
//			value <-libunlynx.DecryptInt(privKey, *aggDes)
//		}()
//
//
//	}
//	t.Log(<-value)
//	t.Log(<-value)
//	t.Log(<-value)
//}
