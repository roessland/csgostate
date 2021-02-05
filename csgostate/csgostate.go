package csgostate

import (
	_ "embed"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

//go:embed "gamestate_integration_csgostate.cfg"
var cfgData []byte
var cfgName = "gamestate_integration_csgostate.cfg"

type Listener struct {
	Updates chan State
}

func NewListener() *Listener {
	listener := &Listener{}

	listener.Updates = make(chan State)

	return listener
}

func (listener *Listener) Listen() {
	http.HandleFunc("/", listener.HandlerFunc)
	log.Print("listening on http://127.0.0.1:3528/")
	log.Fatal(http.ListenAndServe("127.0.0.1:3528", nil))
}

func InstallCfg() error {
	cfgDir, err := FindCfgDir()
	if err != nil {
		return errors.New("csgostate: couldn't autodetect csgo/cfg directory")
	}
	WriteCfg(cfgDir)
	return nil
}

func FindCfgDir() (string, error) {
	return `C:\Program Files (x86)\Steam\steamapps\common\Counter-Strike Global Offensive\csgo\cfg`, nil
}

func WriteCfg(cfgDir string) {
	if len(cfgData) == 0 {
		log.Fatal("csgostate: embedding failed during build, cfgData has length 0")
	}
	cfgPath := path.Join(cfgDir, cfgName)
	f, err := os.Create(cfgPath)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(cfgData)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
	log.Printf("wrote %s", cfgPath)
}

func (listener *Listener) getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte("<html><body>CS:GO State HTTP server</body></html>"))
	return
}

func (listener *Listener) postHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("csgostate: reading body: %s", err)
		return
	}

	var change State
	change.RawJson = body
	err = json.Unmarshal(body, &change)
	if err != nil {
		log.Printf("csgostate: unmarshaling body: %s", err)
		return
	}

	// Discard spectator player updates
	if change.Provider.SteamID == change.Player.SteamID {
		listener.Updates <- change
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
}

func (listener *Listener) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		listener.postHandler(w, r)
	}
	listener.getHandler(w, r)
}
