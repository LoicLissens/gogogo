package ports

import "context"

type BucketManager interface {
	Setup(interface{}, context.Context)
	UploadFile(context.Context, string, interface{}) error
}
