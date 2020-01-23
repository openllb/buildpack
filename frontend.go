package buildpack

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
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

	res, err = c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})

	img := specs.Image{
		Config: specs.ImageConfig{
			Env:        st.Env(),
			Entrypoint: st.GetArgs(),
			WorkingDir: st.GetDir(),
		},
	}

	config, err := json.Marshal(img)
	if err != nil {
		return nil, err
	}

	res.AddMeta(exptypes.ExporterImageConfigKey, config)
	return res, nil
}
