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

var mastodon = `data:image/jpeg;base64,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`
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
