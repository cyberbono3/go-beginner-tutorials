package main

/*
 * [Turn on Tor, for use proxy]:
 *     $ sudo systemctl start tor.service
 * [Example]:
 *     $ go build main.go 
 *     $ ./main_parser mail.ru -t a -a href -tp -ua
 *     $ ./main_parser mainl.ru--tag a --attr href --tor-proxy --user-agent
 */

import (
    "os"
    "fmt"
    "time"
    "regexp"
    "net/url"
    "net/http"
)

/* $ go get golang.org/x/net/proxy */
import "golang.org/x/net/proxy"

const (
    BUFF = 512
    TORS_PROXY = "socks5://127.0.0.1:9050"
    USER_AGENT = "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1547.2 Safari/537.36"
)

var (
    TAG_NAME = ""
    ATR_NAME = ""

    GET_TORS_PROXY = false
    GET_USER_AGENT = false
)

func main() {
    check_args(os.Args)
    check_func(urlopen("http://" + os.Args[1]))
}

func check_func (html string) {
    if TAG_NAME != "" {
        for _, result := range get_object(html) {
            fmt.Println(result)
        }
    } else {
        fmt.Println(html)
    }
}

func check_args (args []string) {
    var (
        flag_tag bool = false
        flag_atr bool = false
    )

    switch len(args) {
        case 1: 
            get_error("no url specified")
        case 2:
            if args[1] == "-i" || args[1] == "--info" {
                get_info()
            } else { return }
        default:
            for _, value := range args[2:] {

                switch value {
                    case "-tp", "--tor-proxy":  GET_TORS_PROXY = true; continue
                    case "-ua", "--user-agent": GET_USER_AGENT = true; continue
                    case "-t", "--tag":  flag_tag = true; continue
                    case "-a", "--attr": flag_atr = true; continue
                }

                if flag_tag {
                    TAG_NAME = value
                    flag_tag = false

                } else if flag_atr {
                    ATR_NAME = value
                    flag_atr = false
                }
            }
    }

}

func get_info() {
    fmt.Println(
`Modules:
    -tp || --tor-proxy  -> Turn on tor proxy
    -ua || --user-agent -> Turn on user-agent
    -t  || --tag        -> Tag name
    -a  || --attr       -> Attribute name
Example:
    $ parse google.com --tag a --attr href -tp -ua`)
    os.Exit(0)
}

func urlopen(url_str string) string {
    var (
        html_page string
        buffer []byte
        length int
        err error
    )

    var (
        client *http.Client
        req *http.Request
        resp *http.Response
    )

    if GET_TORS_PROXY {
        torProxyUrl, err := url.Parse(TORS_PROXY)
        check_error(err)

        torDialer, err := proxy.FromURL(torProxyUrl, proxy.Direct)
        check_error(err)

        torTransport := &http.Transport {
            Dial: torDialer.Dial,
        }

        client = &http.Client {
            Transport: torTransport, 
            Timeout: time.Second * 10,
        }
    } else {
        client = &http.Client{}
    }
    
    req, err = http.NewRequest("GET", url_str, nil)
    check_error(err)

    req.Header.Add("Accept", "text/html")

    if GET_USER_AGENT {
        req.Header.Add("User-Agent", USER_AGENT)
    } 

    resp, err = client.Do(req)
    check_error(err)
    defer resp.Body.Close()

    buffer = make([]byte, BUFF)
    for {
        length, err = resp.Body.Read(buffer)
        html_page += string(buffer[:length])
        if length == 0 || err != nil{ break }
    }

    return html_page
}

func get_object (html string) []string {
    var (
        result [][]string
        slice_result []string
        regular *regexp.Regexp 
    )

    if TAG_NAME != "" {
        if ATR_NAME != "" {
            regular = regexp.MustCompile(`<`+TAG_NAME+`.*?`+ATR_NAME+`\s*=\s*['"]([^\s'"]+)[\s'"]`)
            result = regular.FindAllStringSubmatch(html, -1)
            for _, slice := range result {
                slice_result = append(slice_result, slice[1])
            }
            TAG_NAME = ""; ATR_NAME = ""

        } else {
            regular = regexp.MustCompile(`<`+TAG_NAME+`[^>]*>.+</`+TAG_NAME+`>`)
            result = regular.FindAllStringSubmatch(html, -1)
            for _, slice := range result {
                slice_result = append(slice_result, slice[0])
            }
            TAG_NAME = ""
        }
    } 

    return slice_result
}

func check_error (err error) {
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}

func get_error (err string) {
    fmt.Println("Error:", err)
    os.Exit(1)
}
