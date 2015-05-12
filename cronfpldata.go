package cronfpldata

import (
    "html/template"
    "net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"log"
	"google.golang.org/appengine/datastore"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"sync"
	"runtime"
	"golang.org/x/net/context"
	"sort"
)

type FplData struct {
	FplTeamData							[]FplTeamData
	FplPlayerTransferInData				[]FplPlayerData
	FplPlayerTransferOutData			[]FplPlayerData
}

type FplTeamData struct {
	TeamName							string		
	Current_fixture						string
	Next_fixture						string
}

type Event_explain struct {
	Name							string		
	NumberOfMins					int
	NumberOfHalf					int
}


type FplPlayerData struct {
		Id								int			`json:"id"`
		Photo							string		`json:"photo"`
		Web_name						string		`json:"web_name"`
		//Event_explain					[]Event_explain			`json:"event_explain"`
		//Fixture_history					[][]interface{} `json:"fixture_history"`
		//Season_history					[][]interface{} `json:"season_history"`
		//Fixtures						[][]interface{} `json:"fixtures"`
		Event_total						int			`json:"event_total"`
		TypeName						string		`json:"type_name"`
		TeamName						string		`json:"team_name"`
		Selected_by						string		`json:"selected_by"`
		Total_points					int			`json:"total_points"`
		Current_fixture					string		`json:"current_fixture"`
		Next_fixture					string		`json:"next_fixture"`
		Team_code						int			`json:"team_code"`
		News							string		`json:"news"`
		Team_id							int			`json:"team_id"`
		Status							string		`json:"status"`
		Code							int			`json:"code"`
		First_name						string		`json:"first_name"`
		Second_name						string		`json:"second_name"`
		Now_cost						int			`json:"now_cost"`
		Chance_of_playing_this_round	int			`json:"chance_of_playing_this_round"`
		Chance_of_playing_next_round	int			`json:"chance_of_playing_next_round"`
		Value_form						string		`json:"value_form"`
		Value_season					string		`json:"value_season"`
		Cost_change_start				int			`json:"cost_change_start"`
		Cost_change_event				int			`json:"cost_change_event"`
		Cost_change_start_fall			int			`json:"cost_change_start_fall"`
		Cost_change_event_fall			int			`json:"cost_change_event_fall"`
		In_dreamteam					bool			`json:"in_dreamteam"`
		Dreamteam_count					int			`json:"dreamteam_count"`
		Selected_by_percent				string		`json:"selected_by_percent"`
		Form							string		`json:"form"`
		Transfers_out					int			`json:"transfers_out"`
		Transfers_in					int			`json:"transfers_in"`
		Transfers_out_event				int			`json:"transfers_out_event"`
		Transfers_in_event				int			`json:"transfers_in_event"`
		Loans_in						int			`json:"loans_in"`
		Loans_out						int			`json:"loans_out"`
		Loaned_in						int			`json:"loaned_in"`
		Loaned_out						int			`json:"loaned_out"`
		Event_points					int			`json:"event_points"`
		Points_per_game					string		`json:"points_per_game"`
		Ep_this							string		`json:"ep_this"`
		Ep_next							string		`json:"ep_next"`
		Special							bool		`json:"special"`
		Minutes							int			`json:"minutes"`
		Goals_scored					int			`json:"goals_scored"`
		Assists							int			`json:"assists"`
		Clean_sheets					int			`json:"clean_sheets"`
		Goals_conceded					int			`json:"goals_conceded"`
		Own_goals						int			`json:"own_goals"`
		Penalties_saved					int			`json:"penalties_saved"`
		Penalties_missed				int			`json:"penalties_missed"`
		Yellow_cards					int			`json:"yellow_cards"`
		Red_cards						int			`json:"red_cards"`
		Saves							int			`json:"saves"`
		Bonus							int			`json:"bonus"`
		Ea_index						int			`json:"ea_index"`
		Bps								int			`json:"bps"`
		Element_type					int			`json:"element_type"`
		Team							int			`json:"team"`
}

type By func(p1, p2 *FplPlayerData) bool

type fplPlayerDataSorter struct {
	fplPlayerData []FplPlayerData
	by      func(p1, p2 *FplPlayerData) bool 
}

func (by By) Sort(fplPlayerData []FplPlayerData) {
	fs := &fplPlayerDataSorter{
		fplPlayerData: fplPlayerData,
		by:      by, 
	}
	sort.Sort(fs)
}

func (s *fplPlayerDataSorter) Len() int {
	return len(s.fplPlayerData)
}

func (s *fplPlayerDataSorter) Swap(i, j int) {
	s.fplPlayerData[i], s.fplPlayerData[j] = s.fplPlayerData[j], s.fplPlayerData[i]
}

func (s *fplPlayerDataSorter) Less(i, j int) bool {
	return s.by(&s.fplPlayerData[i], &s.fplPlayerData[j])
}

func init() {

	http.HandleFunc("/cronfpldata", cronfpldata)
	http.HandleFunc("/cronfplTeamdata", cronfplTeamdata)
	http.HandleFunc("/retrievefpldata", retrievefpldata)
	http.HandleFunc("/retrieveFplDataByTeam", retrieveFplDataByTeam)
	http.HandleFunc("/retrieveFplDataByTrend", retrieveFplDataByTrend)
	http.HandleFunc("/retrieveFplDataByAvailability", retrieveFplDataByAvailability)
    http.HandleFunc("/", root)
	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/js/", http.FileServer(http.Dir(".")))
}


func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
       
        
		q := datastore.NewQuery("FplTeamData").Order("TeamName")

        fplDatas := make([]FplTeamData, 0, 20)
        if _, err := q.GetAll(c, &fplDatas); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
		
        
		fplStatsForm.ExecuteTemplate(w, "fantasyPremierLeague.htm", fplDatas)
}

func retrieveFplDataByTeam(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
       
		q := datastore.NewQuery("FplPlayerData").Limit(20)
        fplTeamDatas := make([]FplPlayerData, 0, 20)
        if _, err := q.GetAll(c, &fplTeamDatas); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        
		fplStatsTeamForm.ExecuteTemplate(w, "fantasyPremierLeagueTeam.htm", fplTeamDatas)
		
}

func retrieveFplDataByAvailability(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
       
		q := datastore.NewQuery("FplPlayerData").Filter("News >", "")

        fplDatas := make([]FplPlayerData, 0, 600)
        if _, err := q.GetAll(c, &fplDatas); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
		
		
		TeamName := func(p1, p2 *FplPlayerData) bool {
			return p1.TeamName < p2.TeamName
		}
		
        By(TeamName).Sort(fplDatas)
		
		fplStatsAvailabilityForm.ExecuteTemplate(w, "fantasyPremierLeagueAvailability.htm", fplDatas)
		
}

func retrieveFplDataByTrend(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
		runtime.GOMAXPROCS(2)
		var wg sync.WaitGroup
		wg.Add(2)
        
		fplPlayerTransferInDatas := make([]FplPlayerData, 0, 10)
		fplPlayerTransferOutDatas := make([]FplPlayerData, 0, 10)
		
		   
		go func(){
			defer wg.Done()
			fplPlayerTransferInDatas = queryFplDataOrderBy(w,c, "-Transfers_in_event")
		}()
		
		go func(){
			defer wg.Done()
			fplPlayerTransferOutDatas = queryFplDataOrderBy(w,c, "-Transfers_out_event")
		}()
		wg.Wait()
		
		fplData := &FplData{
		  //FplTeamData: fplTeamDatas,
		  FplPlayerTransferInData: fplPlayerTransferInDatas,
		  FplPlayerTransferOutData: fplPlayerTransferOutDatas,
		}
		
		fplStatsTrendForm.ExecuteTemplate(w, "fantasyPremierLeagueTrend.htm", fplData)
		
}


func queryFplDataOrderBy(w http.ResponseWriter, c context.Context, orderBy string) []FplPlayerData {
	
	fplPlayerData := make([]FplPlayerData, 0, 10)
	q := datastore.NewQuery("FplPlayerData").Order(orderBy).Limit(10)
    if _, err := q.GetAll(c, &fplPlayerData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return nil
    }
	return fplPlayerData
}

func retrievefpldata(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
       
		q := datastore.NewQuery("FplPlayerData").Filter("TeamName =", "Arsenal").Limit(10)
        fplDatas := make([]FplPlayerData, 0, 10)
        if _, err := q.GetAll(c, &fplDatas); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
		
        
		fplStatsForm.ExecuteTemplate(w, "fantasyPremierLeague.htm", fplDatas)
		
}

func fplDataKey(c context.Context) *datastore.Key {
        return datastore.NewKey(c, "FplPlayerData", "FplPlayerData", 0, nil)
}

func fplTeamDataKey(c context.Context) *datastore.Key {
        return datastore.NewKey(c, "FplTeamData", "FplTeamData", 0, nil)
}

func clearCronFplPlayerData(c context.Context){

	fplDatas := make([]FplPlayerData, 0, 1000)
    q, err := datastore.NewQuery("FplPlayerData").Ancestor(fplDataKey(c)).Limit(1000).GetAll(c, &fplDatas)
	if err != nil {
		log.Fatal(err)
	}
	
	datastore.DeleteMulti(c, q)

}

func clearCronFplTeamData(c context.Context){

	fplTeamData := make([]FplTeamData, 0, 1000)
    q, err := datastore.NewQuery("FplTeamData").Ancestor(fplTeamDataKey(c)).Limit(1000).GetAll(c, &fplTeamData)
	if err != nil {
		log.Fatal(err)
	}
	
	datastore.DeleteMulti(c, q)

}

func insertCronFplPlayerData(w http.ResponseWriter, c context.Context){
	
	
	client := urlfetch.Client(c)
	
	for i := 1; i < 700; i=i+1 {
		go func (i int) {
           insertCronFplPlayerDataIndividually(w, c, client, i)
        } (i);
			
	}
	
	log.Println("Finish Running insertCronFplPlayerData")
}

func insertCronFplPlayerDataIndividually(w http.ResponseWriter, c context.Context, client *http.Client, i int){
		
		stringNumber := strconv.Itoa(i)
		str := "http://fantasy.premierleague.com/web/api/elements/" + stringNumber
		resp, err := client.Get(str)
		
		if err != nil {
			log.Fatal(err)
		}
		
		robots, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		
		res := &FplPlayerData{}
		json.Unmarshal(robots, &res)
		
		//Returns if unable to parse
		if res.Id == 0 {
			return
		}
		
		key := datastore.NewIncompleteKey(c, "FplPlayerData", fplDataKey(c))
		if _, err := datastore.Put(c, key, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
}

func insertCronFplTeamData(w http.ResponseWriter, c context.Context){
	
	q := datastore.NewQuery("FplPlayerData").Project("TeamName", "Current_fixture", "Next_fixture").Distinct()

    fplPlayerDatas := make([]FplPlayerData, 0, 20)
    if _, err := q.GetAll(c, &fplPlayerDatas); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
    }
	
	//insert
	res := &FplTeamData{}
	
	for _,element := range fplPlayerDatas {
		res.TeamName = element.TeamName
		res.Current_fixture = element.Current_fixture
		res.Next_fixture = element.Next_fixture
		
		key := datastore.NewIncompleteKey(c, "FplTeamData", fplTeamDataKey(c))
		if _, err := datastore.Put(c, key, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	
	
}

func cronfpldata(w http.ResponseWriter, r *http.Request) {
	
	c := appengine.NewContext(r)
	
	clearCronFplPlayerData(c)
	insertCronFplPlayerData(w, c)
	

}

func cronfplTeamdata(w http.ResponseWriter, r *http.Request) {
	
	c := appengine.NewContext(r)
	clearCronFplTeamData(c)
	insertCronFplTeamData(w, c)
	
}


var fplStatsForm = template.Must(template.New("").ParseFiles("fantasyPremierLeague.htm"))
var fplStatsTeamForm = template.Must(template.New("").ParseFiles("fantasyPremierLeagueTeam.htm"))
var fplStatsTrendForm = template.Must(template.New("").ParseFiles("fantasyPremierLeagueTrend.htm"))
var fplStatsAvailabilityForm = template.Must(template.New("").ParseFiles("fantasyPremierLeagueAvailability.htm"))