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

package monitor

const (
	EIP_UN_USED  = "EIP_UNUSED"
	DISK_UN_USED = "DISK_UNUSED"
	LB_UN_USED   = "LB_UNUSED"

	DRIVER_ACTION = "DELETE"

	EIP_UNUSED_START_DELETE = "start_delete"
	EIP_UNUSED_DELETE_FAIL  = "delete_fail"
)

type MonitorSuggest string

type MonitorResourceType string

const (
	EIP_MONITOR_RES_TYPE  = MonitorResourceType("弹性EIP")
	DISK_MONITOR_RES_TYPE = MonitorResourceType("云硬盘")
	LB_MONITOR_RES_TYPE   = MonitorResourceType("负载均衡实例")
)

const (
	EIP_MONITOR_SUGGEST  = MonitorSuggest("释放未使用的EIP")
	DISK_MONITOR_SUGGEST = MonitorSuggest("释放未使用的Disk")
	LB_MONITOR_SUGGEST   = MonitorSuggest("释放未使用的LB")
)
