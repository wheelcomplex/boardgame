package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/jkomoros/boardgame"
	"github.com/jkomoros/boardgame/server/api/users"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	managers managerMap
	storage  *ServerStorageManager
	//We store the last error so that next time viewHandler is called we can
	//display it. Yes, this is a hack.
	lastErrorMessage string
	config           *ConfigMode

	gameVersionCacheLock sync.RWMutex
	gameVersionCache     map[string]int
}

type Renderer struct {
	c            *gin.Context
	rendered     bool
	cookieCalled bool
	cookieValue  string
}

type Config struct {
	Dev  *ConfigMode
	Prod *ConfigMode
}

type ConfigMode struct {
	AllowedOrigins    string
	DefaultPort       string
	FirebaseProjectId string
	AdminUserIds      []string
	//This is a dangerous config. Only enable in Dev!
	DisableAdminChecking bool
	StorageConfig        map[string]string
}

type MoveForm struct {
	Name        string
	Description string
	Fields      []*MoveFormField
}

type MoveFormFieldType int

type MoveFormField struct {
	Name         string
	Type         boardgame.PropertyType
	DefaultValue interface{}
}

func (c *ConfigMode) Validate() error {
	if c.DefaultPort == "" {
		return errors.New("No default port provided")
	}
	//AllowedOrigins will just be default allow
	if c.AllowedOrigins == "" {
		log.Println("No AllowedOrigins found. Defaulting to '*'")
		c.AllowedOrigins = "*"
	}
	if c.StorageConfig == nil {
		c.StorageConfig = make(map[string]string)
	}
	return nil
}

const (
	configFileName = "config.SECRET.json"
)

type managerMap map[string]*boardgame.GameManager

/*

Overview of the types of handlers and methods

server.fooHandler take a context. They grab all of the dependencies and pass them to the doers.
server.doFoo takes a renderer and all dependencies that come from context. It may fetch additional items from e.g. storage. It renders the result.
server.getRequestFoo fetches an argument from the context's request and nothing else
server.getFoo grabs a thing that was stored in Context and nothing else
server.setFoo sets a thing into context and nothing else
server.calcFoo takes dependencies and returns a result, with no touching context.
*/

/*

NewServer returns a new server. Get it to run by calling Start(). storage
should a *ServerStorageManager, which can be created either from
NewDefaultStorageManager or NewServerStorageManager.

Use it like so:

	func main() {
		storage := server.NewDefaultStorageManager()
		defer storage.Close()
		server.NewServer(storage, mygame.NewManager(storage)).Start()
	}

*/
func NewServer(storage *ServerStorageManager, managers ...*boardgame.GameManager) *Server {

	result := &Server{
		managers:         make(managerMap),
		storage:          storage,
		gameVersionCache: make(map[string]int),
	}

	storage.server = result

	for _, manager := range managers {
		name := manager.Delegate().Name()
		result.managers[name] = manager
		if manager.Storage() != storage {
			log.Println("The storage for one of the managers was not the same item passed in as major storage.")
			return nil
		}
	}

	return result

}

func NewRenderer(c *gin.Context) *Renderer {
	return &Renderer{
		c,
		false,
		false,
		"",
	}
}

func (r *Renderer) Error(message string) {
	if r.rendered {
		panic("Error called on already-rendered renderer")
	}

	r.writeCookie()

	r.c.JSON(http.StatusOK, gin.H{
		"Status": "Failure",
		"Error":  message,
	})

	r.rendered = true
}

func (r *Renderer) Success(keys gin.H) {

	if r.rendered {
		panic("Success called on alread-rendered renderer")
	}

	r.writeCookie()

	if keys == nil {
		keys = gin.H{}
	}

	result := gin.H{}

	for key, val := range keys {
		result[key] = val
	}

	result["Status"] = "Success"

	r.c.JSON(http.StatusOK, result)

	r.rendered = true
}

func (r *Renderer) writeCookie() {
	if r.rendered {
		return
	}
	if !r.cookieCalled {
		return
	}

	//TODO: might need to set the domain in production.

	if r.cookieValue == "" {
		//Unset the cookie
		r.c.SetCookie(cookieName, "", int(time.Now().Add(time.Hour*10000*-1).Unix()), "", "", false, false)
		return
	}

	r.c.SetCookie(cookieName, r.cookieValue, int(time.Now().Add(time.Hour*100).Unix()), "", "", false, false)
}

//SetAuthCookie will set the auth cookie to the specified value. If called
//multiple times for a single request will only actually write headers for the
//last one.
func (r *Renderer) SetAuthCookie(value string) {

	//We don't write the cookies to the response yet because we might get
	//multiple SetAuthCookie calls in one response.

	r.cookieCalled = true
	r.cookieValue = value

}

func (s *Server) userSetup(c *gin.Context) {
	cookie := s.getRequestCookie(c)

	if cookie == "" {
		log.Println("No cookie set")
		return
	}

	user := s.storage.GetUserByCookie(cookie)

	if user == nil {
		log.Println("No user associated with that cookie")
		return
	} else {
		user.LastSeen = time.Now().UnixNano()
		s.storage.UpdateUser(user)
	}

	s.setUser(c, user)

	s.setAdminAllowed(c, s.calcAdminAllowed(user))
}

func (s *Server) gameFromId(gameId, gameName string) *boardgame.Game {

	manager := s.managers[gameName]

	if manager == nil {
		log.Println("Couldnt' find manager for", gameName)
		return nil
	}

	game := manager.Game(gameId)

	//TODO: figure out a way to return a meaningful error

	if game == nil {
		log.Println("Couldn't find game with id", gameId)
		return nil
	}

	if game.Name() != gameName {
		log.Println("The name of the game was not what we were expecting. Wanted", gameName, "got", game.Name())
		return nil
	}

	return game
}

//gameAPISetup fetches the game configured in the URL and puts it in context.
func (s *Server) gameAPISetup(c *gin.Context) {

	id := s.getRequestGameId(c)

	gameName := s.getRequestGameName(c)

	game := s.gameFromId(id, gameName)

	if game == nil {
		return
	}

	s.setGame(c, game)

	userIds := s.storage.UserIdsForGame(id)

	if userIds == nil {
		log.Println("No userIds associated with game")
	}

	user := s.getUser(c)

	if user == nil {
		log.Println("No user provided")
		//The rest of the flow will handle a nil user fine
	}

	effectiveViewingAsPlayer, emptySlots := s.calcViewingAsPlayerAndEmptySlots(userIds, user, game.Agents())

	if user != nil && effectiveViewingAsPlayer == boardgame.ObserverPlayerIndex && len(emptySlots) > 0 && len(emptySlots) == game.NumPlayers()-game.NumAgentPlayers() {
		//Special case: we're the first player, we likely just created it. Just join the thing!

		slot := emptySlots[0]

		if err := s.storage.SetPlayerForGame(game.Id(), slot, user.Id); err != nil {
			log.Println("Tried to set the user as player " + slot.String() + " but failed: " + err.Error())
			return
		} else {
			s.setHasEmptySlots(c, false)
			effectiveViewingAsPlayer = slot
		}

	} else {
		s.setHasEmptySlots(c, len(emptySlots) != 0)
	}
	s.setViewingAsPlayer(c, effectiveViewingAsPlayer)

}

//Checks to make sure the user is logged in, fails if not.
func (s *Server) requireLoggedIn(c *gin.Context) {

	r := NewRenderer(c)

	user := s.getUser(c)

	if user == nil {
		r.Error("Not logged in")
		c.Abort()
		return
	}

	//All good!
}

func (s *Server) gameStatusHandler(c *gin.Context) {
	//This handler is designed to be a very simple status marker for the
	//current version of the specific game. It will be hit hard by all
	//clients, repeatedly, so it should be very fast.

	id := s.getRequestGameId(c)
	name := s.getRequestGameName(c)

	r := NewRenderer(c)

	s.doGameStatus(r, id, name)

}

func (s *Server) doGameStatus(r *Renderer, gameId, gameName string) {

	s.gameVersionCacheLock.RLock()

	version, ok := s.gameVersionCache[gameId]

	s.gameVersionCacheLock.RUnlock()

	if !ok {
		//Guess for whatever reason the game's version wasn't in the cache. Fetch it.

		game := s.gameFromId(gameId, gameName)

		if game == nil {
			r.Error("Couldn't find game with that name")
			return
		}

		s.gameVersionCacheLock.Lock()
		s.gameVersionCache[game.Id()] = game.Version()
		s.gameVersionCacheLock.Unlock()

		version = game.Version()

	}

	r.Success(gin.H{
		"Version": version,
	})
}

func (s *Server) joinGameHandler(c *gin.Context) {
	r := NewRenderer(c)

	game := s.getGame(c)

	if game == nil {
		r.Error("No such game")
		return
	}

	user := s.getUser(c)

	userIds := s.storage.UserIdsForGame(game.Id())

	viewingAsPlayer, emptySlots := s.calcViewingAsPlayerAndEmptySlots(userIds, user, game.Agents())

	s.doJoinGame(r, game, viewingAsPlayer, emptySlots, user)

}

func (s *Server) doJoinGame(r *Renderer, game *boardgame.Game, viewingAsPlayer boardgame.PlayerIndex, emptySlots []boardgame.PlayerIndex, user *users.StorageRecord) {

	if user == nil {
		r.Error("No user provided.")
		return
	}

	if viewingAsPlayer != boardgame.ObserverPlayerIndex {
		r.Error("The given player is already in the game.")
		return
	}

	if len(emptySlots) == 0 {
		r.Error("There aren't any empty slots in the game to join.")
		return
	}

	slot := emptySlots[0]

	if err := s.storage.SetPlayerForGame(game.Id(), slot, user.Id); err != nil {
		r.Error("Tried to set the user as player " + slot.String() + " but failed: " + err.Error())
		return
	}

	r.Success(nil)

}

func (s *Server) newGameHandler(c *gin.Context) {

	r := NewRenderer(c)

	managerId := s.getRequestManager(c)

	numPlayers := s.getRequestNumPlayers(c)

	manager := s.managers[managerId]

	if numPlayers == 0 && manager != nil {
		numPlayers = manager.Delegate().DefaultNumPlayers()
	}

	agents := s.getRequestAgents(c, numPlayers)

	s.doNewGame(r, manager, numPlayers, agents)

}

func (s *Server) doNewGame(r *Renderer, manager *boardgame.GameManager, numPlayers int, agents []string) {

	if manager == nil {
		r.Error("No manager provided")
		return
	}

	game := boardgame.NewGame(manager)

	if game == nil {
		r.Error("No game could be created")
		return
	}

	if err := game.SetUp(numPlayers, agents); err != nil {
		//TODO: communicate the error state back to the client in a sane way
		r.Error("Couldn't set up game: " + err.Error())
		return
	}

	r.Success(gin.H{
		"GameId":   game.Id(),
		"GameName": game.Name(),
	})
}

func (s *Server) listGamesHandler(c *gin.Context) {

	r := NewRenderer(c)
	s.doListGames(r)
}

func (s *Server) doListGames(r *Renderer) {
	r.Success(gin.H{
		"Games": s.storage.ListGames(100),
	})
}

func (s *Server) listManagerHandler(c *gin.Context) {
	r := NewRenderer(c)
	s.doListManager(r)
}

func (s *Server) doListManager(r *Renderer) {
	type agentInfo struct {
		Name        string
		DisplayName string
	}
	type managerInfo struct {
		Name              string
		DisplayName       string
		DefaultNumPlayers int
		Agents            []agentInfo
	}
	var managers []managerInfo
	for name, manager := range s.managers {
		agents := make([]agentInfo, len(manager.Agents()))
		for i, agent := range manager.Agents() {
			agents[i] = agentInfo{
				agent.Name(),
				agent.DisplayName(),
			}
		}
		managers = append(managers, managerInfo{
			Name:              name,
			DisplayName:       manager.Delegate().DisplayName(),
			DefaultNumPlayers: manager.Delegate().DefaultNumPlayers(),
			Agents:            agents,
		})
	}

	r.Success(gin.H{
		"Managers": managers,
	})

}

func (s *Server) gameVersionHandler(c *gin.Context) {

	game := s.getGame(c)

	playerIndex := s.effectivePlayerIndex(c)

	version := s.getRequestGameVersion(c)

	r := NewRenderer(c)

	s.doGameVersion(r, game, version, playerIndex)

}

func (s *Server) doGameVersion(r *Renderer, game *boardgame.Game, version int, playerIndex boardgame.PlayerIndex) {
	if game == nil {
		r.Error("Couldn't find game")
		return
	}

	if playerIndex == invalidPlayerIndex {
		r.Error("Got invalid playerIndex")
		return
	}

	state := game.State(version)

	//If state is nil, JSONForPlayer will basically treat it as just "give the
	//current version" which is a reasonable fallback.

	args := gin.H{
		"Game":            game.JSONForPlayer(playerIndex, state),
		"ViewingAsPlayer": playerIndex,
		"Forms":           s.generateForms(game),
	}

	r.Success(args)
}

//gameInfo is the first payload when a game is loaded, including immutables
//like chest, but also the initial game state payload as a convenience.
func (s *Server) gameInfoHandler(c *gin.Context) {

	game := s.getGame(c)

	playerIndex := s.effectivePlayerIndex(c)

	hasEmptySlots := s.getHasEmptySlots(c)

	r := NewRenderer(c)

	s.doGameInfo(r, game, playerIndex, hasEmptySlots)

}

type playerBoardInfo struct {
	DisplayName string
	IsAgent     bool
	IsEmpty     bool
	PhotoUrl    string
}

func (s *Server) gamePlayerInfo(game *boardgame.Game) []*playerBoardInfo {
	result := make([]*playerBoardInfo, game.NumPlayers())

	userIds := s.storage.UserIdsForGame(game.Id())
	agentNames := game.Agents()

	for i := range result {

		player := &playerBoardInfo{}

		result[i] = player

		if agentNames[i] != "" {
			agent := game.Manager().AgentByName(agentNames[i])

			if agent != nil {
				player.DisplayName = agent.DisplayName()
			}
			player.IsAgent = true
			player.IsEmpty = false
			continue
		}

		userId := userIds[i]

		if userId == "" {
			player.IsEmpty = true
			player.IsAgent = false
			player.DisplayName = ""
			continue
		}

		user := s.storage.GetUserById(userId)

		if user == nil {
			player.IsAgent = false
			player.IsEmpty = false
			player.DisplayName = "Unknown user"
			continue
		}

		player.IsAgent = false
		player.IsEmpty = false
		player.PhotoUrl = user.PhotoUrl
		player.DisplayName = user.EffectiveDisplayName()

		if player.DisplayName == "" {
			player.DisplayName = "Player " + strconv.Itoa(i)
		}

	}

	return result
}

func (s *Server) doGameInfo(r *Renderer, game *boardgame.Game, playerIndex boardgame.PlayerIndex, hasEmptySlots bool) {
	if game == nil {
		r.Error("Couldn't find game")
		return
	}

	if playerIndex == invalidPlayerIndex {
		r.Error("Got invalid playerIndex")
		return
	}

	args := gin.H{
		"Chest":           s.renderChest(game),
		"Forms":           s.generateForms(game),
		"Game":            game.JSONForPlayer(playerIndex, nil),
		"Error":           s.lastErrorMessage,
		"Players":         s.gamePlayerInfo(game),
		"ViewingAsPlayer": playerIndex,
		"HasEmptySlots":   hasEmptySlots,
	}

	s.lastErrorMessage = ""

	r.Success(args)

}

func (s *Server) moveHandler(c *gin.Context) {

	r := NewRenderer(c)

	if c.Request.Method != http.MethodPost {
		r.Error("This method only supports post.")
		return
	}

	game := s.getGame(c)

	if game == nil {
		r.Error("Game not found")
		return
	}

	proposer := s.effectivePlayerIndex(c)

	move, err := s.getMoveFromForm(c, game)

	if move == nil {

		//TODO: move this to doMakeMove once getMoveFromForm is refactored correctly.

		errString := "No move returned"

		if err != nil {
			errString = err.Error()
		}

		r.Error("Couldn't get move: " + errString)
		return
	}

	s.doMakeMove(r, game, proposer, move)

}

func (s *Server) doMakeMove(r *Renderer, game *boardgame.Game, proposer boardgame.PlayerIndex, move boardgame.Move) {

	if err := <-game.ProposeMove(move, proposer); err != nil {
		r.Error("Couldn't make move: " + err.Error())
		return
	}
	//TODO: it would be nice if we could show which fixup moves we made, too,
	//somehow.

	r.Success(nil)
}

func (s *Server) generateForms(game *boardgame.Game) []*MoveForm {

	var result []*MoveForm

	for _, factory := range game.Manager().PlayerMoveFactories() {

		move := factory(game.CurrentState())

		moveItem := &MoveForm{
			Name:        move.Name(),
			Description: move.Description(),
			Fields:      formFields(move),
		}
		result = append(result, moveItem)
	}

	return result
}

func formFields(move boardgame.Move) []*MoveFormField {

	var result []*MoveFormField

	for fieldName, fieldType := range move.ReadSetter().Props() {

		val, _ := move.ReadSetter().Prop(fieldName)

		result = append(result, &MoveFormField{
			Name:         fieldName,
			Type:         fieldType,
			DefaultValue: val,
		})

	}

	return result
}

func (s *Server) renderChest(game *boardgame.Game) map[string][]interface{} {
	//Substantially copied from cli.renderChest().

	deck := make(map[string][]interface{})

	for _, name := range game.Chest().DeckNames() {

		components := game.Chest().Deck(name).Components()

		values := make([]interface{}, len(components))

		for i, component := range components {
			values[i] = struct {
				Index  int
				Values interface{}
			}{
				i,
				component.Values,
			}
		}

		deck[name] = values
	}

	return deck
}

//genericHandler doesn't do much. We just register it so we automatically get
//CORS handlers triggered with the middelware.
func (s *Server) genericHandler(c *gin.Context) {
	r := NewRenderer(c)
	r.Success(gin.H{
		"Message": "Nothing to see here.",
	})
}

//Start is where you start the server, and it never returns until it's time to shut down.
func (s *Server) Start() {

	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		log.Println("Couldn't find a " + configFileName + " in current directory. This file is required. Copy a starter one from boardgame/server/api/config.SAMPLE.json")
		return
	}

	contents, err := ioutil.ReadFile(configFileName)

	if err != nil {
		log.Println("Couldn't read config file:", err)
		return
	}

	var config Config

	if err := json.Unmarshal(contents, &config); err != nil {
		log.Println("couldn't unmarshal config file:", err)
		return
	}

	log.Println("Environment Variables")
	//Dbug print out the current environment
	for _, config := range os.Environ() {
		log.Println("Environ:", config)
	}

	if v := os.Getenv("GIN_MODE"); v == "release" {
		log.Println("Using release mode config")
		s.config = config.Prod
	} else {
		log.Println("Using dev mode config")
		s.config = config.Dev
	}

	if err := s.config.Validate(); err != nil {
		log.Println("The provided config was not valid: ", err)
		return
	}

	name := s.storage.Name()

	storageConfig := s.config.StorageConfig[name]

	log.Println("Connecting to storage", name, "with config '"+storageConfig+"'")

	if err := s.storage.Connect(storageConfig); err != nil {
		log.Println("Couldnt' connect to storage manager: ", err)
		return
	}

	router := gin.New()

	router.Use(gin.Recovery(), gin.LoggerWithWriter(os.Stdout, "/_ah/health"))

	router.NoRoute(s.genericHandler)
	router.Use(cors.Middleware(cors.Config{
		Origins:        s.config.AllowedOrigins,
		RequestHeaders: "content-type, Origin",
		ExposedHeaders: "content-type",
		Methods:        "GET, POST",
		Credentials:    true,
	}))

	//The status endpoint for a game gets POUNDED so we don't want to touch
	//the database if we can help it.
	router.GET("/api/game/:name/:id/status", s.gameStatusHandler)

	//We have everything prefixed by /api just in case at some point we do
	//want to host both static and api on the same logical server.
	mainGroup := router.Group("/api")
	mainGroup.Use(s.userSetup)

	{
		mainGroup.GET("list/game", s.listGamesHandler)
		mainGroup.GET("list/manager", s.listManagerHandler)

		mainGroup.POST("auth", s.authCookieHandler)

		protectedMainGroup := mainGroup.Group("")
		protectedMainGroup.Use(s.requireLoggedIn)
		protectedMainGroup.POST("new/game", s.newGameHandler)

		gameAPIGroup := mainGroup.Group("game/:name/:id")
		gameAPIGroup.Use(s.gameAPISetup)
		{
			gameAPIGroup.GET("info", s.gameInfoHandler)
			gameAPIGroup.GET("version/:version", s.gameVersionHandler)

			//The statusHandler is conceptually here, but becuase we want to
			//optimize it so much we have it congfigured at the top level.

			protectedGameAPIGroup := gameAPIGroup.Group("")
			protectedGameAPIGroup.Use(s.requireLoggedIn)
			protectedGameAPIGroup.POST("move", s.moveHandler)
			protectedGameAPIGroup.POST("join", s.joinGameHandler)
		}
	}

	if p := os.Getenv("PORT"); p != "" {
		router.Run(":" + p)
	} else {
		router.Run(":" + s.config.DefaultPort)
	}

}
