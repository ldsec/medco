package serviceMedCo_test

import (
	"github.com/lca1/medco/services"
	"github.com/lca1/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
	"gopkg.in/dedis/onet.v1/network"
	"strconv"
	"testing"
)

func getParam(nbHosts int) (*onet.Roster, *onet.LocalTest) {

	log.SetDebugVisible(1)
	local := onet.NewLocalTest()
	// generate 3 hosts, they don't connect, they process messages, and they
	// don't register the tree or entity list
	_, el, _ := local.GenTree(nbHosts, true)

	// get query parameters
	return el, local
}

func getClients(nbHosts int, el *onet.Roster) []*serviceMedCo.API {
	clients := make([]*serviceMedCo.API, nbHosts)
	for i := 0; i < nbHosts; i++ {
		clients[i] = serviceMedCo.NewMedCoClient(el.List[i], strconv.Itoa(i))
	}

	return clients
}

func getQueryParams(nbQp int, encKey abstract.Point) libUnLynx.CipherVector {
	listQueryParameters := make(libUnLynx.CipherVector, 0)

	for i := 0; i < nbQp; i++ {
		listQueryParameters = append(listQueryParameters, *libUnLynx.EncryptInt(encKey, int64(i)))
	}

	return listQueryParameters
}

func TestServiceDDT(t *testing.T) {
	el, local := getParam(3)
	clients := getClients(3, el)
	// test the query DDT with 100 query terms
	nbQp := 100
	qt := getQueryParams(nbQp, el.Aggregate)
	defer local.CloseAll()

	proofs := false

	var resultNode1, resultNode1Repeated, resultNode2, resultNode3 []libUnLynx.GroupingKey

	wg := libUnLynx.StartParallelize(len(el.List))

	serviceMedCo.NewMedCoClient(el.List[0], strconv.Itoa(0))

	// the first two threads execute the same operation (repetition) to check that in the end it yields the same result
	go func() {
		defer wg.Done()

		var err error
		_, resultNode1, _, err = clients[0].SendSurveyDDTRequestTerms(el, serviceMedCo.SurveyID("testDDTSurvey_node1"), qt, proofs, true)

		if err != nil {
			t.Fatal("Client", clients[0], " service did not start: ", err)
		}
	}()
	go func() {
		defer wg.Done()

		var err error
		_, resultNode1Repeated, _, err = clients[0].SendSurveyDDTRequestTerms(el, serviceMedCo.SurveyID("testDDTSurvey_node1_repeated"), qt, proofs, true)

		if err != nil {
			t.Fatal("Client", clients[0], " service did not start: ", err)
		}
	}()
	go func() {
		defer wg.Done()

		var err error
		_, resultNode2, _, err = clients[1].SendSurveyDDTRequestTerms(el, serviceMedCo.SurveyID("testDDTSurvey_node2"), qt, proofs, true)

		if err != nil {
			t.Fatal("Client", clients[1], " service did not start: ", err)
		}
	}()

	var err error
	_, resultNode3, _, err = clients[2].SendSurveyDDTRequestTerms(el, serviceMedCo.SurveyID("testDDTSurvey_node3"), qt, proofs, true)

	if err != nil {
		t.Fatal("Client", clients[2], " service did not start: ", err)
	}

	libUnLynx.EndParallelize(wg)

	assert.Equal(t, len(resultNode1), len(qt))
	assert.Equal(t, len(resultNode2), len(qt))
	assert.Equal(t, len(resultNode3), len(qt))

	assert.Equal(t, resultNode1, resultNode1Repeated)

}

func TestServiceAgg(t *testing.T) {
	el, local := getParam(3)
	clients1 := getClients(3, el)
	clients2 := getClients(3, el)
	defer local.CloseAll()

	proofs := false

	secKey1, pubKey1 := libUnLynx.GenKey()
	secKey2, pubKey2 := libUnLynx.GenKey()
	secKey3, pubKey3 := libUnLynx.GenKey()

	aggregate1 := libUnLynx.EncryptInt(el.Aggregate, int64(2))
	aggregate2 := libUnLynx.EncryptInt(el.Aggregate, int64(1))
	aggregate3 := libUnLynx.EncryptInt(el.Aggregate, int64(3))

	aggregate4 := libUnLynx.EncryptInt(el.Aggregate, int64(4))
	aggregate5 := libUnLynx.EncryptInt(el.Aggregate, int64(5))
	aggregate6 := libUnLynx.EncryptInt(el.Aggregate, int64(6))

	var resultNode1, resultNode2, resultNode3, resultNode4, resultNode5, resultNode6 libUnLynx.CipherText

	wg := libUnLynx.StartParallelize(len(el.List) * 2)

	// the first two threads execute the same operation (repetition) to check that in the end it yields the same result
	// surveyID should be the same
	go func() {
		defer wg.Done()

		var err error
		_, resultNode1, _, err = clients1[0].SendSurveyAggRequest(el, serviceMedCo.SurveyID("testAggSurvey1"), pubKey1, *aggregate1, proofs)

		if err != nil {
			t.Fatal("Client", clients1[0], " service did not start: ", err)
		}
	}()
	go func() {
		defer wg.Done()

		var err error
		_, resultNode2, _, err = clients1[1].SendSurveyAggRequest(el, serviceMedCo.SurveyID("testAggSurvey1"), pubKey2, *aggregate2, proofs)

		if err != nil {
			t.Fatal("Client", clients1[1], " service did not start: ", err)
		}
	}()
	go func() {
		defer wg.Done()

		var err error
		_, resultNode3, _, err = clients1[2].SendSurveyAggRequest(el, serviceMedCo.SurveyID("testAggSurvey1"), pubKey3, *aggregate3, proofs)

		if err != nil {
			t.Fatal("Client", clients1[2], " service did not start: ", err)
		}
	}()

	go func() {
		defer wg.Done()

		var err error
		_, resultNode4, _, err = clients2[0].SendSurveyAggRequest(el, serviceMedCo.SurveyID("testAggSurvey2"), pubKey1, *aggregate4, proofs)

		if err != nil {
			t.Fatal("Client", clients2[0], " service did not start: ", err)
		}
	}()
	go func() {
		defer wg.Done()

		var err error
		_, resultNode5, _, err = clients2[1].SendSurveyAggRequest(el, serviceMedCo.SurveyID("testAggSurvey2"), pubKey2, *aggregate5, proofs)

		if err != nil {
			t.Fatal("Client", clients2[1], " service did not start: ", err)
		}
	}()
	go func() {
		defer wg.Done()

		var err error
		_, resultNode6, _, err = clients2[2].SendSurveyAggRequest(el, serviceMedCo.SurveyID("testAggSurvey2"), pubKey3, *aggregate6, proofs)

		if err != nil {
			t.Fatal("Client", clients2[2], " service did not start: ", err)
		}
	}()

	libUnLynx.EndParallelize(wg)

	// Check result
	listResults1 := make([]int64, 0)
	listResults1 = append(listResults1, libUnLynx.DecryptInt(secKey1, resultNode1), libUnLynx.DecryptInt(secKey2, resultNode2), libUnLynx.DecryptInt(secKey3, resultNode3))

	assert.Contains(t, listResults1, int64(2))
	assert.Contains(t, listResults1, int64(1))
	assert.Contains(t, listResults1, int64(3))

	listResults2 := make([]int64, 0)
	listResults2 = append(listResults2, libUnLynx.DecryptInt(secKey1, resultNode4), libUnLynx.DecryptInt(secKey2, resultNode5), libUnLynx.DecryptInt(secKey3, resultNode6))

	assert.Contains(t, listResults2, int64(4))
	assert.Contains(t, listResults2, int64(5))
	assert.Contains(t, listResults2, int64(6))

}

func TestCheckDDTSecrets(t *testing.T) {
	addr := network.NewLocalAddress("local://127.0.0.1:2020")
	_, err := serviceMedCo.CheckDDTSecrets("secrets.toml", addr)
	assert.Nil(t, err, "Error while writing the secrets to the TOML file")

	addr = network.NewLocalAddress("local://127.0.0.1:2010")
	_, err = serviceMedCo.CheckDDTSecrets("secrets.toml", addr)
	assert.Nil(t, err, "Error while writing the secrets to the TOML file")

	addr = network.NewLocalAddress("local://127.0.0.1:2000")
	_, err = serviceMedCo.CheckDDTSecrets("secrets.toml", addr)
	assert.Nil(t, err, "Error while writing the secrets to the TOML file")
}
