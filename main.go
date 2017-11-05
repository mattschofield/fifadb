package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var input chan string
var done chan bool
var futheadBaseURL = "https://www.futhead.com"
var players []Player
var unknown = "Unknown"

func main() {
	input = make(chan string)
	done = make(chan bool)

	players = make([]Player, 0)

	go findAllPlayers()

	for i := 0; i < 7; i++ {
		go findEachPlayerData()
	}

	<-done

	csv, err := gocsv.MarshalString(&players)
	if err != nil {
		panic(err)
	}

	fmt.Println(csv)
}

func visit(url string) *html.Node {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	node, err := html.Parse(r.Body)
	if err != nil {
		panic(err)
	}

	return node
}

func findAllPlayers() {
	for index := 1; index <= 209; index++ {
		node := visit(futheadBaseURL + "/18/players/?page=" + strconv.Itoa(index))
		playerMatcher := scrape.ByClass("player-group-item")
		players := scrape.FindAll(node, playerMatcher)

		for _, player := range players {
			if strings.Contains(scrape.Attr(player, "class"), "list-group-column-headers") {
				continue
			} else {
				input <- findPlayerLink(player)
			}
		}
	}
	close(input)

	done <- true
}

func findEachPlayerData() {
	for url := range input {
		node := visit(url)
		pos := findPlayerPosition(node)

		if strings.Contains(pos, "GK") {
			continue
		} else {
			player := Player{
				ID:                      extractPlayerIDFromURL(url),
				FirstName:               findPlayerFirstName(node),
				LastName:                findPlayerLastName(node),
				Position:                findPlayerPosition(node),
				Club:                    findPlayerClub(node),
				League:                  findPlayerLeague(node),
				Nationality:             findPlayerNationality(node),
				Rating:                  findPlayerRating(node),
				LowestPrice:             findPlayerLowestPrice(node),
				ImageURL:                findPlayerImageURL(node),
				FUTHeadTotalStats:       findPlayerFUTHeadTotalStats(node),
				FUTHeadTotalInGameStats: findPlayerFUTHeadTotalInGameStats(node),
				HeightInCM:              findPlayerHeight(node),
				StrongFoot:              findPlayerStrongFoot(node),
				WorkRateOffensive:       findPlayerWorkRateOffensive(node),
				WorkRateDefensive:       findPlayerWorkRateDefensive(node),
				SkillMoves:              findPlayerSkillMoves(node),
				WeakFoot:                findPlayerWeakFoot(node),
				Pace:                    findPlayerStat(node, "Pace"),
				PaceAttributes: PaceAttributes{
					Acceleration: findPlayerStatAttribute(node, "Acceleration"),
					SprintSpeed:  findPlayerStatAttribute(node, "Sprint Speed"),
				},
				Shooting: findPlayerStat(node, "Shooting"),
				ShootingAttributes: ShootingAttributes{
					Positioning: findPlayerStatAttribute(node, "Positioning"),
					Finishing:   findPlayerStatAttribute(node, "Finishing"),
					ShotPower:   findPlayerStatAttribute(node, "Shot Power"),
					LongShots:   findPlayerStatAttribute(node, "Long Shots"),
					Volleys:     findPlayerStatAttribute(node, "Volleys"),
					Penalties:   findPlayerStatAttribute(node, "Penalties"),
				},
				Passing: findPlayerStat(node, "Passing"),
				PassingAttributes: PassingAttributes{
					Vision:       findPlayerStatAttribute(node, "Vision"),
					Crossing:     findPlayerStatAttribute(node, "Crossing"),
					FreeKick:     findPlayerStatAttribute(node, "Free Kick"),
					ShortPassing: findPlayerStatAttribute(node, "Short Passing"),
					LongPassing:  findPlayerStatAttribute(node, "Long Passing"),
					Curve:        findPlayerStatAttribute(node, "Curve"),
				},
				Dribbling: findPlayerStat(node, "Dribbling"),
				DribblingAttributes: DribblingAttributes{
					Agility:     findPlayerStatAttribute(node, "Agility"),
					Balance:     findPlayerStatAttribute(node, "Balance"),
					Reactions:   findPlayerStatAttribute(node, "Reactions"),
					BallControl: findPlayerStatAttribute(node, "Ball Control"),
					Dribbling:   findPlayerStatAttribute(node, "Dribbling"),
					Composure:   findPlayerStatAttribute(node, "Composure"),
				},
				Defending: findPlayerStat(node, "Defending"),
				DefendingAttributes: DefendingAttributes{
					Interceptions:  findPlayerStatAttribute(node, "Interceptions"),
					Heading:        findPlayerStatAttribute(node, "Heading"),
					Marking:        findPlayerStatAttribute(node, "Marking"),
					StandingTackle: findPlayerStatAttribute(node, "Standing Tackle"),
					SlidingTackle:  findPlayerStatAttribute(node, "Sliding Tackle"),
				},
				Physical: findPlayerStat(node, "Physical"),
				PhysicalAttributes: PhysicalAttributes{
					Jumping:    findPlayerStatAttribute(node, "Jumping"),
					Stamina:    findPlayerStatAttribute(node, "Stamina"),
					Strength:   findPlayerStatAttribute(node, "Strength"),
					Aggression: findPlayerStatAttribute(node, "Aggression"),
				},
			}

			players = append(players, player)
		}
	}
}

func findPlayerName(node *html.Node) string {
	m := scrape.ByClass("player-name")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	return scrape.Text(n)
}

func findPlayerLink(node *html.Node) string {
	m := scrape.ByTag(atom.A)
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	return fmt.Sprintf("https://www.futhead.com%s", scrape.Attr(n, "href"))
}

func extractPlayerIDFromURL(url string) int {
	s := strings.Replace(url, "https://www.futhead.com/18/players/", "", 1)
	id, err := strconv.Atoi(strings.Split(s, "/")[0])
	if err != nil {
		panic(err)
	}

	return id
}

func findPlayerFirstName(node *html.Node) string {
	m := scrape.ByClass("firstname")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	return scrape.Text(n)
}

func findPlayerLastName(node *html.Node) string {
	m := scrape.ByClass("playercard-name")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	return scrape.Text(n)
}

func findPlayerPosition(node *html.Node) string {
	m := scrape.ByClass("playercard-position")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	return scrape.Text(n)
}

func findPlayerClub(node *html.Node) string {
	m := matcherForSideBarStat("Club")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	return scrape.Text(n)
}

func findPlayerLeague(node *html.Node) string {
	m := matcherForSideBarStat("League")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	return scrape.Text(n)
}

func findPlayerNationality(node *html.Node) string {
	m := matcherForSideBarStat("Nation")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	return scrape.Text(n)
}

func findPlayerStrongFoot(node *html.Node) string {
	m := matcherForSideBarStat("Strong Foot")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	return scrape.Text(n)
}

func findPlayerHeight(node *html.Node) int {
	m := matcherForSideBarStat("Height")
	n, ok := scrape.Find(node, m)
	if !ok {
		return 0
	}

	height := strings.Split(scrape.Text(n), " | ")[0]
	height = strings.Replace(height, "cm", "", 1)

	return atoi(height)
}

func findPlayerWorkRateOffensive(node *html.Node) string {
	m := matcherForSideBarStat("Workrates")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	wr := strings.Split(scrape.Text(n), " / ")[0]

	return wr
}

func findPlayerWorkRateDefensive(node *html.Node) string {
	m := matcherForSideBarStat("Workrates")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	wr := strings.Split(scrape.Text(n), " / ")[1]

	return wr
}

func matcherForSideBarStat(stat string) scrape.Matcher {
	return func(n *html.Node) bool {
		if n.DataAtom == atom.Div && strings.Contains(scrape.Attr(n, "class"), "player-sidebar-value") {
			return strings.Contains(scrape.Text(n.Parent), stat)
		}
		return false
	}
}

func findPlayerSkillMoves(node *html.Node) int {
	m := scrape.ByClass("playercard-skill-move")
	n, ok := scrape.Find(node, m)
	if !ok {
		return 0
	}

	sm := scrape.Text(n)
	sm = strings.Replace(sm, " SM", "", 1)

	return atoi(sm)
}

func findPlayerWeakFoot(node *html.Node) int {
	m := scrape.ByClass("playercard-weak-foot")
	n, ok := scrape.Find(node, m)
	if !ok {
		return 0
	}

	wf := scrape.Text(n)
	wf = strings.Replace(wf, " WF", "", 1)

	return atoi(wf)
}

func findPlayerRating(node *html.Node) int {
	m := scrape.ByClass("playercard-rating")
	n, ok := scrape.Find(node, m)
	if !ok {
		return 0
	}

	return atoi(scrape.Text(n))
}

func findPlayerLowestPrice(node *html.Node) int {
	m := scrape.ByClass("player-info-price-lowest-bin")
	n, ok := scrape.Find(node, m)
	if !ok {
		return 0
	}

	return atoi(scrape.Text(n))
}

func findPlayerImageURL(node *html.Node) string {
	m := scrape.ByClass("playercard-picture")
	n, ok := scrape.Find(node, m)
	if !ok {
		return unknown
	}

	m = scrape.ByTag(atom.Img)
	n, _ = scrape.Find(n, m)

	return scrape.Attr(n, "src")
}

func findPlayerFUTHeadTotalStats(node *html.Node) int {
	m := func(n *html.Node) bool {
		if n.DataAtom == atom.Div && strings.Contains(scrape.Attr(n.Parent, "class"), "player-stat-group-seven") {
			child, _ := scrape.Find(n, scrape.ByTag(atom.H5))
			return scrape.Text(child) == "Total Stats"
		}
		return false
	}

	n, ok := scrape.Find(node, m)
	if !ok {
		return 0
	}

	stat := scrape.ByClass("standalone-stat")
	value, _ := scrape.Find(n, stat)

	return atoi(scrape.Text(value))
}

func findPlayerFUTHeadTotalInGameStats(node *html.Node) int {
	m := func(n *html.Node) bool {
		if n.DataAtom == atom.Div && strings.Contains(scrape.Attr(n.Parent, "class"), "player-stat-group-seven") {
			child, _ := scrape.Find(n, scrape.ByTag(atom.H5))
			return scrape.Text(child) == "Total IGS"
		}
		return false
	}

	n, ok := scrape.Find(node, m)
	if !ok {
		return 0
	}

	stat := scrape.ByClass("standalone-stat")
	value, ok := scrape.Find(n, stat)
	if !ok {
		return 0
	}

	return atoi(scrape.Text(value))
}

func matcherForStat(stat string) scrape.Matcher {
	return func(n *html.Node) bool {
		if n.DataAtom == atom.Div && strings.Contains(scrape.Attr(n, "class"), "igs-group") {
			return strings.Contains(scrape.Text(n), stat)
		}
		return false
	}
}

func matcherForStatAttribute(stat string) scrape.Matcher {
	return func(n *html.Node) bool {
		// Find the parent Div of a Span which contains the stat name
		if n.DataAtom == atom.Div && strings.Contains(scrape.Attr(n, "class"), "player-stat-row") {
			attr := scrape.ByClass("player-stat-title")
			child, _ := scrape.Find(n, attr)

			return scrape.Text(child) == stat
		}
		return false
	}
}

func findPlayerStat(node *html.Node, stat string) int {
	attr, ok := scrape.Find(node, matcherForStat(stat))
	if !ok {
		// panic(fmt.Errorf("could not find stat %s", stat))
		return 0
	}

	base, ok := scrape.Find(attr, scrape.ByClass("chembot-delta"))
	if !ok {
		// panic(fmt.Errorf("could not find stat value for %s", stat))
		return 0
	}

	value := scrape.Attr(base, "data-chembot-base")

	return atoi(value)
}

func findPlayerStatAttribute(node *html.Node, stat string) int {
	attr, ok := scrape.Find(node, matcherForStatAttribute(stat))
	if !ok {
		// panic(fmt.Errorf("could not find stat attribute %s", stat))
		return 0
	}

	value, ok := scrape.Find(attr, scrape.ByClass("player-stat-value"))
	if !ok {
		// panic(fmt.Errorf("could not find attribute value for %s", stat))
		return 0
	}

	return atoi(scrape.Text(value))
}

func atoi(text string) int {
	n, err := strconv.Atoi(text)
	if err != nil {
		n = 0
	}

	return n
}
