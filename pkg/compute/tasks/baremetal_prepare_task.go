package tasks

import (
	"context"
	"fmt"

	"yunion.io/x/jsonutils"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/compute/models"
)

type BaremetalPrepareTask struct {
	SBaremetalBaseTask
}

func init() {
	taskman.RegisterTask(BaremetalPrepareTask{})
}

func (self *BaremetalPrepareTask) OnInit(ctx context.Context, obj db.IStandaloneModel, body jsonutils.JSONObject) {
	baremetal := obj.(*models.SHost)
	url := fmt.Sprintf("/baremetals/%s/prepare", baremetal.Id)
	headers := self.GetTaskRequestHeader()
	self.SetStage("OnSyncConfigComplete", nil)
	_, err := baremetal.BaremetalSyncRequest(ctx, "POST", url, headers, self.Params)
	if err != nil {
		self.OnFailure(ctx, baremetal, err.Error())
	}
}

func (self *BaremetalPrepareTask) OnFailure(ctx context.Context, baremetal *models.SHost, reason string) {
	baremetal.SetStatus(self.UserCred, api.BAREMETAL_PREPARE_FAIL, reason)
	self.SetStageFailed(ctx, reason)
}

func (self *BaremetalPrepareTask) OnSyncConfigComplete(ctx context.Context, baremetal *models.SHost, body jsonutils.JSONObject) {
	baremetal.ClearSchedDescCache()
	self.SetStageComplete(ctx, nil)
}

func (self *BaremetalPrepareTask) OnSyncConfigCompleteFailed(ctx context.Context, baremetal *models.SHost, body jsonutils.JSONObject) {
	reason, _ := body.GetString("__reason__")
	self.OnFailure(ctx, baremetal, reason)
}
