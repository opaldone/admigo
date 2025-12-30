package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"admigo/common"
	"admigo/config"

	"github.com/julienschmidt/httprouter"
)

type LookupResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
	ShortName   string `json:"short_name"`
	MiddleName  string `json:"middle_name"`
}

func MapLoca(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	e := config.Env(false)

	info := map[string]string{
		"city": e.Map.City,
	}

	setFrontContent(w, r, "loca", info, nil,
		"map/ix/loca",
		"map/ix/_mus",
		"map/ix/_logs",
		"map/ix/_ros",
	)
}

func MapWs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	e := config.Env(false)

	link := fmt.Sprintf("%s/ws/0", e.Map.Ws)

	js := map[string]string{
		"link":       link,
		"startpoint": e.Map.StartPoint,
		"routeurl":   e.Map.RouteURL,
		"routekey":   e.Map.RouteKey,
	}

	output, _ := json.Marshal(js)

	w.Write(output)
}

func prepLookupResult(list []*LookupResult) {
	for _, li := range list {
		ar := strings.Split(li.DisplayName, ",")
		li.ShortName = fmt.Sprintf("%s, %s", ar[0], ar[1])
		li.MiddleName = fmt.Sprintf("%s,%s,%s,%s", ar[0], ar[1], ar[2], ar[3])
	}
}

func MapShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	e := config.Env(false)

	qv := r.URL.Query()

	ci := ps.ByName("ci")
	qfi := qv.Get("fi")

	pv := fmt.Sprintf("%s,%s", qfi, ci)

	url_in := url.URL{
		Scheme: "https",
		Host:   "nominatim.openstreetmap.org",
		Path:   "/search",
	}

	q := url_in.Query()
	q.Set("q", pv)
	q.Set("format", "json")

	url_in.RawQuery = q.Encode()

	rq, err := http.NewRequest("GET", url_in.String(), nil)
	if err != nil {
		APIError(w, err)
		return
	}

	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Accept-Language", e.Map.Lang)
	rq.Header.Set("User-Agent", "Chrome")

	re, err := http.DefaultClient.Do(rq)
	if err != nil {
		APIError(w, err)
		return
	}

	defer re.Body.Close()

	if re.StatusCode != http.StatusOK {
		APIError(w, fmt.Errorf("External API error %d", re.StatusCode))
		return
	}

	var stru []*LookupResult
	err = json.NewDecoder(re.Body).Decode(&stru)
	if err != nil {
		APIError(w, errors.New("Error decoding API response"))
		return
	}

	prepLookupResult(stru)

	co := ""
	if len(stru) > 0 {
		co, err = getFrontContentAjax(r, stru, nil, "map/ix/_slist")
		if err != nil {
			APIError(w, err)
			return
		}
	}

	ans := common.AjaxAnswer{
		Cont: co,
	}

	output, _ := json.Marshal(ans)

	w.Write(output)
}
