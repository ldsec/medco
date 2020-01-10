package servicesmedco_test

import (
	"github.com/ldsec/medco-unlynx/services"
	"github.com/ldsec/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
	"strconv"
	"sync"
	"testing"
)

func getParam(nbServers int) (*onet.Roster, *onet.LocalTest) {

	log.SetDebugVisible(2)
	local := onet.NewLocalTest(libunlynx.SuiTe)
	// generate 3 hosts, they don't connect, they process messages, and they
	// don't register the tree or entity list
	_, el, _ := local.GenTree(nbServers, true)

	// get query parameters
	return el, local
}

func getClients(nbHosts int, el *onet.Roster) []*servicesmedco.API {
	clients := make([]*servicesmedco.API, nbHosts)
	for i := 0; i < nbHosts; i++ {
		clients[i] = servicesmedco.NewMedCoClient(el.List[i%len(el.List)], strconv.Itoa(i))
	}

	return clients
}

func getQueryParams(nbQp int, encKey kyber.Point) libunlynx.CipherVector {
	listQueryParameters := make(libunlynx.CipherVector, 0)

	for i := 0; i < nbQp; i++ {
		listQueryParameters = append(listQueryParameters, *libunlynx.EncryptInt(encKey, int64(i)))
	}

	return listQueryParameters
}

func TestServiceDDT(t *testing.T) {
	// test with 10 servers
	nbrServers := 3
	el, local := getParam(nbrServers)
	// test with as many clients as servers
	clients := getClients(nbrServers, el)
	// the first two threads execute the same operation (repetition) to check that in the end it yields the same result
	clients[1] = clients[0]
	// test the query DDT with 500 query terms
	nbQp := 100
	qt := getQueryParams(nbQp, el.Aggregate)
	defer local.CloseAll()

	proofs := false

	results := make(map[string][]libunlynx.GroupingKey)

	// sanitization tests
	// no SurveyID
	_, _, _, err := clients[0].SendSurveyDDTRequestTerms(el, "", qt, proofs, true)
	assert.Error(t, err)
	// no Roster
	emptyRoster := *el
	emptyRoster.List = nil
	_, _, _, err = clients[0].SendSurveyDDTRequestTerms(&emptyRoster, "testDDTSurvey_", qt, proofs, true)
	assert.Error(t, err)
	// no terms to tag
	_, _, _, err = clients[0].SendSurveyDDTRequestTerms(el, "testDDTSurvey_", nil, proofs, true)

	wg := libunlynx.StartParallelize(len(clients))
	var mutex = sync.Mutex{}
	for i, client := range clients {
		go func(i int, client *servicesmedco.API) {
			defer wg.Done()

			var err error
			_, res, tr, err := client.SendSurveyDDTRequestTerms(el, servicesmedco.SurveyID("testDDTSurvey_"+client.ClientID), qt, proofs, true)
			mutex.Lock()
			results["testDDTSurvey_"+client.ClientID] = res
			mutex.Unlock()

			if err != nil {
				t.Fatal("Client", client.ClientID, " service did not start: ", err)
			}
			log.Lvl1("Time:", tr.MapTR)
		}(i, client)
	}
	libunlynx.EndParallelize(wg)

	for _, result := range results {
		assert.Equal(t, len(qt), len(result))
	}
	assert.Equal(t, results["testDDTSurvey_"+clients[0].ClientID], results["testDDTSurvey_"+clients[1].ClientID])
}

func TestServiceKS(t *testing.T) {
	// test with 10 servers
	nbrServers := 3
	el, local := getParam(nbrServers)
	// test with as many clients as servers
	nbrClients := nbrServers
	clients := getClients(nbrClients, el)
	defer local.CloseAll()

	proofs := false

	secKeys := make([]kyber.Scalar, 0)
	pubKeys := make([]kyber.Point, 0)
	targetData := make(libunlynx.CipherVector, 0)
	results := make([][]int64, nbrClients)

	for i := 0; i < nbrClients; i++ {
		_, sK, pK := libunlynx.GenKeys(1)
		secKeys = append(secKeys, sK[0])
		pubKeys = append(pubKeys, pK[0])

		targetData = append(targetData, *libunlynx.EncryptInt(el.Aggregate, int64(i)))
	}

	// sanitization tests
	// no SurveyID
	_, _, _, err := clients[0].SendSurveyKSRequest(el, "", pubKeys[0], targetData, proofs)
	assert.Error(t, err)
	// no Roster
	emptyRoster := *el
	emptyRoster.List = nil
	_, _, _, err = clients[0].SendSurveyKSRequest(&emptyRoster, "testKSRequest", pubKeys[0], targetData, proofs)
	assert.Error(t, err)
	// no target pubKey
	_, _, _, err = clients[0].SendSurveyKSRequest(el, "testKSRequest", nil, targetData, proofs)
	assert.Error(t, err)
	// no terms to key switch
	_, _, _, err = clients[0].SendSurveyKSRequest(el, "testKSRequest", pubKeys[0], nil, proofs)
	assert.Error(t, err)

	wg := libunlynx.StartParallelize(nbrClients)
	var mutex = sync.Mutex{}
	for i, client := range clients {
		go func(i int, client *servicesmedco.API) {
			defer wg.Done()

			_, res, tr, err := client.SendSurveyKSRequest(el, servicesmedco.SurveyID("testKSRequest_"+client.ClientID), pubKeys[i], targetData, proofs)
			if err != nil {
				t.Fatal("Client", client.ClientID, " service did not start: ", err)
			}

			decRes := make([]int64, 0)
			for _, val := range res {
				decRes = append(decRes, libunlynx.DecryptInt(secKeys[i], val))
			}
			mutex.Lock()
			results[i] = decRes
			mutex.Unlock()
			log.Lvl1("Time:", tr.MapTR)
		}(i, client)
	}
	libunlynx.EndParallelize(wg)

	// Check result
	for _, res := range results {
		for i := 0; i < nbrClients; i++ {
			assert.Equal(t, res[i], int64(i))
		}
	}
}

func TestServiceAgg(t *testing.T) {
	// test with 10 servers
	nbrServers := 3
	el, local := getParam(nbrServers)
	// test with as many clients as servers
	nbrClients := nbrServers
	clients := getClients(nbrClients, el)
	defer local.CloseAll()

	proofs := false

	secKeys := make([]kyber.Scalar, 0)
	pubKeys := make([]kyber.Point, 0)
	targetData := *libunlynx.EncryptInt(el.Aggregate, int64(1))
	results := make([]int64, nbrClients)

	for i := 0; i < nbrClients; i++ {
		_, sK, pK := libunlynx.GenKeys(1)
		secKeys = append(secKeys, sK[0])
		pubKeys = append(pubKeys, pK[0])
	}

	// sanitization tests
	// no SurveyID
	_, _, _, err := clients[0].SendSurveyAggRequest(el, "", pubKeys[0], targetData, proofs)
	assert.Error(t, err)
	// no Roster
	emptyRoster := *el
	emptyRoster.List = nil
	_, _, _, err = clients[0].SendSurveyAggRequest(&emptyRoster, "testAggRequest", pubKeys[0], targetData, proofs)
	assert.Error(t, err)
	// no target pubKey
	_, _, _, err = clients[0].SendSurveyAggRequest(el, "testAggRequest", nil, targetData, proofs)
	assert.Error(t, err)
	// no terms to aggregate
	emptyData := libunlynx.CipherText{}
	_, _, _, err = clients[0].SendSurveyAggRequest(el, "testAggRequest", pubKeys[0], emptyData, proofs)
	assert.Error(t, err)

	wg := libunlynx.StartParallelize(nbrClients)
	var mutex = sync.Mutex{}
	for i, client := range clients {
		go func(i int, client *servicesmedco.API) {
			defer wg.Done()

			_, res, tr, err := client.SendSurveyAggRequest(el, "testAggRequest", pubKeys[i], targetData, proofs)
			if err != nil {
				t.Fatal("Client", client.ClientID, " service did not start: ", err)
			}

			mutex.Lock()
			results[i] = libunlynx.DecryptInt(secKeys[i], res)
			mutex.Unlock()
			log.Lvl1("Time:", tr.MapTR)
		}(i, client)
	}
	libunlynx.EndParallelize(wg)

	// Check result
	for _, res := range results {
		assert.Equal(t, res, int64(nbrServers))
	}
}

func TestServiceShuffle(t *testing.T) {
	// test with 10 servers
	nbrServers := 3
	el, local := getParam(nbrServers)
	// test with as many clients as servers
	nbHosts := nbrServers
	clients := getClients(nbHosts, el)
	defer local.CloseAll()

	proofs := false

	secKeys := make([]kyber.Scalar, 0)
	pubKeys := make([]kyber.Point, 0)
	targetData := make(libunlynx.CipherVector, 0)
	results := make([]int64, nbHosts)

	for i := 0; i < nbHosts; i++ {
		_, sK, pK := libunlynx.GenKeys(1)
		secKeys = append(secKeys, sK[0])
		pubKeys = append(pubKeys, pK[0])

		targetData = append(targetData, *libunlynx.EncryptInt(el.Aggregate, int64(i)))
	}

	// sanitization tests
	// no SurveyID
	_, _, _, err := clients[0].SendSurveyShuffleRequest(el, "", pubKeys[0], &targetData[0], proofs)
	assert.Error(t, err)
	// no Roster
	emptyRoster := *el
	emptyRoster.List = nil
	_, _, _, err = clients[0].SendSurveyShuffleRequest(&emptyRoster, "testShuffleRequest", pubKeys[0], &targetData[0], proofs)
	assert.Error(t, err)
	// no target pubKey
	_, _, _, err = clients[0].SendSurveyShuffleRequest(el, "testShuffleRequest", nil, &targetData[0], proofs)
	assert.Error(t, err)
	// no terms to aggregate
	_, _, _, err = clients[0].SendSurveyShuffleRequest(el, "testShuffleRequest", pubKeys[0], nil, proofs)
	assert.Error(t, err)

	wg := libunlynx.StartParallelize(nbHosts)
	var mutex = sync.Mutex{}
	for i, client := range clients {
		go func(i int, client *servicesmedco.API) {
			defer wg.Done()

			var err error
			_, res, tr, err := client.SendSurveyShuffleRequest(el,
				servicesmedco.SurveyID("testShuffleRequest"), pubKeys[i], &targetData[i], proofs)
			if err != nil {
				t.Fatal("Client", client.ClientID, " service did not start: ", err)
			}

			mutex.Lock()
			results[i] = libunlynx.DecryptInt(secKeys[i], res)
			mutex.Unlock()
			log.Lvl1(i, "Time:", tr.MapTR)
		}(i, client)
	}
	libunlynx.EndParallelize(wg)

	// Check result
	for i := 0; i < nbHosts; i++ {
		assert.Contains(t, results, int64(i))
	}
}

func TestCheckDDTSecrets(t *testing.T) {
	addr := network.NewLocalAddress("local://127.0.0.1:2020")
	_, err := servicesmedco.CheckDDTSecrets("secrets.toml", addr, nil)
	assert.Nil(t, err, "Error while writing the secrets to the TOML file")

	addr = network.NewLocalAddress("local://127.0.0.1:2010")
	_, err = servicesmedco.CheckDDTSecrets("secrets.toml", addr, nil)
	assert.Nil(t, err, "Error while writing the secrets to the TOML file")

	addr = network.NewLocalAddress("local://127.0.0.1:2000")
	_, err = servicesmedco.CheckDDTSecrets("secrets.toml", addr, nil)
	assert.Nil(t, err, "Error while writing the secrets to the TOML file")
}
