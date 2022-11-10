package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
	"nuclei/pkg/nuclei"
	"strings"
)

type PocService struct {
}

// Poc grpc服务
func (s *PocService) Poc(ctx context.Context, req *PocScanReq) (*PocScanResult, error) {
	res := pocScanResults(strings.Split(req.Urls, ","))
	return &PocScanResult{JsonResults: res}, nil
}

// PocService 输入urls 返回json
func pocScanResults(targets []string) string {
	options := &types.Options{
		Targets: targets,
	}
	if wrapper, err := nuclei.NewWrapper(options); err != nil {
		panic(err)
	} else {
		results, scanErr := wrapper.RunEnumeration(context.Background(), targets)
		if scanErr != nil {
			panic(scanErr)
		}
		if len(results) > 0 {
			var retList []PocResults
			for _, r := range results {
				retList = append(retList, PocResults{
					ID:   r.TemplateID,
					Name: r.Template,
				})
			}
			fmt.Println(retList)
			ret, err := json.Marshal(retList)
			if err != nil {
				panic(err)
			}
			return string(ret)
		}
	}
	return ""
}

type PocResults struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
