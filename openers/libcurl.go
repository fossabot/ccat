package openers

import (
	"ccat/log"
	"ccat/utils"
	"io"
	"regexp"
	"strings"
	"time"

	curl "github.com/andelf/go-curl"
)

var curlOpenerName = "curl"
var curlOpenerDescription = "get URL via libcurl bindings\n           " +
	curl.Version() + "\n           protocols: " +
	strings.Join(curl.VersionInfo(0).Protocols, ",")

type curlOpener struct {
	easy              *curl.CURL
	name, description string
}

func init() {
	register(&curlOpener{
		name:        curlOpenerName,
		description: curlOpenerDescription,
	})
}

func (f *curlOpener) easyHandlerInit() {
	// we don't cleanup curl stuff when ending because we don't care (we only use one)

	//curl.GlobalInit(curl.GLOBAL_DEFAULT)
	//defer curl.GlobalCleanup()
	f.easy = curl.EasyInit()
	f.easy.Setopt(curl.OPT_VERBOSE, false)
	f.easy.Setopt(curl.OPT_FOLLOWLOCATION, true)
	f.easy.Setopt(curl.OPT_MAXREDIRS, 10)
	f.easy.Setopt(curl.OPT_CONNECTTIMEOUT, 10)
	f.easy.Setopt(curl.OPT_WRITEFUNCTION, func(ptr []byte, userdata interface{}) bool {
		pipe := userdata.(*io.PipeWriter)
		if _, err := pipe.Write(ptr); err != nil {
			return false
		}
		return true
	})

	step := time.Now().Unix()
	dlstep := 0.0
	f.easy.Setopt(curl.OPT_NOPROGRESS, false)
	f.easy.Setopt(curl.OPT_PROGRESSFUNCTION, func(dltotal, dlnow, ultotal, ulnow float64, _ interface{}) bool {
		if time.Now().Unix()-step > 2 {
			log.Debugf("downloaded: %3.2f%%, speed: %.1fKiB/s \r", dlnow/dltotal*100, (dlnow-dlstep)/1000/float64((time.Now().Unix()-step)))
			step = time.Now().Unix()
			dlstep = dlnow
		}
		return true
	})
}

func (f curlOpener) Name() string {
	return f.name
}
func (f curlOpener) Description() string {
	return f.description
}
func (f *curlOpener) Open(s string, _ bool) (io.ReadCloser, error) {

	r, w := io.Pipe()
	go func() {
		log.Debugln(" curl goroutine started")

		if f.easy == nil {
			f.easyHandlerInit()
		}
		//defer easy.Cleanup()

		s = tryTransformUrl(s)

		f.easy.Setopt(curl.OPT_URL, s)

		f.easy.Setopt(curl.OPT_WRITEDATA, w)

		if err := f.easy.Perform(); err != nil {
			println(" curl ERROR", err.Error())
			w.CloseWithError(err)
		}
		w.Close()
		log.Debugln(" curl goroutine ended")
	}()

	return r, nil
}

func (f curlOpener) Evaluate(s string) float32 {
	// https://everything.curl.dev/protocols/curl
	// The latest curl (as of this writing) supports these protocols:
	// DICT, FILE, FTP, FTPS, GOPHER, GOPHERS, HTTP, HTTPS, IMAP, IMAPS, LDAP, LDAPS,
	// MQTT, POP3, POP3S, RTMP, RTSP, SCP, SFTP, SMB, SMBS, SMTP, SMTPS, TELNET, TFTP
	arr := strings.SplitN(s, "://", 2)
	before := arr[0]
	//log.Printf("before=%s found=%v s=%v", before, found, s)
	if utils.StringInSlice(before, curl.VersionInfo(0).Protocols) {
		return 0.1
	}
	return 0
}

func tryTransformUrl(s string) string {
	// ease life by checking urls

	r := regexp.MustCompile(`^https://github.com/(.+)/blob(/.+)$`)
	matches := r.FindStringSubmatch(s)

	if len(matches) > 0 {
		url := "https://raw.githubusercontent.com/" + matches[1] + matches[2]
		log.Debugf("%s looks like a github url, transforming it to get the raw version: %s", s, url)
		return url
	}
	return s
}
