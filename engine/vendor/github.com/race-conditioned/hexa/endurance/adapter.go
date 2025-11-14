package endurance

import "github.com/race-conditioned/hexa/horizon/ports/inbound"

func Transport[C inbound.Ctx, Com inbound.Command, Res inbound.Result](
	usecase func(C, Com) (Res, error),
	pre func(ctx C, meta inbound.RequestMeta, cmd Com),
	post func(ctx C, meta inbound.RequestMeta, cmd Com, res Res, err error),
) inbound.UnaryHandler[C, Com, Res] {
	return func(ctx C, meta inbound.RequestMeta, cmd Com) (Res, error) {
		if pre != nil {
			pre(ctx, meta, cmd)
		}
		res, err := usecase(ctx, cmd)
		if post != nil {
			post(ctx, meta, cmd, res, err)
		}
		return res, err
	}
}
