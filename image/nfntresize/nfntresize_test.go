package nfntresize

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/jpeg"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

var _ imageserver_image.Processor = &Processor{}

func TestProcessor(t *testing.T) {
	prc := &Processor{}
	for _, tc := range []struct {
		params             imageserver.Params
		expectedWidth      int
		expectedHeight     int
		expectedParamError string
	}{
		// no size
		{
			params:         imageserver.Params{},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			params:         imageserver.Params{Param: imageserver.Params{}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":  0,
				"height": 0,
			}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		// with size
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width": 100,
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"height": 100,
			}},
			expectedHeight: 100,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":  100,
				"height": 100,
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		// interpolation
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":         100,
				"interpolation": "nearest_neighbor",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":         100,
				"interpolation": "bilinear",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":         100,
				"interpolation": "bicubic",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":         100,
				"interpolation": "mitchell_netravali",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":         100,
				"interpolation": "lanczos2",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":         100,
				"interpolation": "lanczos3",
			}},
			expectedWidth: 100,
		},
		// mode
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "resize",
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "thumbnail",
			}},
			expectedWidth:  100,
			expectedHeight: 79, // 819 * 100 / 1024
		},
		// error
		{
			params:             imageserver.Params{Param: "invalid"},
			expectedParamError: Param,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width": "invalid",
			}},
			expectedParamError: Param + ".width",
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"height": "invalid",
			}},
			expectedParamError: Param + ".height",
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width": -1,
			}},
			expectedParamError: Param + ".width",
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"height": -1,
			}},
			expectedParamError: Param + ".height",
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":         100,
				"interpolation": false,
			}},
			expectedParamError: Param + ".interpolation",
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width":         100,
				"interpolation": "invalid",
			}},
			expectedParamError: Param + ".interpolation",
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width": 100,
				"mode":  false,
			}},
			expectedParamError: Param + ".mode",
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width": 100,
				"mode":  "invalid",
			}},
			expectedParamError: Param + ".mode",
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
			im, err := imageserver_image.Decode(imageserver_testdata.Medium)
			if err != nil {
				t.Fatal(err)
			}
			im, err = prc.Process(im, tc.params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && err.Param == tc.expectedParamError {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedWidth != 0 && im.Bounds().Dx() != tc.expectedWidth {
				t.Fatalf("unexpected width: got %d, want %d", im.Bounds().Dx(), tc.expectedWidth)
			}
			if tc.expectedHeight != 0 && im.Bounds().Dy() != tc.expectedHeight {
				t.Fatalf("unexpected height: got %d, want %d", im.Bounds().Dy(), tc.expectedHeight)
			}
		}()
	}
}

func TestProcessorChange(t *testing.T) {
	prc := &Processor{}
	for _, tc := range []struct {
		params   imageserver.Params
		expected bool
	}{
		{
			params:   imageserver.Params{},
			expected: false,
		},
		{
			params:   imageserver.Params{Param: imageserver.Params{}},
			expected: false,
		},
		{
			params:   imageserver.Params{Param: 666},
			expected: true,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width": 100,
			}},
			expected: true,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"height": 100,
			}},
			expected: true,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"interpolation": "lanczos3",
			}},
			expected: true,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"mode": "resize",
			}},
			expected: true,
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"foo": "bar",
			}},
			expected: false,
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
			c := prc.Change(tc.params)
			if c != tc.expected {
				t.Fatalf("unexpected result: got %t, want %t", c, tc.expected)
			}
		}()
	}
}
