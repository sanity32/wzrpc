package wzrpc

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

func (f Figure[TReq, TResp]) norm(data any) (resp TResp, err error) {
	if err = deployDataTo(data, &resp); err != nil {
		return resp, ErrUnableToDeployData
	}
	return resp, nil
}
