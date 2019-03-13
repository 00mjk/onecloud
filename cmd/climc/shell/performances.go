package shell

import (
	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
)

func init() {
	type Top5Options struct {
		NODE_LABELS string `help:"Service tree tree-node labels"`
	}
	R(&Top5Options{}, "performance-top5", "Show performance top5", func(s *mcclient.ClientSession, args *Top5Options) error {
		params := jsonutils.NewDict()
		params.Add(jsonutils.NewString(args.NODE_LABELS), "node_labels")

		result, err := modules.Performances.GetTop5(s, params)
		if err != nil {
			return err
		}
		printObject(result)
		return nil
	})
}
