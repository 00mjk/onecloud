// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"context"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/sqlchemy"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/util/stringutils2"
)

type SGroupguestManager struct {
	SGroupJointsManager
	SGuestResourceBaseManager
}

var GroupguestManager *SGroupguestManager

func init() {
	db.InitManager(func() {
		GroupguestManager = &SGroupguestManager{
			SGroupJointsManager: NewGroupJointsManager(
				SGroupguest{},
				"guestgroups_tbl",
				"groupguest",
				"groupguests",
				GuestManager,
			),
		}
		GroupguestManager.SetVirtualObject(GroupguestManager)
	})
}

type SGroupguest struct {
	SGroupJointsBase

	Tag     string `width:"256" charset:"ascii" nullable:"true" list:"user" update:"user" create:"optional"` // Column(VARCHAR(256, charset='ascii'), nullable=True)
	GuestId string `width:"36" charset:"ascii" nullable:"false" list:"user" create:"required"`               // Column(VARCHAR(36, charset='ascii'), nullable=False)
}

func (manager *SGroupguestManager) GetSlaveFieldName() string {
	return "guest_id"
}

func (joint *SGroupguest) Master() db.IStandaloneModel {
	return db.JointMaster(joint)
}

func (joint *SGroupguest) Slave() db.IStandaloneModel {
	return db.JointSlave(joint)
}

func (self *SGroupguest) GetExtraDetails(
	ctx context.Context,
	userCred mcclient.TokenCredential,
	query jsonutils.JSONObject,
	isList bool,
) (api.GroupguestDetails, error) {
	return api.GroupguestDetails{}, nil
}

func (manager *SGroupguestManager) FetchCustomizeColumns(
	ctx context.Context,
	userCred mcclient.TokenCredential,
	query jsonutils.JSONObject,
	objs []interface{},
	fields stringutils2.SSortedStrings,
	isList bool,
) []api.GroupguestDetails {
	rows := make([]api.GroupguestDetails, len(objs))

	groupRows := manager.SGroupJointsManager.FetchCustomizeColumns(ctx, userCred, query, objs, fields, isList)

	guestIds := make([]string, len(rows))
	for i := range rows {
		rows[i] = api.GroupguestDetails{
			GroupJointResourceDetails: groupRows[i],
		}
		guestIds[i] = objs[i].(SGroupguest).GuestId
	}

	guestIdMaps, err := db.FetchIdNameMap2(GuestManager, guestIds)
	if err != nil {
		log.Errorf("FetchIdNameMap2 fail %s", err)
		return rows
	}

	for i := range rows {
		if name, ok := guestIdMaps[guestIds[i]]; ok {
			rows[i].Guest = name
			rows[i].Server = name
		}
	}

	return rows
}

func (self *SGroupguest) Delete(ctx context.Context, userCred mcclient.TokenCredential) error {
	return db.DeleteModel(ctx, userCred, self)
}

func (self *SGroupguest) Detach(ctx context.Context, userCred mcclient.TokenCredential) error {
	return db.DetachJoint(ctx, userCred, self)
}

func (self *SGroupguestManager) FetchByGuestId(guestId string) ([]SGroupguest, error) {
	q := self.Query().Equals("guest_id", guestId)
	return self.fetchByQuery(q)
}

func (self *SGroupguestManager) fetchByQuery(q *sqlchemy.SQuery) ([]SGroupguest, error) {
	joints := make([]SGroupguest, 0, 1)
	err := db.FetchModelObjects(self, q, &joints)
	if err != nil {
		return nil, err
	}
	return joints, err
}

func (self *SGroupguestManager) FetchByGroupId(groupId string) ([]SGroupguest, error) {
	q := self.Query().Equals("group_id", groupId)
	return self.fetchByQuery(q)
}

func (self *SGroupguestManager) Attach(ctx context.Context, groupId, guestId string) (*SGroupguest, error) {

	joint := &SGroupguest{}
	joint.GuestId = guestId
	joint.GroupId = groupId

	err := self.TableSpec().Insert(joint)
	if err != nil {
		return nil, err
	}
	joint.SetModelManager(self, joint)
	return joint, nil
}

func (manager *SGroupguestManager) ListItemFilter(
	ctx context.Context,
	q *sqlchemy.SQuery,
	userCred mcclient.TokenCredential,
	query api.GroupguestListInput,
) (*sqlchemy.SQuery, error) {
	var err error

	q, err = manager.SGroupJointsManager.ListItemFilter(ctx, q, userCred, query.GroupJointsListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SGroupJointsManager.ListItemFilter")
	}
	q, err = manager.SGuestResourceBaseManager.ListItemFilter(ctx, q, userCred, query.ServerFilterListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SGuestResourceBaseManager.ListItemFilter")
	}

	if len(query.Tag) > 0 {
		q = q.In("tag", query.Tag)
	}

	return q, nil
}

func (manager *SGroupguestManager) OrderByExtraFields(
	ctx context.Context,
	q *sqlchemy.SQuery,
	userCred mcclient.TokenCredential,
	query api.GroupguestListInput,
) (*sqlchemy.SQuery, error) {
	var err error

	q, err = manager.SGroupJointsManager.OrderByExtraFields(ctx, q, userCred, query.GroupJointsListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SGroupJointsManager.OrderByExtraFields")
	}
	q, err = manager.SGuestResourceBaseManager.OrderByExtraFields(ctx, q, userCred, query.ServerFilterListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SGuestResourceBaseManager.OrderByExtraFields")
	}

	return q, nil
}
