package buildpack

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

const (
	contextKey = "context"
)

func Frontend(ctx context.Context, c client.Client) (*client.Result, error) {
	inputs, err := c.Inputs(ctx)
	if err != nil {
		return nil, err
	}

	context, ok := inputs[contextKey]
	if !ok {
		return nil, fmt.Errorf(`must provide frontend input "context"`)
	}
	delete(inputs, contextKey)

	def, err := context.Marshal(llb.WithCustomName("Detecting buildpack"))
	if err != nil {
		return nil, err
	}

	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, err
	}

	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}

	var keys []string
	for key := range inputs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var detect string
	for _, key := range keys {
		_, err = ref.StatFile(ctx, client.StatRequest{
			Path: key,
		})
		if err != nil {
			if os.IsNotExist(errors.Cause(err)) {
				continue
			}
			return nil, err
		}
		detect = key
		break
	}

	if detect == "" {
		return nil, fmt.Errorf("failed to detect buildpack")
	}

	st := inputs[detect]
	def, err = st.Marshal()
	if err != nil {
		return nil, err
	}

	return c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
}
