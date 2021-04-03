package main

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

var mastodon = `data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAASABIAAD/4QCMRXhpZgAATU0AKgAAAAgABQESAAMAAAABAAEAAAEaAAUAAAABAAAASgEbAAUAAAABAAAAUgEoAAMAAAABAAIAAIdpAAQAAAABAAAAWgAAAAAAAABIAAAAAQAAAEgAAAABAAOgAQADAAAAAQABAACgAgAEAAAAAQAAAFqgAwAEAAAAAQAAAE4AAAAA/+0AOFBob3Rvc2hvcCAzLjAAOEJJTQQEAAAAAAAAOEJJTQQlAAAAAAAQ1B2M2Y8AsgTpgAmY7PhCfv/AABEIAE4AWgMBEQACEQEDEQH/xAAfAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgv/xAC1EAACAQMDAgQDBQUEBAAAAX0BAgMABBEFEiExQQYTUWEHInEUMoGRoQgjQrHBFVLR8CQzYnKCCQoWFxgZGiUmJygpKjQ1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4eLj5OXm5+jp6vHy8/T19vf4+fr/xAAfAQADAQEBAQEBAQEBAAAAAAAAAQIDBAUGBwgJCgv/xAC1EQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkICQkKDA8MCgsOCwkJDRENDg8QEBEQCgwSExIQEw8QEBD/2wBDAQMDAwQDBAgEBAgQCwkLEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBD/3QAEAAz/2gAMAwEAAhEDEQA/AP07/wBZ/wBNvO/4D5+P/Qdv60AMmmhjhkuLiVDEVLyO7BFlVRkkk8IFHfvTjFydluJtRV2fL3xP/bGkhvbjSPhhp9pdwq2yTWb9GaO4x0MMAIJUdmZhnqAR1/Rsp4FU4KrmUmm/sR3X+KWtn5JPzaZ+b5vx44TdLLIppfblqn/himm15tpdVdHiuu/Hn4x+IpJX1H4h6rGk2N8VmUtY+OwEaggfia+vw/DmU4b4KEW+7vJ/i7fgfHYjiXN8S3z4iSXaNor8Ff8A8m+Zyd14k8SXxZr7xNrdyW+8ZtUuHz9cvXpQwmHp/BTivSMf8jzJYzEz+KrN+s5//JFB5ZZG3yTSu3955GZvzJJrdJJWS/AwcpN3bb9W2/vbuWoNY1m1/wCPXW9Uh/656hOv8nrOVCjP4oRfrGP+RrDE16fwVJr/ALfn/wDJHQaV8W/inoeRpfxF8QxKcZWS9M6nHTIm31wVsky3E/xcPB/K3/pNjuoZ5mmG/hYma9Xzf+lqR6n4M/bG8e6XcrD440yx8Q2MhAneCMWl3gdCCD5bkehCfWvnMfwNgq8W8JJ05efvR/8Akl+J9Nl/HmOoSUcZFVI+Xuy/Plf/AJKfWfhHxf4f8faDb+I/DV+uoWV8CM42O7r96N1PMToeCpxX5jjsDXy6vLD4iNpL+k0+qfRn6hgsbQzGhHEYaV4v+mmujXVPY2f9Z/0287/gPn4/9B2/rXIdYfbP+oz/AOS//wBagD//0P07/wBZ/wBNvO/4D5+P/Qdv60AfNH7Y3xNuNPsLL4Z6RdkPrEQvdWlQ4MlsrYiiI/hDsGJH8SoR0NfoXAuVRqzlmNVaQdo/4ur+S27N36H53x5m7pU45bSes1zS/wAKeiflJ790mup8lkknJOSa/Tj8u3O9+HPwO+JPxTja88KaGg09GMbahez+RbFxwVVsMzkHrtUgYIJB4rxszz/AZQ+TEz97+VK79Xqkvm7+Vj2cr4fzDOFz4WHu6+9J2jddFo2/krbq91Y7XW/2N/jLpNk95ZLoWrsiljBZ3rpKfZRKgVj/AMCWvKocb5XWnyz5oebSa+dndfcz2MRwNm1GDlDkn5JtP5XVn98TxO+s7vTLuew1K1mtLm1doriGdDG8Lr95XB5Uj37c9Oa+spzjVip03dPZrVP0PkKlOdGThUXK1uno1bv2/wAtdtT1HwJ+zF8W/HunRa1Z6TZ6Vp9woaCfVp2haVSOGWJVZ9p7Ftue3qfncw4ryzL6jpSk5yW6ir28rtpfdc+jy/hLNMxpqrGKhF7Obab8+VJu3q4vy7z+M/2VfjD4M0+XVm0ux1u0gUtK2kztLKigZLGF1ViB/slj7VGB4uyvHTVPmcG/5kkvvTa++y8y8dwdmuBpupyqaX8jbfrytJ/c5PyZ5ACrAMpBB5BByCK+natufLJpq6PYP2ZPiXceBviDbaFdzFtE8TSx2N3Cz7UW4PEE3sd2EJ7hlB6CvluLcpjmOAlWiv3lNNrzX2l/7cvR92fW8IZvLLsfGhN/u6rUX5S+y/8A2197xXRH3Zy/B/emX/gPn4/9B2/rX4wftIfbP+oyP/Af/wCtQB//0f07P7wf89vO/wCA+fj/ANB2/rQB+fn7RmsTa18a/FU80/mi1uYrCI4xhIoU4x/vO9fuPC1FUMooJdU5P1bf6JH4XxXWdfOK7fRqK9FFP85M5n4d+Em8d+O9A8HiVo01a/S3ldThkhALyke/lo4B7Eg16OZ4z+z8HVxW/JFteuiX4tHm5Xgv7RxtLCXspySfok5O3yTXle/Q/SpR4b8C+GeBa6TomiWhPGEitreNf0AUV+Dfv8fX6yqTfq23/mz9+/cYDD9IU4L0SSX4JI8h0r9sn4PanrUekyPrNlBPKsMd/d2Oy3yTgFuS6Lk/eZQB1OK+mrcE5pSpOqlFtK9lK7+XRvyT9D5elxxlNWqqTcopu3M4tLtr1S82kjd8e/s+eFPH/wASfD/xDvhGBp/OpWuwFdREfzW+8/7Ddf7y4B4FceX8R4nLsBVwMPtfC/5b/Fb1X3PU7cx4awuZY+ljqn2fiXSVtY3/AML189ndB8S/2lPhv8LNb/4RnVn1C/1SONJZ7bT7cSfZ1b7u9mKqCQMhc5xg4wRRlfC+Pzaj9YpWjDZOTte29lq/nsLNeKsvyit9WrNudk2oq9k9rvZbbXv5HXfDr4k+FPin4fHiPwleSS26ytBNFNGY5reUAEo6HocEEdQQQQSDXm5llmJymv7DEqz3VtU13TPUy3M8Nm1D2+FldbPo0+zT27+mux8bftdfD7T/AAT8S4tW0a2S3svE9s980Ma7VjukcLMQOwffG2P7xY96/UeDMxnjsA6VV3lTaV/7rV191mvS3Y/KuNctp4DMFVpK0aqcrf3k0pP/ALe5k/W76s8OeeW1U3duxWW3HnxsOqunzKfwZQa+vUVN8ktno/R6P8D42U5UoupF2cdV6rVfikfqBo162raPY6iy+YdQtYZyPu+cSgb/AIDtz+NfznXp+yqyh2bX3M/pOlP2lOM+6TLn2z/qMj/wH/8ArVkaH//S/Tv/AFn/AE287/gPn4/9B2/rQB+dPxsbf8YPGjeZvzrc/wA2MZ+SPtX7zw//AMirD/4F+bPwTiP/AJG+J/x/+2wOg/ZbAb48+FgQDj7YR7H7M/NcfFumTVv+3f8A0pHXwhrnVD0n/wCkn1z+1FK8XwD8YmNipayRDj0MqAj8q/NOE0pZzh7/AM36M/TOL5OOR4pr+Rn54Xah0nQ9GVwfxBr9vpuzT9D8Nrq8ZryZ+nvwyu57/wCHPhe9uXLyz6NZyOx6ljCuTX8+5pBU8dWhHZSl+bP6Jy2bq4OlOW7jH8kfnr8Z7qa8+L3jaa4cs/8Abt0mT/dXaij8FUD8K/bcjgqeV4dR/kX43Z+GZ9UlUzbFOX87/CMUfQf7B0shs/G8BY7FurFwvYExOCfyUflXxPiFFc+Hl5S/M+38OpN08VHpzR/9JRD+3eo3+CmxznUBn2xDVeHz/wB4X+H9ReIe2Gf+L8kfJd3/AMedxxn9y/8A6Ca/S6fxx9V+Z+YVv4U/R/kz9OvCI3+E9EGPM83TbXjOPOxEv/fO39a/nXG/7zU/xP8ANn9JYb+BD0X5Gt9s/wCoyP8AwH/+tXMbn//T/Tv/AFn/AE283/gPn4/9B2/rQB+dXxtz/wALh8aEvvzrUx3Yxn5I+1fvHD7vlWH/AMC/Nn4JxH/yN8T/AI//AG2Bv/stf8l58L/S9/8ASdq5OLf+RNW/7d/9KOvhD/kdUfSf5I+tv2pv+SBeL/8Ar0j/APRyV+a8Jf8AI5w/r+jP0vjD/kR4r/Cz89LgE+aB1Iav22PQ/D6uvN8z9E/hh8SPh7Y/DjwtZ3njrQILiDR7OOWKTUYVeNxEoKsC2QQeMGvw7Ncsxs8dWlGjJpzl9l935H7zleY4OGBoxlVimox+0uyPhP4pXdrf/E/xffWNxHcW1xrl3LDNEwZJELDDKRwQfUV+w5TGVPLqEJqzUI3R+L5xKM8zxM4u6c3Z99In0Z+wb/qfHH/XxYf+i5K+H8Qt8N6S/NH3Xh18OK/xR/8ASUM/bv6+Cv8Aev8A+UNHh/viP+3f1DxD2w3rL8kfJV5/x6XA/wCmMn/oJr9Lp/GvVfmfmFZ2pT9H+TP1C0C3+z6Fptpt8zFnBHj7vnbI1H/Adv61/OWIlz1py7t/mf0pRjy04x7JfkXvtn/UZH/gP/8AWrE1P//U/Ts/vB/z283/AID5+P8A0Hb+tAH56/tAwmD41eMlL792pLLu24yHt4TnH+elfunDUufKMO/7v5SkfhHE6tnGJX95fjCBN+ztr2i+GfjJ4e1zxFqtrpun2wu/OurqURxR7oGVdzHgZJAHvS4lw9XFZVVo0IuUny2SV3o+wcMYmjhM2pV8RNRglK7bsldaavufS37RPxe+Fvif4L+KNC8PfELw/qWo3VqiwWtrfxySysJUJCqDk8An8K+C4aybMcJmtGtXoTjFPVuLSWj6n33FGd5Zi8nxFChiISnKLSSkm36JO58SOQXYj1NfrK2PySWsmQm3tySTbwknkkxqSf0q+eS6v7zJ0aTd3Ffcv8h4AUBVAAHAAGAKlu5aSSsj6V/Y38feCPA0Pi8eMfFmlaKb2aza2F9dJD5wVHDFdx5wSM/UV8Hxtl2LzB0PqtKU7KV+VN2u1vY+/wCBsyweXRxCxdWMOaUbc0kr2jbS7GftjePfBPjj/hEj4O8V6VrX2M3v2j7DdJN5W4Rbd20nGcHH0NHBWXYvAe3+tUpQvy25k1e172uHG+ZYPMVQ+qVYz5XK/K07XStex4B4Z0ibX/Euj6FbxNLJqOo2tqEUZLB5lDf+O7ifYGvtMXWWGw9Ss/sxk/uTt+Nj4nB0HicTSor7U4r75K/4Xb8kz9OQAw2qPNEnAH3fOC/+g7f1r+d3qf0aH2z/AKjI/wDAf/61AH//1f07/wBaP+e3nf8AAfPx/wCg7f1oA+Gv2t9Lew+M11fklo9X020u0k24EhUPE5x7bFH5V+ycFVlVymML6wlJffZr9T8Y43o+yzhz/nhF/c5J/p+B4yCQcg4NfWHyKdth26R/lyze3WlZIq8paDaZIUAFAChmX7rEZ9DRYabWxa0zTNV13UodG0ewutQ1Cc4itLaNpZmz/sjkD3OB6msqtWlh6bq1ZKMVu3ovv/Ra+RpSpVcTUVGknKb2S1f3dF5uy7s+vP2dv2cJ/BVzB458dQxS64yN/Z+nK4ZbPIwzM44M2CRxwoJAJJJP5bxPxTHMYvB4L+H1e3Nbou0euur3dtj9W4W4VeXNY3Gr970jvyX3u+sraaaJXSvdt/Qn+s/6a+b/AMB8/H/oO39a+GPug+2f9Rkf+A//ANagD//W/Tv/AFn/AE287/gPn4/9B2/rQB82ftp+Dn1DQND8fWymRtMnewu5cY3QTEbGI7bZVVf+Bk1+gcBY9U69TBS+2uZesd//ACVv7j8+4/wHtcPSx0fsPlfpK1v/ACZR+TZ8j1+oH5WegfAXTfCmtfFbRdC8aaPa6npmqLcWf2e5+4Z2j3RHPY5RgD1ycd68PiOrisPllStg5OM42d1va9n+av6Hu8NUcJic0p0cbBShJSjZ6rmauvnZSs/PzPqLUv2RPgzqOGsrDWLEvnZ9k1SULJ64WTcFx+tfnNHjbNqWkpRl6xX5qx+k1uB8nq/DGUfScvyvY52b9ibwPKwNn438Sqrk+WCls2/1xmMYx79a748f4xL3qMH/AOBL/wBuPPfh7gr+7XqL/wAA/wDkB1r+xP4EV0afxp4lulbOxVFtF5uOuD5eVx+tKfH+NatGjBf+BP8A9uHHw+wSd51qj8vcX5RudDpf7InwXsniefTdX1bnKrdapKqT+uVQqAB6Hg1wVeNs3qfDKMfSK/N3O+lwNk1P4oSl6zlb5q6T9Nj07wv4K8I+DbMaf4Q8N6dptvL1W0gWE3OP7x6jHuea+dxePxWPlz4qo5vzdz6TCYDC4CHs8LTjBdkkja/1n/TXzf8AgPn4/wDQdv61yHWH+s/6beb/AMB8/H/oO39aAD7Z/wBRkf8AgP8A/WoA/9f9O/8AWf8ATbzv+A+fj/0Hb+tAFDX9C0zxVol94e1q2F7Y6tA9vcRn5ftCEYOP7hXqD6gVvhsTUwdaNei7Si7p+aMMThqWMoyw9ZXjJNNeTPz7+LPwm8RfCfXjYaor3OmXUhGmapsxHeJ2B7LKP4kPUjK5B4/cclzqhnVD2lLSa+KPVf5x7P5PXf8AC88ySvklfkq6wb92X83l/i7rrvHTRcTHJJDIssMjxSRsHR0YqyMpyGBHIIIBBHIIFeu0pKzV0zxk3F3Ts18mrdn0aeqfRn0h4I/bN1nTbGPT/H3hc60yKEkvrGdIZp8dDJG2E3epQgH0FfA5jwJRrVHUwVTkT+zJNpejV3b1XzZ+g5bx9Vo01Tx1Lna+1FpN+sXZX9HZ9lsd3F+2j8MJULXeg+J0ZwN6G1ibzMdMlZMLj2614suA8xT92cH82vzR7UeP8sa9+FRP/Df8m0UNT/bZ8FxhxpfgnxBfsw+Zp5ILZZMdM/MWXHsK3pcAYuX8WtBenM/0t+JjV8QsHG/saM5evLH85X/A9X+DvxDuvir4Hh8XXmkRWM13d3Nu1rHMZBIIZSgO8gYIAz0Ga+Zz3K45PjHhYy5klF3tbdJ7an1GR5o84wSxco8rbkrXv8Mmt9Ox23+s/wCm3m/8B8/H/oO39a8c9cP9Z/0283/gPn4/9B2/rQAf6z/pt5v/AAHz8f8AoO39aAD7Z/1GR/4D/wD1qAP/0P07m/d/aPO/eeVs83t5ufu/7uPbrQATfu/tHm/vPK2eb283P3f93Ht1oAp65o+k63p99pPiDTbbU7J1Vbq3uIg8c4PTKngY9q2oV6uGqKrRk4yWzTszKvQpYmm6VaKlF7pq6fyPFPFX7HXwv1B7iXQ77WNBkiIJS2nE8J3dAI5Q2wD0BFfXYTjnMqKtWUanqrP742v87nx+L4Eyys70HKn6O6+6SkkvJWOE1X9h/wAQ28lwdK+ImnTxQ7dv2rTXRyD6lJMfpXtUvEGg1+9w7T8pL9V+p4lTw8rpv2eITXS8Hf52kl+CKI/Yk8dLJMs3jfw8qw7clbS4YnPoCwrWXH+DSuqM/vj/AJGUfD3Ft614f+Ay/wDkzp9G/Yl0eykkl8XeO769SAoXg061S2Vw3be5dh9Qfyrz8Tx/VkmsNQS85Ny/BWR6GG8PaMdcVXcvKKUV+PNL7me9+DvBHhv4ceHW8LeGbJ4dPs3MsiPM8jzSSNuLM7EnOTmvicfmGIzOu8RiXeTstktFsrI+4y/L8PleHWGwqtFXerbd27t3bb1ZvTfu/tHm/vPK2eb283P3f93Ht1riO0Jv3f2jzf3nlbPN7ebn7v8Au49utABN+7+0eb+88rZ5vbzc/d/3ce3WgC59mv8A/oI/+QVoA//Z`
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
	usr, _, err := twitterClient.Users.Show(&twitter.UserShowParams{
		ScreenName: "kiwiidb",
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
