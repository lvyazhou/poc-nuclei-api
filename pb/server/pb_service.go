package pb

import (
	"context"
	"fmt"
)

type PocService struct {
}

func (s *PocService) Poc(ctx context.Context, req *PocReq) (*PocResult, error) {
	fmt.Println("poc service ... ")

	return &PocResult{JsonResults: req.Urls}, nil

}

func (s *PocService) mustEmbedUnimplementedRunPocServiceServer() {
	//TODO implement me
	panic("implement me")
}
