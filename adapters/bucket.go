package adapters

import (
	"context"
	"io"
	"jiva-guildes/settings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type BucketManager struct {
	cld *cloudinary.Cloudinary
}

func (b *BucketManager) Setup(config interface{}) {
	cld, err := cloudinary.NewFromURL(settings.AppSettings.BUCKET_API_KEY)
	if err != nil {
		panic(err)
	}
	b.cld = cld
}

func (b *BucketManager) UploadFile(ctx context.Context, file io.Reader, params uploader.UploadParams) error {
	_, err := b.cld.Upload.Upload(ctx, file, params)
	if err != nil {
		return err
	}
	return nil
}
