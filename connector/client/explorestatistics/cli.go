package explorestatisticsclient

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	medcoclient "github.com/ldsec/medco/connector/client"
	medcomodels "github.com/ldsec/medco/connector/models"

	"github.com/ldsec/medco/connector/restapi/models"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
)

// ClientResultElement holds the information for the CLI whole susrvival analysis loop
type ClientResultElement struct {
	ClearTimePoint     string
	EncEventOfInterest string
	EncCensoringEvent  string
}

// ExecuteClientExploreStatistics creates an explore statistics form parameters, and makes a call to the API to executes this query
func ExecuteClientExploreStatistics(token, username, password, cohortPanelStr, timing, query string, nbBuckets int64, disableTLSCheck bool, resultFile, timersFile string) (err error) {

	err = inputValidation(query, cohortPanelStr, timing)
	if err != nil {
		logrus.Error(err)
		return
	}

	//initialize objects and channels
	clientTimers := medcomodels.NewTimers()
	var accessToken string
	tokenChan := make(chan string, 1)
	errChan := make(chan error)
	signal := make(chan struct{})
	wait := &sync.WaitGroup{}
	wait.Add(1)
	go func() {
		wait.Wait()
		signal <- struct{}{}
	}()

	// --- get token
	logrus.Info("requesting access token")
	go func() {
		defer wait.Done()
		accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
		if err != nil {
			errChan <- err
			return
		}
		tokenChan <- accessToken
		logrus.Info("access token received")
		logrus.Tracef("token %s", accessToken)

	}()

	var parameters *Parameters

	select {
	case err = <-errChan:
		logrus.Error(err)
		return
	case <-time.After(time.Duration(utilclient.TokenTimeoutSeconds) * time.Second):
		err = fmt.Errorf("timeout %d seconds", utilclient.TokenTimeoutSeconds)
		logrus.Error(err)
		return

	case <-signal:
		accessToken = <-tokenChan
		//TODO Parse multiple panels from the query
		var queryPanel []*models.PanelConceptItemsItems0
		queryPanel, _, err = medcoclient.ParseQueryItem(query)
		if err != nil {
			logrus.Error("while parsing start item: ", err.Error())
			return
		}

		err = panelValidation(queryPanel)
		if err != nil {
			logrus.Error("while validating start item", err)
			return
		}

		var cohortPanel []*models.Panel
		cohortPanel, err = medcoclient.ParseQueryString(cohortPanelStr)
		if err != nil {
			logrus.Error("Error parsing cohort panel ", err)
			return
		}

		concepts := []string{}

		//parsing concepts of the query

		//parsing modifiers of the query
		var modifiers []*modifier
		for _, panelItem := range queryPanel {
			concepts = append(concepts, *panelItem.QueryTerm)

			if mod := panelItem.Modifier; mod != nil {
				modifiers = append(modifiers, &modifier{
					ModifierKey: *mod.ModifierKey,
					AppliedPath: *mod.AppliedPath,
				})

				continue
			}

		}

		parameters = &Parameters{
			cohortPanel,
			concepts,
			modifiers,
			nbBuckets,
		}
	}

	// --- execute query
	timer := time.Now()
	logrus.Info("CLI executing explore stats query ", parameters.String())
	results, userPrivateKey, err := executeQuery(accessToken, parameters, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while executing explore statistics results: %s", err.Error())
		logrus.Error(err)
		return
	}
	clientTimers.AddTimers("medco-connector-explore-statistics-query-remote-execution", timer, nil)
	logrus.Info("query executed")
	logrus.Tracef("encrypted results: %+v", results)
	logrus.Tracef("timers: %v", results[0])

	// --- decrypt result
	timer = time.Now()
	logrus.Info("decrypting results")
	clearResults := make([]*NodeClearResults, len(results))
	var nodesTimers []medcomodels.Timers = make([]medcomodels.Timers, len(results))
	for idx, encryptedResults := range results {
		clearResults[idx], err = encryptedResults.Decrypt(userPrivateKey)
		if err != nil {
			err = fmt.Errorf("while decrypting explore statistics results: %s", err.Error())
			logrus.Error(err)
			return
		}

		nodesTimers[idx] = encryptedResults.Timers
	}
	clientTimers.AddTimers("medco-connector-decryptions", timer, nil)
	logrus.Info("results decrypted")
	logrus.Tracef("clear results: %+v", clearResults)

	// --- printing results
	printResults(clearResults, nodesTimers, clientTimers, resultFile, timersFile)

	logrus.Info("Operation completed")
	return

}

func executeQuery(accessToken string, parameters *Parameters, disableTLSCheck bool) (results []*EncryptedResults, userPrivateKey string, err error) {

	query, err := NewExploreStatisticsQuery(
		accessToken,
		parameters,
		disableTLSCheck,
	)
	userPrivateKey = query.userPrivateKey
	if err != nil {
		return
	}

	errChan := make(chan error)
	resultChan := make(chan struct {
		Results []*EncryptedResults
	}, 1)

	resTimeout := time.After(time.Duration(utilclient.ExploreStatisticsTimeoutSeconds) * time.Second)
	resultTicks := time.Tick(time.Duration(utilclient.WaitTickSeconds) * time.Second)
	go func() {
		results, err := query.Execute()
		if err != nil {
			logrus.Error(err)
			errChan <- err
			return
		}
		resultChan <- struct {
			Results []*EncryptedResults
		}{Results: results}

	}()

	tickTime := 0
resLoop:
	for {
		select {
		case <-resTimeout:
			err = fmt.Errorf("timeout %d", utilclient.ExploreStatisticsTimeoutSeconds)
			return
		case err = <-errChan:
			logrus.Error(err)
			return
		case res := <-resultChan:
			results = res.Results
			break resLoop
		case <-resultTicks:
			tickTime += int(utilclient.WaitTickSeconds)
			logrus.Infof("waiting for response (%d seconds)", tickTime)
		}
	}

	return

}

func printResults(nodesClearResults []*NodeClearResults, nodesTimers []medcomodels.Timers, clientTimers medcomodels.Timers, resultFile, timerFile string) (err error) {
	logrus.Info("printing results")
	csv, err := utilclient.NewCSV(resultFile)
	if err != nil {
		err = fmt.Errorf("while creating CSV file handler: %s", err)
		logrus.Error(err)
		return
	}
	err = csv.Write([]string{"lower bound", "higher bound", "count", "nodeIndex"})
	if err != nil {
		err = fmt.Errorf("while writing result headers:%s", err.Error())
		logrus.Error(err)
		return
	}
	for nodeIdx, nodeClearResults := range nodesClearResults {
		for _, clearResult := range nodeClearResults.Intervals {
			csv.Write([]string{
				*clearResult.LowerBound,
				*clearResult.HigherBound,
				strconv.FormatInt(*clearResult.Count, 10),
				strconv.Itoa(nodeIdx),
			})
			if err != nil {
				err = fmt.Errorf("while writing a record: %s", err.Error())
				logrus.Error(err)
				return
			}

		}

	}
	err = csv.Flush()
	if err != nil {
		err = fmt.Errorf("while flushing buffer to result file: %s", err.Error())
		logrus.Error(err)
		return
	}

	err = csv.Close()
	if err != nil {
		err = fmt.Errorf("while closing result file: %s", err.Error())
		logrus.Error(err)
		return
	}
	logrus.Info("results printed")

	// print timers
	logrus.Info("dumping timers")
	dumpCSV, err := utilclient.NewCSV(timerFile)
	if err != nil {
		err = fmt.Errorf("while creating CSV file handler: %s", err)
		logrus.Error(err)
		return
	}
	dumpCSV.Write([]string{"type", "timer_description", "duration_milliseconds"})
	if err != nil {
		err = fmt.Errorf("while writing headers for timer file: %s", err)
		logrus.Error(err)
		return
	}
	// each remote time profilings
	for nodeIdx, nodeTimers := range nodesTimers {
		serverSortedTimers := nodeTimers.SortTimers()
		for _, duration := range serverSortedTimers {
			dumpCSV.Write([]string{
				"server" + strconv.Itoa(nodeIdx),
				duration[0], //description
				duration[1], //duration in ms
			})
			if err != nil {
				err = fmt.Errorf("while writing record for timer file: %s", err)
				logrus.Error(err)
				return
			}
		}

	}
	// and local
	localSortedTimers := clientTimers.SortTimers()
	for _, duration := range localSortedTimers {
		dumpCSV.Write([]string{
			"client",
			duration[0], //description
			duration[1], //duration in ms
		})
		if err != nil {
			err = fmt.Errorf("while writing record for timer file: %s", err)
			logrus.Error(err)
			return
		}
	}

	err = dumpCSV.Flush()
	if err != nil {
		err = fmt.Errorf("while flushing timer file: %s", err)
		logrus.Error(err)
		return
	}
	logrus.Info()
	err = dumpCSV.Close()
	if err != nil {
		err = fmt.Errorf("while closing timer file: %s", err)
		logrus.Error(err)
		return
	}

	logrus.Info("timers dumped")
	return

}

func inputValidation(query, cohortPanels, timing string) error {
	if cohortPanels == "" {
		return fmt.Errorf("cohort panels -c is not set")
	}

	if timing == "" {
		return fmt.Errorf("the timing that will be used to define the cohort is not defined")
	}
	if query == "" {
		return fmt.Errorf("the analyte of the query is not set")
	}

	return nil
}

func panelValidation(panel []*models.PanelConceptItemsItems0) (err error) {
	if len(panel) == 0 {
		err = fmt.Errorf("panels are empty. Was the start item string empry ?")
		logrus.Error(err)
		return
	}
	if len(panel) > 1 {
		err = fmt.Errorf("multiple items retrieved from the item string. Was a file provided ? Only clear concept should be used")
		logrus.Error(err)
		return
	}
	if *(panel[0].Encrypted) {
		err = fmt.Errorf("encrypted concept found, only clear concept should be used")
		logrus.Error(err)
	}
	return
}
