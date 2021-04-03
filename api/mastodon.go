package handler

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/koding/multiconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var mastodon = ` data:image/jpeg;base64,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`
var twitterClient *twitter.Client

type Config struct {
	ClientID     string
	ClientSecret string
	TokenUrl     string `default:"https://api.twitter.com/oauth2/token"`
}

func init() {
	//load in conf from env vars
	conf := &Config{}
	multiconfig.New().Load(&conf)

	//construct twitter client
	config := &clientcredentials.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		TokenURL:     conf.TokenUrl,
	}

	httpClient := config.Client(oauth2.NoContext)
	twitterClient = twitter.NewClient(httpClient)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	usernames, ok := r.URL.Query()["username"]
	if !ok || len(usernames[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Add your twitter username as a query parameter: https://mastodon-in-twitter-avatar.vercel.app/api/mastodon?username=<YOUR_TWITTER_USERNAME>")
		return
	}
	usr, _, err := twitterClient.Users.Show(&twitter.UserShowParams{
		ScreenName: usernames[0],
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Oops")
		return
	}
	avatar := strings.Replace(usr.ProfileImageURLHttps, "_normal", "", 1)
	result, err := combineImages(avatar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Oops")
		return
	}
	err = png.Encode(w, result)
	if err != nil {
		fmt.Println(err)
	}
}

type ImageLayer struct {
	Image image.Image
	XPos  int
	YPos  int
}

func combineImages(imageUrl string) (result *image.RGBA, err error) {

	resp, err := http.Get(imageUrl)
	if err != nil {
		return nil, err
	}
	avatarImg, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	i := strings.Index(mastodon, ",")
	if i < 0 {
		return nil, err
	}
	// pass reader to NewDecoder
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(mastodon[i+1:]))
	mastodonImg, _, err := image.Decode(dec)
	if err != nil {
		return nil, err
	}
	//create image's background
	bgImg := image.NewRGBA(image.Rect(0, 0, avatarImg.Bounds().Dx(), avatarImg.Bounds().Dy()))

	//set the background color
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{color.Opaque}, image.ZP, draw.Src)

	//looping image layer, higher array index = upper layer
	for _, img := range []ImageLayer{
		{
			Image: avatarImg,
			XPos:  0,
			YPos:  0,
		},
		{
			Image: mastodonImg,
			XPos:  avatarImg.Bounds().Dx() - mastodonImg.Bounds().Dx(),
			YPos:  avatarImg.Bounds().Dy() - mastodonImg.Bounds().Dy(),
		},
	} {
		//set image offset
		offset := image.Pt(img.XPos, img.YPos)

		//combine the image
		draw.Draw(bgImg, img.Image.Bounds().Add(offset), img.Image, image.ZP, draw.Over)
	}
	return bgImg, nil

}
