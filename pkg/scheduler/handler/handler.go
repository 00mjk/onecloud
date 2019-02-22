package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/bitly/go-simplejson"
	"gopkg.in/gin-gonic/gin.v1"

	"yunion.io/x/log"
	"yunion.io/x/onecloud/pkg/scheduler/api"
	"yunion.io/x/onecloud/pkg/scheduler/core"
	"yunion.io/x/onecloud/pkg/scheduler/db/models"
	schedman "yunion.io/x/onecloud/pkg/scheduler/manager"
)

// InstallHandler is an interface that registes route and
// handles scheduler's services.
func InstallHandler(r *gin.Engine) {
	r.POST("/scheduler", timer(scheduleHandler))
	r.POST("/scheduler/:action", timer(schedulerActionHandler))
	r.POST("/scheduler/:action/:ident", timer(schedulerActionIdentHandler))
	InstallPingHandler(r)
	InstallVersionHandler(r)
	InstallK8sSchedExtenderHandler(r)
}

func timer(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		bytes, _ := httputil.DumpRequest(c.Request, true)
		log.V(10).Debugf(`
>>>>>>>>>>>>>
HTTP Request:
%s
>>>>>>>>>>>>>`, string(bytes))
		f(c)
		log.Infof("Handler %q cost: %v", c.Request.URL.Path, time.Since(startTime))
	}
}

func scheduleHandler(c *gin.Context) {
	doSyncSchedule(c)
}

func schedulerActionHandler(c *gin.Context) {
	act := c.Param("action")
	switch act {
	case "test":
		doSchedulerTest(c)
	case "candidate-list":
		doCandidateList(c)
	case "cleanup":
		doCleanup(c)
	case "history-list":
		doHistoryList(c)
	case "clean-cache":
		doCleanAllHostCache(c)
	//case "reserved-resources":
	//doReservedResources(c)
	default:
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("action: %s not support", act))
	}
}

func schedulerActionIdentHandler(c *gin.Context) {
	act := c.Param("action")
	id := c.Param("ident")
	switch act {
	case "clean-cache":
		doCleanHostCache(c, id)
	case "candidate-detail":
		doCandidateDetail(c, id)
	case "history-detail":
		doHistoryDetail(c, id)
	case "completed":
		doCompleted(c, id)
	default:
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("action: %s not support", act))
	}
}

func doSchedulerTest(c *gin.Context) {
	if !schedman.IsReady() {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Global scheduler not init"))
		return
	}

	sjson, err := simplejson.NewFromReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	schedInfo, err := api.NewSchedInfo(sjson, true)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	schedInfo.IsSuggestion = true
	result, err := schedman.Schedule(schedInfo)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transToSchedTestResult(result, schedInfo.SuggestionLimit))
}

func transToSchedTestResult(result *core.SchedResultItemList, limit int64) interface{} {
	return &api.SchedTestResult{
		Data:   result.Data,
		Total:  int64(result.Len()),
		Limit:  limit,
		Offset: 0,
	}
}

func doCandidateList(c *gin.Context) {
	sjson, err := simplejson.NewFromReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	args, err := api.NewCandidateListArgs(sjson)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := schedman.GetCandidateList(args)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func doCandidateDetail(c *gin.Context, id string) {
	hs, err := models.FetchHostByIDs([]string{id})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(hs) == 0 {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("Candidate %s not found.", id))
		return
	}

	host := hs[0].(*models.Host)

	args := new(api.CandidateDetailArgs)
	args.ID = id
	if !host.IsHypervisor() {
		args.Type = api.HostTypeBaremetal
	} else {
		args.Type = api.HostTypeHost
	}

	result, err := schedman.GetCandidateDetail(args)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func doCleanup(c *gin.Context) {
	sjson, err := simplejson.NewFromReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	args, err := api.NewCleanupArgs(sjson)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := schedman.Cleanup(args)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func doHistoryList(c *gin.Context) {
	sjson, err := simplejson.NewFromReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	args, err := api.NewHistoryArgs(sjson)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := schedman.GetHistoryList(args)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func doHistoryDetail(c *gin.Context, id string) {
	sjson, err := simplejson.NewFromReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	args, err := api.NewHistoryDetailArgs(sjson, id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := schedman.GetHistoryDetail(args)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func doSyncSchedule(c *gin.Context) {
	if !schedman.IsReady() {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Global scheduler not init"))
		return
	}
	sjson, err := simplejson.NewFromReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	schedInfo, err := api.NewSchedInfo(sjson, false)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := schedman.Schedule(schedInfo)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	count := schedInfo.Data.Count
	var resp interface{}
	if schedInfo.Data.Backup {
		resp = transToBackupSchedResult(result, schedInfo.Data.HostID, schedInfo.Data.BackupHostID, count)
	} else {
		resp = transToRegionSchedResult(result.Data, count)
	}

	c.JSON(http.StatusOK, resp)
}

func transToRegionSchedResult(result []*core.SchedResultItem, count int64) interface{} {
	apiResults := make([]api.SchedResultItem, 0)
	succCount := 0
	for _, nr := range result {
		for {
			if nr.Count <= 0 {
				break
			}
			tr := api.SchedSuccItem{
				Candidate: api.SchedNormalResultItem{
					ID:   nr.ID,
					Name: nr.Name,
					Data: nr.Data,
				},
			}
			apiResults = append(apiResults, tr)
			nr.Count--
			succCount++
		}
	}

	for {
		if int64(succCount) >= count {
			break
		}
		er := api.SchedErrItem{Error: "Out of resource"}
		apiResults = append(apiResults, er)
		succCount++
	}

	return regionResponse(apiResults)
}

func regionResponse(v interface{}) interface{} {
	return struct {
		Result interface{} `json:"scheduler"`
	}{Result: v}
}

func newExpireArgsByHostIDs(ids []string) (*api.ExpireArgs, error) {
	hs, err := models.FetchHostByIDs(ids)
	if err != nil {
		return nil, err
	}
	if len(hs) == 0 {
		return nil, fmt.Errorf("Hostscache %v not found", ids)
	}

	expireArgs := &api.ExpireArgs{
		DirtyBaremetals: []string{},
		DirtyHosts:      []string{},
	}
	for _, obj := range hs {
		host := obj.(*models.Host)
		if !host.IsHypervisor() {
			expireArgs.DirtyBaremetals = append(expireArgs.DirtyBaremetals, host.ID)
		} else {
			expireArgs.DirtyHosts = append(expireArgs.DirtyHosts, host.ID)
		}
	}
	return expireArgs, nil
}

func doCleanAllHostCache(c *gin.Context) {
	ids, err := models.AllIDs(models.Hosts)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	args, err := newExpireArgsByHostIDs(ids)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	doCleanHostCacheByArgs(c, args)
}

func doCleanHostCache(c *gin.Context, hostID string) {
	args, err := newExpireArgsByHostIDs([]string{hostID})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	doCleanHostCacheByArgs(c, args)
}

func doCleanHostCacheByArgs(c *gin.Context, args *api.ExpireArgs) {
	result, err := schedman.Expire(args)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, regionResponse(result))
}

func doCompleted(c *gin.Context, id string) {
	sjson, err := simplejson.NewFromReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	completedNotifyArgs, err := api.NewCompletedNotifyArgs(sjson, id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := schedman.CompletedNotify(completedNotifyArgs)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
