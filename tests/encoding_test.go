package tests

import (
	"math/rand"
	"testing"

	"github.com/sameepkat/ushort/internal/encoding"
)

func randInt() int {
	num := rand.Intn(900000) + 100000
	return num
}

func TestEnconding(t *testing.T) {
	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://openai.com",
		"https://github.com",
		"https://golang.org",
		"https://stackoverflow.com",
		"https://reddit.com",
		"https://news.ycombinator.com",
		"https://docker.com",
		"https://kubernetes.io",
		"https://netflix.com",
		"https://amazon.com",
		"https://facebook.com",
		"https://twitter.com",
		"https://linkedin.com",
		"https://youtube.com",
		"https://microsoft.com",
		"https://apple.com",
		"https://mozilla.org",
		"https://ubuntu.com",
		"https://archlinux.org",
		"https://debian.org",
		"https://python.org",
		"https://rust-lang.org",
		"https://vuejs.org",
		"https://reactjs.org",
		"https://nextjs.org",
		"https://vitejs.dev",
		"https://tailwindcss.com",
		"https://bootstrap.com",
		"https://cloudflare.com",
		"https://digitalocean.com",
		"https://linode.com",
		"https://heroku.com",
		"https://vercel.com",
		"https://supabase.com",
		"https://planet-scale.com",
		"https://stripe.com",
		"https://paypal.com",
		"https://bitbucket.org",
		"https://gitlab.com",
		"https://medium.com",
		"https://dev.to",
		"https://npmjs.com",
		"https://pypi.org",
		"https://rubygems.org",
		"https://crates.io",
		"https://godoc.org",
		"https://pkg.go.dev",
	}

	for i, _ := range urls {
		go func() {
			code := encoding.Encode(uint64(randInt() * i))
			t.Logf("%v\n", code)
		}()
	}
}
