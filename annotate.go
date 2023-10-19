package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/api/vision/v1"
)

// Annotate image from a url using Google Vision API
func Annotate(imageUrl string) ([]*vision.EntityAnnotation, error) {
	ctx := context.Background()

	// Create a Vision API service instance
	service, err := vision.NewService(ctx)
	if err != nil {
		return nil, err
	}

	// Fire a GET request to the url
	resp, err := http.Get(imageUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching image: %v", err)
	}
	defer resp.Body.Close()

	// Read data from response and convert it to base64-encoded string
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading image data: %v", err)
	}
	base64String := base64.StdEncoding.EncodeToString(imageBytes)

	// Prepare request to vision api
	req := &vision.AnnotateImageRequest{
		Image: &vision.Image{
			Content: base64String,
		},
		Features: []*vision.Feature{
			{
				Type:       "LABEL_DETECTION",
				MaxResults: 5,
			},
		},
	}

	// Call vision api
	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}
	res, err := service.Images.Annotate(batch).Do()
	if err != nil {
		return nil, fmt.Errorf("error calling vision API: %v", err)
	}
	
	annotations := res.Responses[0].LabelAnnotations
	if len(annotations) == 0 {
		return nil, fmt.Errorf("no label found: %s", imageUrl)		
	}

	return annotations, nil
}
