package data_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/regx"
	"github.com/britannic/blacklist/utils"
)

func TestExclusions(t *testing.T) {
	b, err := config.Get(config.Testdata2, global.Area.Root)
	if err != nil {
		t.Error("Couldn't load config.Testdata")
	}

	globex := data.GetExcludes(*b)
	ex := data.GetExcludes(*b)
	dex := make(config.Dict)

	for _, s := range src {
		f := fmt.Sprintf("../testdata/tdata.%v.%v", s.Type, s.Name)
		testdata, err := utils.Getfile(f)
		if err != nil {
			t.Errorf("Cannot open %v", f)
		}

		var tdata string
		for _, l := range testdata {
			if len(l) > 0 {
				tdata += l + "\n"
			}
		}

		gdata := data.Process(s, globex, dex, tdata)

		for k := range gdata.List {
			i := strings.Count(k, ".")
			if i < 1 {
				t.Errorf("key: %v has . count of %v", k, i)
			}

			switch {
			case i == 1:
				if ex.KeyExists(k) {
					t.Errorf("Exclusion failure, found matching key: %v", k)
				}
			case i > 1:
				if ex.SubKeyExists(k) {
					t.Errorf("Exclusion failure, found submatch for key: %v", k)
				}
			}
		}
	}
}

func TestGetHTTP(t *testing.T) {
	type tdata struct {
		body  []byte
		err   error
		prcsd *config.Src
	}

	h := &tdata{}
	d := []*tdata{}
	rx := regx.Regex

	b, err := config.Get(config.Testdata, global.Area.Root)
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

	a := data.GetURLs(*b)
	ex := make(config.Dict)
	dex := make(config.Dict)
	for k := range a {
		for _, u := range a[k] {
			if len(u.URL) > 0 {
				h.body, h.err = data.GetHTTP(u.URL)
				d = append(d, h)
				h.prcsd = data.Process(u, ex, dex, string(h.body[:]))
			}
		}
	}

	for _, z := range d {
		for got := range z.prcsd.List {
			want := rx.FQDN.FindStringSubmatch(got)[1]
			if strings.Compare(got, want) != 0 {
				t.Errorf("wanted: %v - got: %v", want, got)
			}
		}
	}
}

func TestGetUrls(t *testing.T) {
	blist, err := config.Get(config.Testdata, global.Area.Root)
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

	b := *blist
	a := data.GetURLs(b)

	for k := range b {
		for _, url := range a[k] {
			if g, ok := b[k].Source[url.Name]; ok {
				want := g.URL
				got := url.URL
				if want != url.URL {
					t.Errorf("%v URL mismatch:", url.Name)
					fmt.Printf("Wanted %v\nGot: %v", want, got)
				}
			}
		}
	}
}

func TestProcess(t *testing.T) {
	for _, s := range src {
		ex := make(config.Dict)
		dex := make(config.Dict)
		f := fmt.Sprintf("../testdata/tdata.%v.%v", s.Type, s.Name)
		testdata, err := utils.Getfile(f)
		if err != nil {
			t.Errorf("Cannot open %v", f)
		}

		var tdata string
		for _, l := range testdata {
			if len(l) > 0 {
				tdata += l + "\n"
			}
		}

		f = fmt.Sprintf("../testdata/sdata.%v.%v", s.Type, s.Name)
		staticdata, err := utils.Getfile(f)
		if err != nil {
			t.Errorf("Cannot open %v", f)
		}

		var wdata string
		for _, l := range staticdata {
			if len(l) > 0 {
				wdata += l + "\n"
			}
		}

		gdata := string(data.GetList(data.Process(s, ex, dex, tdata))[:])

		if !utils.CmpHash([]byte(wdata), []byte(gdata)) {
			mismatch := []*struct {
				d string
				f string
			}{
				{
					d: wdata,
					f: fmt.Sprintf("/tmp/want.%v.%v", s.Type, s.Name),
				},
				{
					d: gdata,
					f: fmt.Sprintf("/tmp/got.%v.%v", s.Type, s.Name),
				},
			}

			for _, m := range mismatch {
				utils.WriteFile(m.f, []byte(m.d))
			}
			t.Errorf("data mismatch between standard and data.Processed data for %q.", s.Name)
		}
	}
}

func TestPurgeFiles(t *testing.T) {
	whatOS := runtime.GOOS
	if whatOS == "darwin" {
		global.DmsqDir = "/tmp"
		global.Logfile = "/tmp/blacklist.log"
	}

	b, err := config.Get(config.Testdata, global.Area.Root)
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

	urls := data.GetURLs(*b)
	if err := data.PurgeFiles(urls); err != nil {
		t.Errorf("Error removing unused conf files: %v", err)
	}
}

func TestStripPrefix(t *testing.T) {
	rx := regx.Regex()
	tline := `[This line should be delimited by "[]" only.]`

	for _, s := range src {
		var l string
		switch s.Prfx {
		case "http":
			l = s.Prfx + "://" + tline
		default:
			l = s.Prfx + tline
		}

		r, ok := data.StripPrefix(l, s.Prfx, rx)
		switch {
		case tline != r:
			t.Errorf("stripPrefix() failed for %v", s.Name)
			fmt.Printf("Want: %v\nGot: %v\n", tline, r)
		case !ok:
			t.Errorf("stripPrefix() failed for %v", s.Name)
		}
	}
}

// http://play.golang.org/p/KAwluDqGIl
var src = []*config.Src{
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "pre-configured",
		Type:    "domains",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "malc0de",
		Prfx:    "zone ",
		Type:    "domains",
		URL:     "http://malc0de.com/bl/ZONES",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "pre-configured",
		Type:    "hosts",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "adaway",
		Prfx:    "127.0.0.1 ",
		Type:    "hosts",
		URL:     "http://adaway.org/hosts.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "malwaredomainlist",
		Prfx:    "127.0.0.1 ",
		Type:    "hosts",
		URL:     "http://www.malwaredomainlist.com/hostslist/hosts.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "openphish",
		Prfx:    "http",
		Type:    "hosts",
		URL:     "https://openphish.com/feed.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "someonewhocares",
		Prfx:    "0.0.0.0",
		Type:    "hosts",
		URL:     "http://someonewhocares.org/hosts/zero/",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "volkerschatz",
		Prfx:    "http",
		Type:    "hosts",
		URL:     "http://www.volkerschatz.com/net/adpaths",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "winhelp2002",
		Prfx:    "0.0.0.0 ",
		Type:    "hosts",
		URL:     "http://winhelp2002.mvps.org/hosts.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "yoyo",
		Type:    "hosts",
		URL:     "http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext",
	},
}
