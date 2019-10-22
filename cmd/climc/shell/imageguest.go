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

package shell

import (
	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
	"yunion.io/x/onecloud/pkg/mcclient/options"
)

func init() {

	type GuestImageCreateOptions struct {
		NAME        string `help:"Name of guest image"`
		ImageNumber int    `help:"common image number of guest image"`
		Protected   bool   `help:"if guest image is protected"`
	}

	R(&GuestImageCreateOptions{}, "guest-image-create", "Create guest image's metadata", func(s *mcclient.ClientSession,
		args *GuestImageCreateOptions) error {

		params := jsonutils.NewDict()
		params.Add(jsonutils.NewString(args.NAME), "name")
		if args.ImageNumber > 0 {
			params.Add(jsonutils.NewInt(int64(args.ImageNumber)), "image_number")
		}
		if args.Protected {
			params.Add(jsonutils.JSONTrue, "protected")
		}
		ret, err := modules.GuestImages.Create(s, params)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	},
	)

	type GuestImageListOptions struct {
		options.BaseListOptions

		Name string `help:"Name filter"`
	}

	R(&GuestImageListOptions{}, "guest-image-list", "List guest images", func(s *mcclient.ClientSession,
		args *GuestImageListOptions) error {

		params, err := args.Params()
		if err != nil {
			return err
		}

		if len(args.Name) > 0 {
			params.Add(jsonutils.NewString(args.Name), "name")
		}
		params.Add(jsonutils.JSONTrue, "details")
		rets, err := modules.GuestImages.List(s, params)
		if err != nil {
			return err
		}
		printList(rets, modules.GuestImages.GetColumns(s))
		return nil
	},
	)

	type GuestImageDeleteOptions struct {
		ID                    []string `help:"Image ID or name"`
		OverridePendingDelete *bool    `help:"Delete image directly instead of pending delete" short-token:"f"`
	}
	R(&GuestImageDeleteOptions{}, "guest-image-delete", "Delete a image", func(s *mcclient.ClientSession,
		args *GuestImageDeleteOptions) error {
		params, err := options.StructToParams(args)
		if err != nil {
			return err
		}
		ret := modules.GuestImages.BatchDeleteWithParam(s, args.ID, params, nil)
		printBatchResults(ret, modules.GuestImages.GetColumns(s))
		return nil
	})

	type GuestImageCancelDeleteOptions struct {
		ID string `help:"Guest Image id or name"`
	}
	R(&GuestImageCancelDeleteOptions{}, "guest-image-cancel-delete", "Cancel pending delete images",
		func(s *mcclient.ClientSession,
			args *GuestImageCancelDeleteOptions) error {
			if image, e := modules.GuestImages.PerformAction(s, args.ID, "cancel-delete", nil); e != nil {
				return e
			} else {
				printObject(image)
			}
			return nil
		})

	type GuestImageOptions struct {
		ID string `help:"Guest Image id or name"`
	}
	R(&GuestImageOptions{}, "guest-image-mark-protected", "Mark image protected", func(s *mcclient.ClientSession,
		args *GuestImageOptions) error {

		params := jsonutils.NewDict()
		params.Add(jsonutils.JSONTrue, "protected")
		result, err := modules.GuestImages.Update(s, args.ID, params)
		if err != nil {
			return err
		}
		printObject(result)
		return nil
	},
	)
	R(&GuestImageOptions{}, "guest-image-mark-unprotected", "Mark image protected", func(s *mcclient.ClientSession,
		args *GuestImageOptions) error {

		params := jsonutils.NewDict()
		params.Add(jsonutils.JSONFalse, "protected")
		result, err := modules.GuestImages.Update(s, args.ID, params)
		if err != nil {
			return err
		}
		printObject(result)
		return nil
	},
	)

}
