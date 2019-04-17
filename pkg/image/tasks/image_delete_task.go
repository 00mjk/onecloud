package tasks

import (
	"context"
	"fmt"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"

	api "yunion.io/x/onecloud/pkg/apis/image"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/image/models"
	"yunion.io/x/onecloud/pkg/image/options"
)

type ImageDeleteTask struct {
	taskman.STask
}

func init() {
	taskman.RegisterTask(ImageDeleteTask{})
}

func (self *ImageDeleteTask) OnInit(ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {
	image := obj.(*models.SImage)

	// imageStatus, _ := self.Params.GetString("image_status")
	isPurge := jsonutils.QueryBoolean(self.Params, "purge", false)
	isOverridePendingDelete := jsonutils.QueryBoolean(self.Params, "override_pending_delete", false)

	if options.Options.EnablePendingDelete && !isPurge && !isOverridePendingDelete {
		// imageStatus == models.IMAGE_STATUS_ACTIVE
		if image.PendingDeleted {
			self.SetStageComplete(ctx, nil)
			return
		}
		self.startPendingDeleteImage(ctx, image)
	} else {
		self.startDeleteImage(ctx, image)
	}
}

func (self *ImageDeleteTask) startPendingDeleteImage(ctx context.Context, image *models.SImage) {
	image.StopTorrents()
	image.DoPendingDelete(ctx, self.UserCred)
	self.SetStageComplete(ctx, nil)
}

func (self *ImageDeleteTask) startDeleteImage(ctx context.Context, image *models.SImage) {
	err := image.RemoveFiles()
	if err != nil {
		msg := fmt.Sprintf("fail to remove %s %s", image.GetPath(""), err)
		log.Errorf(msg)
		self.SetStageFailed(ctx, msg)
		return
	}

	image.SetStatus(self.UserCred, api.IMAGE_STATUS_DELETED, "delete")

	image.RealDelete(ctx, self.UserCred)

	self.SetStageComplete(ctx, nil)
}
