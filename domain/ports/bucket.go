package ports

import "context"

type BucketManager interface {
	Setup(interface{})
	UploadFile(context.Context, string, interface{}) error
}
