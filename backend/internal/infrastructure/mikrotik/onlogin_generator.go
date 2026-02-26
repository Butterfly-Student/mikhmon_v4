package mikrotik

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// OnLoginGenerator generates and parses on-login scripts for MikroTik user profiles
type OnLoginGenerator struct{}

// NewOnLoginGenerator creates a new on-login script generator
func NewOnLoginGenerator() *OnLoginGenerator {
	return &OnLoginGenerator{}
}

// Generate creates an on-login script for user profile
// This script is executed every time a user logs in
// INI ADALAH FITUR PALING PENTING DARI MIKHMON!
func (g *OnLoginGenerator) Generate(req *dto.ProfileRequest) string {
	// Build script components
	components := g.buildScriptComponents(req)

	// Combine all components into final script
	return g.assembleScript(components)
}

// ScriptComponents holds all parts of the on-login script
type ScriptComponents struct {
	Header     string // Output metadata for Mikhmon
	Expiration string // Expiration logic
	Record     string // Sales recording to /system/script
	LockUser   string // MAC address lock
	LockServer string // Server lock
	Footer     string // Script closing
}

// buildScriptComponents builds all parts of the script
func (g *OnLoginGenerator) buildScriptComponents(req *dto.ProfileRequest) *ScriptComponents {
	comp := &ScriptComponents{}

	// 1. Header with metadata output
	comp.Header = g.buildHeader(req)

	// 2. Expiration logic
	comp.Expiration = g.buildExpirationLogic(req)

	// 3. Sales recording (if remc or ntfc mode)
	if req.ExpireMode == "remc" || req.ExpireMode == "ntfc" {
		comp.Record = g.buildRecordScript(req)
	}

	// 4. MAC Lock script
	if req.LockUser == "Enable" {
		comp.LockUser = g.buildLockUserScript()
	}

	// 5. Server Lock script
	if req.LockServer == "Enable" {
		comp.LockServer = g.buildLockServerScript()
	}

	// 6. Footer based on expire mode
	comp.Footer = g.buildFooter(req)

	return comp
}

// buildHeader creates the script header with metadata output
// Format: :put (",expire_mode,price,validity,selling_price,,lock_user,lock_server,");
func (g *OnLoginGenerator) buildHeader(req *dto.ProfileRequest) string {
	// Determine mode character
	var mode string
	switch req.ExpireMode {
	case "ntf", "ntfc":
		mode = "N"
	case "rem", "remc":
		mode = "X"
	default:
		mode = ""
	}

	// Build :put output with metadata
	return fmt.Sprintf(
		`:put (",%s,%.0f,%s,%.0f,,%s,%s,"); :local mode "%s"; {`,
		req.ExpireMode,
		req.Price,
		req.Validity,
		req.SellingPrice,
		req.LockUser,
		req.LockServer,
		mode,
	)
}

// buildExpirationLogic creates the expiration handling script
func (g *OnLoginGenerator) buildExpirationLogic(req *dto.ProfileRequest) string {
	if req.ExpireMode == "0" || req.ExpireMode == "" {
		return ""
	}

	// Script to add scheduler, calculate expiration, and update comment
	script := fmt.Sprintf(`
    :local date [/system clock get date];
    :local year [:pick $date 7 11];
    :local month [:pick $date 0 3];
    :local comment [/ip hotspot user get [/ip hotspot user find where name="$user"] comment];
    :local ucode [:pic $comment 0 2];
    
    :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={
        /sys sch add name="$user" disable=no start-date=$date interval="%s";
        :delay 2s;
        :local exp [/sys sch get [/sys sch find where name="$user"] next-run];
        :local getxp [len $exp];
        
        :if ($getxp = 16) do={
            :local d [:pic $exp 0 6];
            :local t [:pic $exp 7 15];
            :local s ("/");
            :local exp ("$d$s$year $t");
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };
        
        :if ($getxp = 8) do={
            /ip hotspot user set comment="$date $exp $mode" [find where name="$user"];
        };
        
        :if ($getxp > 16) do={
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };
        
        /sys sch remove [find where name="$user"];
    }
`, req.Validity)

	return script
}

// buildRecordScript creates the sales recording script
// Records transaction to /system/script for reporting
func (g *OnLoginGenerator) buildRecordScript(req *dto.ProfileRequest) string {
	// Script name format: $date-|-$time-|-$user-|-$price-|-$address-|-$mac-|-$validity-|-$profile-|-$comment
	return fmt.Sprintf(`
    :local mac $"mac-address";
    :local time [/system clock get time];
    /system script add name="$date-|-$time-|-$user-|-%.0f-|-$address-|-$mac-|-%s-|-%s-|-$comment" owner="$month$year" source=$date comment=mikhmon;`,
		req.Price,
		req.Validity,
		req.Name,
	)
}

// buildLockUserScript creates the MAC address lock script
func (g *OnLoginGenerator) buildLockUserScript() string {
	return `; [:local mac $"mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]]`
}

// buildLockServerScript creates the server lock script
func (g *OnLoginGenerator) buildLockServerScript() string {
	return `; [:local mac $"mac-address"; :local srv [/ip hotspot host get [find where mac-address="$mac"] server]; /ip hotspot user set server=$srv [find where name=$user]]`
}

// buildFooter creates the script footer based on expire mode
func (g *OnLoginGenerator) buildFooter(req *dto.ProfileRequest) string {
	switch req.ExpireMode {
	case "rem", "remc":
		// Remove mode - close script with }
		return "}"
	case "ntf", "ntfc":
		// Notice mode - same as remove
		return "}"
	case "0":
		// No expiration - wrap lock scripts if any exist
		if req.LockUser == "Enable" || req.LockServer == "Enable" {
			return "}"
		}
		return ""
	default:
		return "}"
	}
}

// assembleScript combines all components into final script
func (g *OnLoginGenerator) assembleScript(comp *ScriptComponents) string {
	var parts []string

	// Start with header
	parts = append(parts, comp.Header)

	// Add expiration logic if exists
	if comp.Expiration != "" {
		parts = append(parts, comp.Expiration)
	}

	// Add record script if exists (before locks)
	if comp.Record != "" {
		parts = append(parts, comp.Record)
	}

	// Add lock scripts
	if comp.LockUser != "" {
		parts = append(parts, comp.LockUser)
	}
	if comp.LockServer != "" {
		parts = append(parts, comp.LockServer)
	}

	// Add footer
	if comp.Footer != "" {
		parts = append(parts, comp.Footer)
	}

	// Join all parts
	script := strings.Join(parts, "")

	// Clean up whitespace
	script = strings.TrimSpace(script)

	return script
}

// Parse extracts Mikhmon metadata from an existing on-login script
func (g *OnLoginGenerator) Parse(script string) *dto.ProfileRequest {
	req := &dto.ProfileRequest{}

	if script == "" {
		return req
	}

	// Parse :put header for metadata
	// Generated format: :put (",ntf,10000,1d,15000,,Enable,Disable,");
	// The double-quote wraps the entire CSV: :put ("<csv>");
	putPattern := regexp.MustCompile(`:put \(",([\w]*),([\d\.]*),([^,]*),([\d\.]*),,([^,]*),([^,]*),"\)`)
	matches := putPattern.FindStringSubmatch(script)

	if len(matches) >= 7 {
		req.ExpireMode = matches[1]
		req.Price = parseFloat(matches[2])
		req.Validity = matches[3]
		req.SellingPrice = parseFloat(matches[4])
		req.LockUser = matches[5]
		req.LockServer = matches[6]
	}

	// Detect record script presence
	if strings.Contains(script, "/system script add") {
		// If has record script but mode is rem -> remc, ntf -> ntfc
		if req.ExpireMode == "rem" {
			req.ExpireMode = "remc"
		} else if req.ExpireMode == "ntf" {
			req.ExpireMode = "ntfc"
		}
	}

	return req
}

// GenerateExpiredAction generates the action script for when user expires
// This is used when setting up expiration monitoring
func (g *OnLoginGenerator) GenerateExpiredAction(expireMode string) string {
	switch expireMode {
	case "rem", "remc":
		// Remove user when expired
		return "/ip hotspot user remove [find name=$user]"
	case "ntf", "ntfc":
		// Set limit to 1s (effectively disable)
		return "/ip hotspot user set limit-uptime=1s [find name=$user]"
	default:
		return ""
	}
}

// GenerateExpireMonitorScript generates the global scheduler script used by
// "Mikhmon-Expire-Monitor" to enforce expired users handling.
func (g *OnLoginGenerator) GenerateExpireMonitorScript() string {
	return `:local dateint do={:local montharray ("jan","feb","mar","apr","may","jun","jul","aug","sep","oct","nov","dec"); :local days [:pick $d 4 6]; :local month [:pick $d 0 3]; :local year [:pick $d 7 11]; :local monthint ([:find $montharray $month]); :local month ($monthint + 1); :if ([len $month] = 1) do={:local zero ("0"); :return [:tonum ("$year$zero$month$days")];} else={:return [:tonum ("$year$month$days")];}}; :local timeint do={:local hours [:pick $t 0 2]; :local minutes [:pick $t 3 5]; :return ($hours * 60 + $minutes);}; :local date [/system clock get date]; :local time [/system clock get time]; :local today [$dateint d=$date]; :local curtime [$timeint t=$time]; :local tyear [:pick $date 7 11]; :local lyear ($tyear - 1); :foreach i in=[/ip hotspot user find where comment~"/$tyear" || comment~"/$lyear"] do={:local comment [/ip hotspot user get $i comment]; :local limit [/ip hotspot user get $i limit-uptime]; :local name [/ip hotspot user get $i name]; :local gettime [:pick $comment 12 20]; :if ([:pick $comment 3] = "/" and [:pick $comment 6] = "/") do={:local expd [$dateint d=$comment]; :local expt [$timeint t=$gettime]; :if ((($expd < $today and $expt < $curtime) or ($expd < $today and $expt > $curtime) or ($expd = $today and $expt < $curtime)) and $limit != "00:00:01") do={:if ([:pick $comment 21] = "N") do={/ip hotspot user set limit-uptime=1s $i; /ip hotspot active remove [find where user=$name];} else={/ip hotspot user remove $i; /ip hotspot active remove [find where user=$name];}}}}`
}

// Example outputs for documentation:
/*

Example 1: Profile "Premium" with remove mode, lock user, record sales
------------------------------------------------------------------------
Input:
  Name: Premium
  ExpireMode: remc
  Validity: 30d
  Price: 5000
  SellingPrice: 5500
  LockUser: Enable
  LockServer: Disable

Output Script:
:put (",remc,5000,30d,5500,,Enable,Disable,"); :local mode "X"; {
    :local date [/system clock get date];
    :local year [:pick $date 7 11];
    :local month [:pick $date 0 3];
    :local comment [/ip hotspot user get [/ip hotspot user find where name="$user"] comment];
    :local ucode [:pic $comment 0 2];

    :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={
        /sys sch add name="$user" disable=no start-date=$date interval="30d";
        :delay 2s;
        :local exp [/sys sch get [/sys sch find where name="$user"] next-run];
        :local getxp [len $exp];

        :if ($getxp = 15) do={
            :local d [:pic $exp 0 6];
            :local t [:pic $exp 7 16];
            :local s ("/");
            :local exp ("$d$s$year $t");
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };

        :if ($getxp = 8) do={
            /ip hotspot user set comment="$date $exp $mode" [find where name="$user"];
        };

        :if ($getxp > 15) do={
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };

        /sys sch remove [find where name="$user"];
    :local mac $"mac-address";
    :local time [/system clock get time];
    /system script add name="$date-|-$time-|-$user-|-5000-|-$address-|-$mac-|-30d-|-Premium-|-$comment" owner="$month$year" source=$date comment=mikhmon;
; [:local mac $"mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]]
}}

Example 2: Profile "Basic" with notice mode, no lock, no record
-----------------------------------------------------------------
Input:
  Name: Basic
  ExpireMode: ntf
  Validity: 7d
  Price: 2000
  SellingPrice: 2500
  LockUser: Disable
  LockServer: Disable

Output Script:
:put (",ntf,2000,7d,2500,,Disable,Disable,"); :local mode "N"; {
    :local date [/system clock get date];
    :local year [:pick $date 7 11];
    :local month [:pick $date 0 3];
    :local comment [/ip hotspot user get [/ip hotspot user find where name="$user"] comment];
    :local ucode [:pic $comment 0 2];

    :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={
        /sys sch add name="$user" disable=no start-date=$date interval="7d";
        :delay 2s;
        :local exp [/sys sch get [/sys sch find where name="$user"] next-run];
        :local getxp [len $exp];

        :if ($getxp = 15) do={
            :local d [:pic $exp 0 6];
            :local t [:pic $exp 7 16];
            :local s ("/");
            :local exp ("$d$s$year $t");
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };

        :if ($getxp = 8) do={
            /ip hotspot user set comment="$date $exp $mode" [find where name="$user"];
        };

        :if ($getxp > 15) do={
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };

        /sys sch remove [find where name="$user"];
}}

*/
