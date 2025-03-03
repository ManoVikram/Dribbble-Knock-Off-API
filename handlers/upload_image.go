package handlers

import (
	"context"
	"os"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

func UploadToImageKitHandler(fileName string, base64Image string) (string, error) {
	var ctx = context.Background()

	ik := imagekit.NewFromParams(imagekit.NewParams{
		PublicKey:   os.Getenv("IMAGEKIT_PUBLIC_KEY"),
		PrivateKey:  os.Getenv("IMAGEKIT_PRIVATE_KEY"),
		UrlEndpoint: os.Getenv("IMAGEKIT_URL_ENDPOINT"),
	})

	response, err := ik.Uploader.Upload(ctx, base64Image, uploader.UploadParam{
		FileName: fileName,
	})

	return response.Data.Url, err
}
