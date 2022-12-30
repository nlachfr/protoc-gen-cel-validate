package options

import "google.golang.org/protobuf/proto"

func Join(opts ...*Options) *Options {
	switch len(opts) {
	case 0:
		return nil
	case 1:
		return opts[0]
	}
	opt := proto.Clone(opts[0]).(*Options)
	if opts[1] != nil {
		if opts[0] == nil {
			opt = opts[1]
		} else {
			proto.Merge(opt, opts[1])
		}
	}
	return Join(append([]*Options{opt}, opts[1:]...)...)
}
