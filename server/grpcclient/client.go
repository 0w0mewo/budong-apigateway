package grpcclient

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/0w0mewo/budong-apigateway/server"
	"github.com/0w0mewo/budong-apigateway/server/setupb"
	"google.golang.org/grpc"
)

var _ server.Service = &SetuGrpcClient{}

type SetuGrpcClient struct {
	client setupb.SetuServiceClient
	conn   *grpc.ClientConn
}

func NewSetuGrpcClient(addr string) *SetuGrpcClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := setupb.NewSetuServiceClient(conn)
	return &SetuGrpcClient{
		client: c,
		conn:   conn,
	}
}

func (sgc *SetuGrpcClient) GetSetuFromDB(id int) ([]byte, error) {
	stream, err := sgc.client.GetSetuById(context.Background(), &setupb.SetuReq{Id: int64(id)})
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		buf.Write(chunk.Chunk)
	}

	return buf.Bytes(), err

}

func (sgc *SetuGrpcClient) RequestSetu(num int, isR18 bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := sgc.client.Fetch(ctx, &setupb.FetchReq{Amount: uint64(num), R18: isR18})
	if err != nil {
		return err
	}

	return nil

}

func (sgc *SetuGrpcClient) GetInventory(page, pageLimit uint64) ([]*server.SetuInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := sgc.client.GetInventory(ctx, &setupb.InventoryReq{
		Page: page, PageLimit: pageLimit})
	if err != nil {
		return nil, err
	}

	return inventoryToSetuInfos(resp), nil
}

func (sgc *SetuGrpcClient) RandomSetu() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stream, err := sgc.client.Random(ctx, &setupb.RandomReq{R18: false})
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		buf.Write(chunk.Chunk)
	}

	return buf.Bytes(), err

}

func (sgc *SetuGrpcClient) Count() uint64 {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cnt, err := sgc.client.Count(ctx, &setupb.CountReq{})
	if err != nil {
		return 0
	}

	return cnt.Cnt
}

func (sgc *SetuGrpcClient) Shutdown() {
	sgc.conn.Close()
}

func inventoryToSetuInfos(setus *setupb.InventoryResp) []*server.SetuInfo {
	res := make([]*server.SetuInfo, 0, len(setus.Info))

	for _, s := range setus.Info {
		res = append(res, &server.SetuInfo{
			Id:    int(s.Id),
			Title: s.Title,
			Uid:   int(s.Uid),
			Url:   s.Url,
			IsR18: s.IsR18,
		})
	}

	return res
}
