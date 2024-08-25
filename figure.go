package wzrpc

import "context"

type Figure[TReq, TResp any] string

func (f Figure[TReq, TResp]) Str() string {
	return string(f)
}

func (f Figure[TReq, TResp]) Go(l Leader, opts any) (resp TResp, err error) {

	if r, err := LeaderAction(l, f.Str(), opts, 0); err != nil {
		return resp, err
	} else {
		return f.norm(r)
	}
}

func (f Figure[TReq, TResp]) GoWithCtx(ctx context.Context, l Leader, opts any) (resp TResp, err error) {
	stopped := false
	defer func() { stopped = true }()

	type respBundle struct {
		Resp TResp
		Err  error
	}

	ch := make(chan respBundle)
	go func() {
		resp, err := f.Go(l, opts)
		if !stopped {
			ch <- respBundle{resp, err}
		}
	}()

	select {
	case <-ctx.Done():
		return resp, ctx.Err()
	case r := <-ch:
		return r.Resp, r.Err
	}
}

func (f Figure[TReq, TResp]) norm(data any) (resp TResp, err error) {
	if err = deployDataTo(data, &resp); err != nil {
		return resp, ErrUnableToDeployData
	}
	return resp, nil
}
