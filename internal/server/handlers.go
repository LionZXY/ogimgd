package server

import (
	"context"
	"errors"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"
	"time"

	"github.com/nDmitry/ogimgd/internal/preview"
)

const timeout = 30 * time.Second

type drawer interface {
	Draw(ctx context.Context, opts preview.Options) (image.Image, error)
}

func getPreview(d drawer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(r.Context(), timeout)

		defer cancel()

		opts := preview.Options{
			CanvasW:    1200,
			CanvasH:    630,
			Opacity:    0.6,
			AvaD:       64,
			IconW:      48,
			IconH:      48,
			TitleSize:  76,
			AuthorSize: 36,
			LabelSize:  40,
			Quality:    84,
		}

		titleParam := r.URL.Query().Get("title")

		if titleParam == "" {
			handleBadRequest(w, errors.New("Missing required title parameter"))
			return
		}

		opts.Title = titleParam

		authorParam := r.URL.Query().Get("author")

		if authorParam == "" {
			handleBadRequest(w, errors.New("Missing required author parameter"))
			return
		}

		opts.Author = authorParam

		bgParam := r.URL.Query().Get("bg")

		if bgParam == "" {
			handleBadRequest(w, errors.New("Missing required bg parameter"))
			return
		}

		opts.BgURL = bgParam

		avaParam := r.URL.Query().Get("ava")

		if avaParam == "" {
			handleBadRequest(w, errors.New("Missing required ava parameter"))
			return
		}

		opts.AvaURL = avaParam

		iconParam := r.URL.Query().Get("icon")

		if iconParam == "" {
			handleBadRequest(w, errors.New("Missing required icon parameter"))
			return
		}

		opts.IconURL = iconParam

		labelLParam := r.URL.Query().Get("label_l")

		if labelLParam != "" {
			opts.LabelL = labelLParam
		}

		labelRParam := r.URL.Query().Get("label_r")

		if labelRParam != "" {
			opts.LabelR = labelRParam
		}

		opacityParam := r.URL.Query().Get("op")

		if opacityParam != "" {
			var err error

			if opts.Opacity, err = strconv.ParseFloat(opacityParam, 64); err != nil {
				handleBadRequest(w, errors.New("Could not parse op parameter"))
				return
			}
		}

		img, err := d.Draw(ctx, opts)

		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		jpeg.Encode(w, img, &jpeg.Options{Quality: opts.Quality})
	}
}
