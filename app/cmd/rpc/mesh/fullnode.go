package mesh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alitto/pond"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/pokt-network/pocket-core/app"
	"github.com/puzpuzpuz/xsync"
	"github.com/robfig/cron/v3"
	"io"
	"io/ioutil"
	log2 "log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// fullNode - represent the pocket client instance running that could handle 1 or N addresses (lean)
type fullNode struct {
	Name             string
	URL              string
	Servicers        *xsync.MapOf[string, *servicer]
	Status           *app.HealthResponse
	BlocksPerSession int64
	Worker           *pond.WorkerPool
	MetricsWorker    *Metrics
	Crons            *cron.Cron
}

// NewWorker - generate a new worker.
func (node *fullNode) NewWorker() {
	node.Worker = NewWorkerPool(
		node.URL,
		app.GlobalMeshConfig.ServicerWorkerStrategy,
		app.GlobalMeshConfig.ServicerMaxWorkers,
		app.GlobalMeshConfig.ServicerMaxWorkersCapacity,
		app.GlobalMeshConfig.ServicerWorkersIdleTimeout,
	)

	node.MetricsWorker = NewWorkerPoolMetrics(
		node.Name,
		node.Worker,
	)
}

// start - start worker and cron jobs
func (node *fullNode) start() {
	logger.Debug(fmt.Sprintf("starting node %s with %d servicers", node.URL, node.Servicers.Size()))
	node.Crons.Start()
	node.NewWorker()
	node.MetricsWorker.Start()
}

// stop - stop worker and crons jobs from node
func (node *fullNode) stop() {
	logger.Debug(fmt.Sprintf("stopping worker pool of node %s", node.URL))
	node.Worker.Stop()
	logger.Debug(fmt.Sprintf("worker pool of node %s stopped!", node.URL))

	logger.Debug(fmt.Sprintf("stopping health cron job of node %s", node.URL))
	node.Crons.Stop()
	logger.Debug(fmt.Sprintf("check cron job of node %s stopped!", node.URL))

	logger.Debug(fmt.Sprintf("stopping metrics of node %s", node.URL))
	node.MetricsWorker.Stop()
	logger.Debug(fmt.Sprintf("metrics of node %s stopped!", node.URL))
}

// checkNodeEndpoint - check node endpoint
func (node *fullNode) checkNodeEndpoint(endpoint string) error {
	requestURL := fmt.Sprintf(
		"%s%s?verify=true",
		node.URL,
		endpoint,
	)
	req, err := http.NewRequest("POST", requestURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(AuthorizationHeader, servicerAuthToken.Value)
	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, err := servicerClient.Do(req)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}(resp.Body)

	// read the body just to allow http 1.x be able to reuse the connection
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		e := errors.New(fmt.Sprintf("Couldn't parse response body. Error: %s", err.Error()))
		return e
	}

	isSuccess := resp.StatusCode == 200

	if !isSuccess {
		return errors.New(
			fmt.Sprintf(
				"error=StatusCode != 200 code=%d",
				resp.StatusCode,
			),
		)
	}

	return nil
}

// runCheck - check that node is able to work as expected
func (node *fullNode) runCheck() error {
	if node.Servicers.Size() == 0 {
		return errors.New(fmt.Sprintf("node %s has 0 servicers load.", node.URL))
	}

	//logger.Debug(fmt.Sprintf("checking node %s connectivity", node.URL))
	servicers := make([]string, 0)

	node.Servicers.Range(func(key string, s *servicer) bool {
		servicers = append(servicers, s.Address.String())
		return true
	})

	payload := CheckPayload{
		Servicers: servicers,
		Chains:    make([]string, 0),
	}

	for _, chain := range chains.M {
		payload.Chains = append(payload.Chains, chain.ID)
	}

	jsonData, e := json.Marshal(payload)
	if e != nil {
		return e
	}

	requestURL := fmt.Sprintf(
		"%s%s",
		node.URL,
		ServicerCheckEndpoint,
	)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	req.Header.Set(AuthorizationHeader, servicerAuthToken.Value)
	if err != nil {
		return err
	}

	if app.GlobalMeshConfig.UserAgent != "" {
		req.Header.Set("User-Agent", app.GlobalMeshConfig.UserAgent)
	}
	resp, err := servicerClient.Do(req)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return // add log here
		}
	}(resp.Body)

	// read the body just to allow http 1.x be able to reuse the connection
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 || !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		if resp.StatusCode != 200 {
			err = errors.New(fmt.Sprintf("node is returning a non 200 code response from %s. code=%d err=%s", requestURL, resp.StatusCode, string(body)))
		} else {
			err = errors.New(fmt.Sprintf("node is returning a non json response from %s. code=%d err=%s", requestURL, resp.StatusCode, string(body)))
		}

		return err
	}

	res := &CheckResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	node.Status = &res.Status
	node.BlocksPerSession = res.BlocksPerSession

	if len(res.WrongChains) > 0 {
		return errors.New(fmt.Sprintf("unable to validate following chains: %s", strings.Join(res.WrongChains[:], ",")))
	}

	if len(res.WrongServicers) > 0 {
		return errors.New(fmt.Sprintf("unable to validate following servicers: %s", strings.Join(res.WrongServicers[:], ",")))
	}

	return nil
}

// scheduleNodeChecks - schedule a node heal pooling
func (node *fullNode) scheduleNodeChecks() {
	_, err := node.Crons.AddFunc(fmt.Sprintf("@every %ds", app.GlobalMeshConfig.NodeCheckInterval), func() {
		e := node.runCheck()
		if e != nil {
			logger.Error(
				fmt.Sprintf(
					"node %s failed check with error=%s",
					node.URL,
					e.Error(),
				),
			)
		}
	})

	if err != nil {
		log2.Fatal(err)
	}
}

// GetLatestSessionBlockHeight - same as pocket core code, just a reimplementation base on the info return by fullNode check call.
func (node *fullNode) GetLatestSessionBlockHeight() (sessionBlockHeight int64) {
	// get the latest block height
	blockHeight := node.Status.Height
	// get the blocks per session
	blocksPerSession := node.BlocksPerSession
	// if block height / blocks per session remainder is zero, just subtract blocks per session and add 1
	if blockHeight%blocksPerSession == 0 {
		sessionBlockHeight = blockHeight - node.BlocksPerSession + 1
	} else {
		// calculate the latest session block height by diving the current block height by the blocksPerSession
		sessionBlockHeight = (blockHeight/blocksPerSession)*blocksPerSession + 1
	}
	return
}

// createNode - returns a fullNode instance
func createNode(urlStr, name string) *fullNode {
	nodeCronJobsWorker := cron.New()

	if name == "" {
		u, err := url.Parse(urlStr)
		if err != nil {
			logger.Error("Unable to parse url of node", "err", err)
			name = urlStr
		} else {
			name = u.Hostname()
		}
	}

	logger.Debug(fmt.Sprintf("new node name=%s url=%s", name, urlStr))

	node := &fullNode{
		Name:      name,
		URL:       urlStr,
		Servicers: xsync.NewMapOf[*servicer](),
		Status:    nil,
		// just a default because on each check request the node will return the right value for the height returned
		// TODO: implement a default from config as "fallback" before get the real node value
		BlocksPerSession: 4,
		Crons:            nodeCronJobsWorker,
	}

	node.scheduleNodeChecks()

	return node
}

// connectivityChecks - run check over critical endpoints that mesh node need to be able to reach on servicer
func connectivityChecks(onlyFor mapset.Set[string]) {
	logger.Info("start connectivity checks")
	totalNodes := nodesMap.Size()
	if onlyFor.Cardinality() > 0 {
		totalNodes = onlyFor.Cardinality()
	}
	connectivityWorkerPool := pond.New(
		totalNodes, totalNodes, pond.MinWorkers(totalNodes),
		pond.Strategy(pond.Eager()),
	)

	var success uint32

	endpoints := []string{ServicerRelayEndpoint, ServicerSessionEndpoint, ServicerCheckEndpoint}

	// check health for all the servicer nodes before start.
	nodesMap.Range(func(key string, node *fullNode) bool {
		// run this check only if something is sent, otherwise (like on first start) it will run for all the nodes.
		if onlyFor.Cardinality() > 0 && !onlyFor.Contains(key) {
			// skip node because
			return true
		}

		for _, endpoint := range endpoints {
			connectivityWorkerPool.Submit(func() {
				ep := endpoint
				e := node.checkNodeEndpoint(ep)
				if e != nil {
					// any connectivity error with the node will stop this mesh client
					log2.Fatal(fmt.Sprintf("unable to reach node %s at endpoint %s. error: %s", node.URL, ep, e.Error()))
				}
				success++
			})
		}
		return true
	})

	// Wait for all HTTP requests to complete.
	connectivityWorkerPool.StopAndWait()

	if success == 0 {
		logger.Error(fmt.Sprintf("any node was able to be reach at endpoints: %s", strings.Join(endpoints[:], ",")))
		log2.Fatal(fmt.Sprintf("nodes=%d; reachable=%d", totalNodes, success))
	}

	if int(success) != totalNodes*len(endpoints) {
		logger.Error(fmt.Sprintf("IMPORTANT!!! few endpoints on nodes are not reachable"))
		logger.Error("you should stop this and fix the connectivity before continue")
	}

	firstCheckWorker := pond.New(
		totalNodes, totalNodes, pond.MinWorkers(totalNodes),
		pond.IdleTimeout(time.Duration(app.GlobalMeshConfig.ServicerWorkersIdleTimeout)*time.Millisecond),
		pond.Strategy(pond.Eager()),
	)

	nodesMap.Range(func(key string, node *fullNode) bool {
		// it will not kill the process because sometimes there is errors on the node side that are solved without
		// need to restart mesh client, like routing one.
		firstCheckWorker.Submit(func() {
			// run first time node check
			e := node.runCheck()
			if e != nil {
				logger.Error(fmt.Sprintf("node %s fail check with: %s", node.URL, e.Error()))
			}

			// start node working
			node.start()
		})

		return true
	})

	firstCheckWorker.StopAndWait()

	logger.Info("connectivity check done")
}

// NodesSize - return how many nodes are registered under nodesMap
func NodesSize() int {
	return nodesMap.Size()
}
