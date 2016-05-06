package core

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"strconv"
)

var _ = log.Printf

type ArenaPosition struct {
	X int
	Y int
}

type P ArenaPosition

type RealPosition struct {
	X float64
	Y float64
}

type RP RealPosition

type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

type Button struct {
	Id string `json:"id"`
	R  Rect   `json:"r"`
}

type RenderInfo struct {
	ArenaCellSize int
	ArenaBorder   int
	PlayerSize    float64
	WebScale      float64
	ButtonWidth   float64
	ButtonHeight  float64
	PlayerSpeed   float64
}

type MatchOptions struct {
	ArenaWidth        int        `json:"arenaWidth"`
	ArenaHeight       int        `json:"arenaHeight"`
	ArenaCellSize     int        `json:"arenaCellSize"`
	ArenaBorder       int        `json:"arenaBorder"`
	Warmup            float64    `json:"warmup"`
	ArenaEntrance     P          `json:"arenaEntrance"`
	ArenaExit         P          `json:"arenaExit"`
	PlayerSize        float64    `json:"playerSize"`
	WebScale          float64    `json:"webScale"`
	ButtonWidth       float64    `json:"buttonWidth"`
	ButtonHeight      float64    `json:"buttonHeight"`
	T1                float64    `json:"t1"`
	T2                float64    `json:"t2"`
	T3                float64    `json:"t3"`
	TRampage          float64    `json:"tRampage"`
	GoldBonus         [2]int     `json:"buttonBonus"`
	TouchPunish       [2]float64 `json:"touchPunish"`
	Mode2InitGold     [4]int     `json:"mode2InitGold"`
	Mode2GoldDropRate [4]int     `json:"mode2GoldDropRate"`
	MaxEnergy         float64    `json:"maxEnergy"`
	Mode1TotalTime    float64    `json:"mode1TotalTime"`
	WallRects         []Rect     `json:"walls"`
	Buttons           []*Button  `json:"buttons"`

	PlayerSpeed           float64       `json:"-"`
	Walls                 [][]int       `json:"-"`
	EnergyBonus           [4][4]float64 `json:"-"`
	InitButtonNum         [4]int        `json:"-"`
	ButtonHideTime        [2]float64    `json:"-"`
	RampageTime           [2]float64    `json:"-"`
	FirstComboInterval    [4]float64    `json:"-"`
	ComboInterval         [4]float64    `json:"-"`
	FirstComboExtra       float64       `json:"-"`
	ComboExtra            float64       `json:"-"`
	LaserSpeed            float64       `json:"-"`
	LaserSpeedup          float64       `json:"-"`
	EnergySpeedup         float64       `json:"-"`
	LaserAppearTime       float64       `json:"-"`
	LaserPauseTime        float64       `json:"-"`
	TileAdjacency         map[int][]int `json:"-"`
	PlayerInvincibleTime  float64       `json:"-"`
	Mode1TouchPunish      float64       `json:"-"`
	Mode2TouchPunish      int           `json:"-"`
	Mode2GoldDropInterval float64       `json:"-"`
	MainArduino           []string      `json:"-"`
	SubArduino            []string      `json:"-"`
}

type ScoreInfo [4]map[string]interface{}

var opt = DefaultMatchOptions()

func GetOptions() *MatchOptions {
	return opt
}

func GetScoreInfo() ScoreInfo {
	return [4]map[string]interface{}{
		map[string]interface{}{
			"time":   opt.T1,
			"status": "T1",
		},
		map[string]interface{}{
			"time":   opt.T2,
			"status": "T2",
		},
		map[string]interface{}{
			"time":   opt.T3,
			"status": "T3",
		},
		map[string]interface{}{
			"time":   opt.TRampage,
			"status": "T4",
		},
	}
}

func DefaultMatchOptions() *MatchOptions {
	var opt MatchOptions
	if _, err := toml.DecodeFile("cfg.toml", &opt); err != nil {
		log.Printf("parse toml error:%v\n", err.Error())
		os.Exit(1)
	}
	opt.buildWallRects()
	opt.buildButtons()
	opt.buildAdjacency()
	return &opt
}

func (m *MatchOptions) buildAdjacency() {
	adj := make(map[int][]int)
	w, h := m.ArenaWidth, m.ArenaHeight
	adjacentWith := func(a int, b int) bool {
		if b < 0 || b >= w*h {
			return false
		}
		if a%w == 0 && b == a-1 || b%w == 0 && a == b-1 {
			return false
		}
		for _, wall := range m.Walls {
			p1 := m.TilePosToInt(P{wall[0], wall[1]})
			p2 := m.TilePosToInt(P{wall[2], wall[3]})
			if a == p1 && b == p2 || a == p2 && b == p1 {
				return false
			}
		}
		return true
	}
	for i := 0; i < w*h; i++ {
		adj[i] = make([]int, 0)
		left := i - 1
		right := i + 1
		top := i - w
		bottom := i + w
		list := [4]int{left, right, top, bottom}
		for _, j := range list {
			if adjacentWith(i, j) {
				adj[i] = append(adj[i], j)
			}
		}
	}
	m.TileAdjacency = adj
}

func (m *MatchOptions) buildWallRects() {
	m.WallRects = make([]Rect, 0)
	for _, wall := range m.Walls {
		horizontal := wall[0] == wall[2]
		var w, h, x, y float64
		if horizontal {
			w = float64(m.ArenaCellSize + 2*m.ArenaBorder)
			h = float64(m.ArenaBorder)
			x = float64(wall[0]*(m.ArenaCellSize+m.ArenaBorder) - m.ArenaBorder/2)
			y = float64(MaxInt(wall[1], wall[3])*(m.ArenaCellSize+m.ArenaBorder) - m.ArenaBorder/2)
		} else {
			w = float64(m.ArenaBorder)
			h = float64(m.ArenaCellSize + 2*m.ArenaBorder)
			y = float64(wall[1]*(m.ArenaCellSize+m.ArenaBorder) - m.ArenaBorder/2)
			x = float64(MaxInt(wall[0], wall[2])*(m.ArenaCellSize+m.ArenaBorder) - m.ArenaBorder/2)
		}
		m.WallRects = append(m.WallRects, Rect{x, y, w, h})
	}
}

func (m *MatchOptions) buildButtons() {
	m.Buttons = make([]*Button, 0)
	// top and bottom wall
	c := m.ArenaCellSize
	b := m.ArenaBorder
	bw := m.ButtonWidth
	bh := m.ButtonHeight
	var x, y, w, h float64
	id := 0
	app := func() {
		m.Buttons = append(m.Buttons, &Button{strconv.Itoa(id), Rect{x, y, w, h}})
		id = id + 1
	}
	for i := 0; i < m.ArenaWidth; i++ {
		if m.ArenaEntrance.Y == 0 && i == m.ArenaEntrance.X {
			continue
		}
		if m.ArenaExit.Y == 0 && i == m.ArenaExit.X {
			continue
		}
		x = float64(c+b)*(float64(i)+0.5) - 0.5*bw
		y = float64(b) * 0.5
		w = bw
		h = bh
		app()
	}
	for i := 0; i < m.ArenaWidth; i++ {
		if m.ArenaEntrance.Y == m.ArenaHeight-1 && i == m.ArenaEntrance.X {
			continue
		}
		if m.ArenaExit.Y == m.ArenaHeight-1 && i == m.ArenaExit.X {
			continue
		}
		x = float64(c+b)*(float64(i)+0.5) - 0.5*bw
		y = float64((c+b)*m.ArenaHeight) - 0.5*float64(b) - bh
		w = bw
		h = bh
		app()
	}
	// left and right wall
	for i := 0; i < m.ArenaHeight; i++ {
		if m.ArenaEntrance.X == 0 && i == m.ArenaEntrance.Y {
			continue
		}
		if m.ArenaExit.X == 0 && i == m.ArenaExit.Y {
			continue
		}
		x = float64(b) * 0.5
		y = float64(c+b)*(float64(i)+0.5) - 0.5*bw
		w = bh
		h = bw
		app()
	}
	for i := 0; i < m.ArenaHeight; i++ {
		if m.ArenaEntrance.X == m.ArenaHeight-1 && i == m.ArenaEntrance.Y {
			continue
		}
		if m.ArenaExit.X == m.ArenaHeight-1 && i == m.ArenaExit.Y {
			continue
		}
		x = float64((c+b)*m.ArenaWidth) - 0.5*float64(b) - bh
		y = float64(c+b)*(float64(i)+0.5) - 0.5*bw
		w = bh
		h = bw
		app()
	}
	// inner wall
	for idx, rect := range m.WallRects {
		wall := m.Walls[idx]
		horizontal := wall[0] == wall[2]
		if horizontal {
			w = bw
			h = bh
			x = rect.X + float64(b) + 0.5*(float64(c)-bw)
			// above
			y = rect.Y - bh
			app()
			// below
			y = rect.Y + bh
			app()
		} else {
			w = bh
			h = bw
			y = rect.Y + float64(b) + 0.5*(float64(c)-bw)
			// left
			x = rect.X - bh
			app()
			// right
			x = rect.X + bh
			app()
		}
	}
}

func (m *MatchOptions) CollideWall(r *Rect) bool {
	for _, rect := range m.WallRects {
		if m.Collide(r, &rect) {
			return true
		}
	}
	return false
}

func (m *MatchOptions) Collide(r1 *Rect, r2 *Rect) bool {
	if r1.X < r2.X+r2.W &&
		r1.X+r1.W > r2.X &&
		r1.Y < r2.Y+r2.H &&
		r1.H+r1.Y > r2.Y {
		return true
	}
	return false
}

func (m *MatchOptions) PressingButtons(r *Rect) []string {
	ret := make([]string, 2)
	i := 0
	for _, btn := range m.Buttons {
		if m.Collide(r, &btn.R) {
			ret[i] = btn.Id
			i += 1
			if i == 2 {
				return ret
			}
		}
	}
	if i == 0 {
		return nil
	}
	return ret
}

func (m *MatchOptions) RealPosition(p P) RP {
	rp := RP{}
	rp.X = float64(m.ArenaCellSize+m.ArenaBorder) * (float64(p.X) + 0.5)
	rp.Y = float64(m.ArenaCellSize+m.ArenaBorder) * (float64(p.Y) + 0.5)
	return rp
}

func (m *MatchOptions) TilePosition(rp RP) (P, bool) {
	u := float64(m.ArenaCellSize + m.ArenaBorder)
	f := func(a float64) (int, bool) {
		i := 0
		for a >= u {
			a -= u
			i += 1
		}
		if (a >= float64(m.ArenaBorder/2)) && (a <= float64(m.ArenaBorder/2+m.ArenaCellSize)) {
			return i, true
		}
		return i, false
	}
	xI, xBool := f(rp.X)
	yI, yBool := f(rp.Y)
	return P{xI, yI}, xBool && yBool
}

func (m *MatchOptions) TilePosToInt(p P) int {
	return p.X + p.Y*m.ArenaWidth
}

func (m *MatchOptions) IntToTile(i int) P {
	return P{i % m.ArenaWidth, i / m.ArenaWidth}
}

func (m *MatchOptions) laserSpeed(energy float64) float64 {
	level := int(energy / m.EnergySpeedup)
	speed := m.LaserSpeed - float64(level)*m.LaserSpeedup
	return float64(m.ArenaCellSize) / 10 / speed
}
