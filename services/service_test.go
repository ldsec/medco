package servicesmedco_test

import (
	"github.com/lca1/medco-unlynx/services"
	"github.com/lca1/unlynx/lib"
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

	log.SetDebugVisible(1)
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
	el, local := getParam(10)
	// test with 10 concurrent clients
	clients := getClients(10, el)
	// the first two threads execute the same operation (repetition) to check that in the end it yields the same result
	clients[1] = clients[0]
	// test the query DDT with 500 query terms
	nbQp := 100
	qt := getQueryParams(nbQp, el.Aggregate)
	defer local.CloseAll()

	proofs := false

	results := make(map[string][]libunlynx.GroupingKey)

	wg := libunlynx.StartParallelize(len(clients))
	var mutex = sync.Mutex{}
	for i, client := range clients {
		go func(i int, client *servicesmedco.API) {
			defer wg.Done()

			_, res, _, err := client.SendSurveyDDTRequestTerms(el, servicesmedco.SurveyID("testDDTSurvey_"+client.ClientID), qt, proofs, true)
			mutex.Lock()
			results["testDDTSurvey_"+client.ClientID] = res
			mutex.Unlock()

			if err != nil {
				t.Fatal("Client", client.ClientID, " service did not start: ", err)
			}
		}(i, client)
	}
	libunlynx.EndParallelize(wg)

	for _, result := range results {
		assert.Equal(t, len(qt), len(result))
	}
	assert.Equal(t, results["testDDTSurvey_"+clients[0].ClientID], results["testDDTSurvey_"+clients[1].ClientID])
}

func TestServiceAgg(t *testing.T) {
	// test with 10 servers
	el, local := getParam(10)
	// test with 10 concurrent clients
	nbHosts := 10
	clients1 := getClients(nbHosts, el)
	//clients2 := getClients(nbHosts, el)
	defer local.CloseAll()

	proofs := false

	secKeys := make([]kyber.Scalar, 0)
	pubKeys := make([]kyber.Point, 0)
	aggregates1 := make([]libunlynx.CipherText, 0)
	aggregates2 := make([]libunlynx.CipherText, 0)

	results1 := make([]libunlynx.CipherText, nbHosts)
	//results2 := make([]libunlynx.CipherText, nbHosts)

	for i := 0; i < nbHosts; i++ {
		_, sK, pK := libunlynx.GenKeys(1)
		secKeys = append(secKeys, sK[0])
		pubKeys = append(pubKeys, pK[0])

		aggregates1 = append(aggregates1, *libunlynx.EncryptInt(el.Aggregate, int64(i)))
		aggregates2 = append(aggregates2, *libunlynx.EncryptInt(el.Aggregate, int64(i)))
	}

	wg := libunlynx.StartParallelize(nbHosts)
	var mutex = sync.Mutex{}
	for i, client := range clients1 {
		go func(i int, client *servicesmedco.API) {
			defer wg.Done()

			_, res, _, err := client.SendSurveyAggRequest(el, servicesmedco.SurveyID("testAggSurvey1"), pubKeys[i], aggregates1[i], proofs)
			mutex.Lock()
			results1[i] = res
			mutex.Unlock()

			if err != nil {
				t.Fatal("Client", client.ClientID, " service did not start: ", err)
			}
		}(i, client)

	}

	/*for i,client := range clients2 {
		go func(i int, client *servicesmedco.API) {
			defer wg.Done()

			_, res, _, err := client.SendSurveyAggRequest(el, servicesmedco.SurveyID("testAggSurvey2"), pubKeys[i], aggregates2[i], proofs)
			mutex.Lock()
			results2[i]=res
			mutex.Unlock()

			if err != nil {
				t.Fatal("Client", client.ClientID, " service did not start: ", err)
			}
		}(i, client)

	}*/

	libunlynx.EndParallelize(wg)

	// Check result
	listResults1 := make([]int64, 0)
	for i, res := range results1 {
		listResults1 = append(listResults1, libunlynx.DecryptInt(secKeys[i], res))
	}

	assert.Contains(t, listResults1, int64(0))
	assert.Contains(t, listResults1, int64(1))
	assert.Contains(t, listResults1, int64(2))

	/*listResults2:= make([]int64, 0)
	for i,res := range results2 {
		listResults2 = append(listResults2, libunlynx.DecryptInt(secKeys[i], res))
	}

	assert.Contains(t, listResults2, int64(0))
	assert.Contains(t, listResults2, int64(1))
	assert.Contains(t, listResults2, int64(2))*/

}

func TestCheckDDTSecrets(t *testing.T) {
	addr := network.NewLocalAddress("local://127.0.0.1:2020")
	_, err := servicesmedco.CheckDDTSecrets("secrets.toml", addr)
	assert.Nil(t, err, "Error while writing the secrets to the TOML file")

	addr = network.NewLocalAddress("local://127.0.0.1:2010")
	_, err = servicesmedco.CheckDDTSecrets("secrets.toml", addr)
	assert.Nil(t, err, "Error while writing the secrets to the TOML file")

	addr = network.NewLocalAddress("local://127.0.0.1:2000")
	_, err = servicesmedco.CheckDDTSecrets("secrets.toml", addr)
	assert.Nil(t, err, "Error while writing the secrets to the TOML file")
}
