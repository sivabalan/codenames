package codenames

import (
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"strconv"
)

const tpl = `
<!DOCTYPE html>
<html>
    <head>
        <title>Codenames - Play Online</title>
        <script src="/static/app.js?v=0.01" type="text/javascript"></script>
        <link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet">
        <link rel="stylesheet" type="text/css" href="/static/game.css" />
        <link rel="stylesheet" type="text/css" href="/static/lobby.css" />
        <link rel="shortcut icon" type="image/png" id="favicon" href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAA8SURBVHgB7dHBDQAgCAPA1oVkBWdzPR84kW4AD0LCg36bXJqUcLL2eVY/EEwDFQBeEfPnqUpkLmigAvABK38Grs5TfaMAAAAASUVORK5CYII="/>

        <script type="text/javascript">
             {{if .SelectedGameID}}
             window.selectedGameID = "{{.SelectedGameID}}";
             {{end}}
             window.autogeneratedGameID = "{{.AutogeneratedGameID}}";
        </script>
    </head>
    <body>
		<script>
		  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
		  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
		  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
		  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

		  ga('create', 'UA-88084599-2', 'auto');
		  ga('send', 'pageview');

		</script>
		<div id="app">
		</div>
    </body>
</html>
`

type templateParameters struct {
	SelectedGameID      string
	AutogeneratedGameID string
}

func (s *Server) handleIndex(rw http.ResponseWriter, req *http.Request) {
	dir, id := filepath.Split(req.URL.Path)
	if dir != "" && dir != "/" {
		http.NotFound(rw, req)
		return
	}

	playerID := getPlayerID(req)

	if playerID == "" {
		playerID := strconv.Itoa(rand.Int())
		rw.Header()["Set-Cookie"] = []string{"player_id=" + playerID}
	}

	autogeneratedID := s.getAutogeneratedID()

	err := s.tpl.Execute(rw, templateParameters{
		SelectedGameID:      id,
		AutogeneratedGameID: autogeneratedID,
	})
	if err != nil {
		http.Error(rw, "error rendering", http.StatusInternalServerError)
	}
}

func (s *Server) getAutogeneratedID() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	autogeneratedID := ""
	for {
		a := strings.ToLower(s.gameIDWords[rand.Intn(len(s.gameIDWords))])
		b := strings.ToLower(s.gameIDWords[rand.Intn(len(s.gameIDWords))])
		autogeneratedID = fmt.Sprintf("%s-%s", a, b)
		if _, ok := s.games[autogeneratedID]; !ok {
			break
		}
	}
	return autogeneratedID
}
