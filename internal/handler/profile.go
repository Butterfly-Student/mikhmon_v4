package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"mikhmon_v4/config"
	ros "mikhmon_v4/internal/routeros"
)

// ProfileHandler manages hotspot user profiles.
type ProfileHandler struct {
	store sessions.Store
	pool  *ros.Pool
}

func NewProfileHandler(store sessions.Store, pool *ros.Pool) *ProfileHandler {
	return &ProfileHandler{store: store, pool: pool}
}

// GetProfiles lists all hotspot user profiles.
func (h *ProfileHandler) GetProfiles(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	profiles, err := ros.RunArgs(client, "/ip/hotspot/user/profile/print")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

// GetProfile retrieves a single profile by .id.
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	profiles, err := ros.RunArgs(client, "/ip/hotspot/user/profile/print", "?.id="+c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(profiles) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, profiles[0])
}

// AddProfile creates a new hotspot user profile.
func (h *ProfileHandler) AddProfile(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var body ProfileRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	onLogin := buildOnLoginScript(body)

	args := []string{
		"/ip/hotspot/user/profile/add",
		"=name=" + sanitizeName(body.Name),
		"=address-pool=" + body.AddressPool,
		"=rate-limit=" + body.RateLimit,
		"=shared-users=" + body.SharedUsers,
		"=status-autorefresh=1m",
		"=on-login=" + onLogin,
		"=parent-queue=" + body.ParentQueue,
	}

	result, err := ros.RunArgs(client, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}

	// Fetch created profile.
	newID := ""
	if len(result) > 0 {
		newID = result[0]["ret"]
	}
	if newID == "" {
		c.JSON(http.StatusOK, gin.H{"message": "success", "data": gin.H{}})
		return
	}

	profiles, _ := ros.RunArgs(client, "/ip/hotspot/user/profile/print", "?.id="+newID)
	data := map[string]string{}
	if len(profiles) > 0 {
		data = profiles[0]
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": data})
}

// UpdateProfile updates an existing hotspot user profile.
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var body ProfileRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	onLogin := buildOnLoginScript(body)
	id := c.Param("id")

	args := []string{
		"/ip/hotspot/user/profile/set",
		"=.id=" + id,
		"=name=" + sanitizeName(body.Name),
		"=address-pool=" + body.AddressPool,
		"=rate-limit=" + body.RateLimit,
		"=shared-users=" + body.SharedUsers,
		"=status-autorefresh=1m",
		"=on-login=" + onLogin,
		"=parent-queue=" + body.ParentQueue,
	}

	_, err = ros.RunArgs(client, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}

	profiles, _ := ros.RunArgs(client, "/ip/hotspot/user/profile/print", "?.id="+id)
	data := map[string]string{}
	if len(profiles) > 0 {
		data = profiles[0]
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": data})
}

// RemoveProfile deletes a hotspot user profile.
func (h *ProfileHandler) RemoveProfile(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = ros.RunArgs(client, "/ip/hotspot/user/profile/remove", "=.id="+c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetAddressPools lists IP pools on the router.
func (h *ProfileHandler) GetAddressPools(c *gin.Context) {
	h.simpleGet(c, "/ip/pool/print")
}

// GetQueues lists simple queues on the router.
func (h *ProfileHandler) GetQueues(c *gin.Context) {
	h.simpleGet(c, "/queue/simple/print")
}

// GetNATRules lists NAT rules on the router.
func (h *ProfileHandler) GetNATRules(c *gin.Context) {
	h.simpleGet(c, "/ip/firewall/nat/print")
}

// GetInterfaces lists interfaces on the router.
func (h *ProfileHandler) GetInterfaces(c *gin.Context) {
	h.simpleGet(c, "/interface/print")
}

func (h *ProfileHandler) simpleGet(c *gin.Context, cmd string) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := ros.RunArgs(client, cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ── On-login script builder ────────────────────────────────────────────────

// ProfileRequest is the JSON body for add/update profile operations.
type ProfileRequest struct {
	Name         string `json:"name"`
	SharedUsers  string `json:"sharedusers"`
	RateLimit    string `json:"ratelimit"`
	ExpMode      string `json:"expmode"`
	Validity     string `json:"validity"`
	Price        string `json:"price"`
	SellingPrice string `json:"sellingprice"`
	AddressPool  string `json:"addresspool"`
	LockUser     string `json:"lockuser"`
	LockServer   string `json:"lockserver"`
	ParentQueue  string `json:"parentqueue"`
}

// buildOnLoginScript generates the RouterOS on-login script from a ProfileRequest.
// Directly ported from post/post_add_userprofile.php and post/post_update_userprofile.php.
func buildOnLoginScript(body ProfileRequest) string {
	price := body.Price
	if price == "" {
		price = "0"
	}
	sprice := body.SellingPrice
	if sprice == "" {
		sprice = "0"
	}

	// MAC lock script segment.
	lock := ""
	if body.LockUser == "Enable" {
		lock = `; [:local mac $"mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]]`
	}

	// Server lock script segment.
	slock := ""
	if body.LockServer != "Disable" && body.LockServer != "" {
		slock = `; [:local mac $"mac-address"; :local srv [/ip hotspot host get [find where mac-address="$mac"] server]; /ip hotspot user set server=$srv [find where name=$user]]`
	}

	expmode := body.ExpMode
	validity := strings.ToLower(body.Validity)

	// Sale recording script segment.
	record := fmt.Sprintf(
		`; :local mac $"mac-address"; :local time [/system clock get time ]; /system script add name="$date-|-$time-|-$user-|-`+
			`%s-|-$address-|-$mac-|-`+validity+`-|-`+sanitizeName(body.Name)+`-|-$comment" owner="$month$year" source=$date comment=mikhmon`,
		price,
	)

	mode := ""
	switch expmode {
	case "ntf", "ntfc":
		mode = "N"
	case "rem", "remc":
		mode = "X"
	}

	if expmode == "" || expmode == "0" {
		if price != "" && price != "0" {
			return fmt.Sprintf(`:put (",,` + price + `,,` + sprice + `,noexp,` + body.LockUser + `,` + body.LockServer + `,")` + lock + slock)
		}
		return ""
	}

	// Base on-login script (expiration monitoring).
	base := fmt.Sprintf(
		`:put (",` + expmode + `,` + price + `,` + validity + `,` + sprice + `,,` + body.LockUser + `,` + body.LockServer + `,"); ` +
			`:local mode "` + mode + `"; ` +
			`{:local date [ /system clock get date ];` +
			`:local year [ :pick $date 7 11 ];` +
			`:local month [ :pick $date 0 3 ];` +
			`:local comment [ /ip hotspot user get [/ip hotspot user find where name="$user"] comment];` +
			` :local ucode [:pic $comment 0 2];` +
			` :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={` +
			` /sys sch add name="$user" disable=no start-date=$date interval="` + validity + `";` +
			` :delay 2s;` +
			` :local exp [ /sys sch get [ /sys sch find where name="$user" ] next-run];` +
			` :local getxp [len $exp];` +
			` :if ($getxp = 15) do={` +
			` :local d [:pic $exp 0 6]; :local t [:pic $exp 7 16]; :local s ("/"); :local exp ("$d$s$year $t");` +
			` /ip hotspot user set comment="$exp ` + mode + `" [find where name="$user"];};` +
			` :if ($getxp = 8) do={ /ip hotspot user set comment="$date $exp ` + mode + `" [find where name="$user"];};` +
			` :if ($getxp > 15) do={ /ip hotspot user set comment="$exp ` + mode + `" [find where name="$user"];};` +
			` /sys sch remove [find where name="$user"]`,
	)

	switch expmode {
	case "rem", "ntf":
		return base + lock + slock + "}}"
	case "remc", "ntfc":
		return base + record + lock + slock + "}}"
	}
	return ""
}

func sanitizeName(s string) string {
	return strings.ReplaceAll(s, " ", "-")
}
