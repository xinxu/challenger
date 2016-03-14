package core

import (
  "strconv"
)

type MatchOptions struct {
  ArenaWidth     int       `json:"arenaWidth"`
  ArenaHeight    int       `json:"arenaHeight"`
  ArenaCellSize  int       `json:"arenaCellSize"`
  ArenaBorder    int       `json:"arenaBorder"`
  Warmup         int       `json:"warmup"`
  ArenaEntrance  P         `json:"arenaEntrance"`
  ArenaExit      P         `json:"arenaExit"`
  PlayerSize     float64   `json:"playerSize"`
  Web_ArenaScale float64   `json:"webScale"`
  ButtonWidth    float64   `json:"buttonWidth"`
  ButtonHeight   float64   `json:"buttonHeight"`
  T1             float64   `json:"t1"`
  T2             float64   `json:"t2"`
  T3             float64   `json:"t3"`
  TRampage       float64   `json:"tRampage"`
  WallRects      []Rect    `json:"walls"`
  Buttons        []*Button `json:"buttons"`
  // private
  playerSpeed   float64
  arenaWallList []W
}

func DefaultMatchOptions() *MatchOptions {
  v := MatchOptions{}
  v.ArenaWidth = 8
  v.ArenaHeight = 6
  v.ArenaCellSize = 140
  v.ArenaBorder = 24
  v.Warmup = 20
  v.ArenaEntrance = P{0, 4}
  v.ArenaExit = P{6, 0}
  v.PlayerSize = 50
  v.Web_ArenaScale = 0.5
  v.ButtonWidth = 60
  v.ButtonHeight = 30
  v.T1 = 2
  v.T2 = 2.2
  v.T3 = 2.5
  v.TRampage = 1
  v.playerSpeed = 200
  v.buildWallPoints()
  v.buildWallRects()
  v.buildButtons()
  return &v
}

func (m *MatchOptions) buildWallPoints() {
  m.arenaWallList = []W{
    W{P{4, 0}, P{5, 0}},
    W{P{1, 0}, P{1, 1}},
    W{P{6, 0}, P{6, 1}},
    W{P{0, 1}, P{1, 1}},
    W{P{2, 1}, P{3, 1}},
    W{P{3, 1}, P{4, 1}},
    W{P{6, 1}, P{7, 1}},
    W{P{2, 1}, P{2, 2}},
    W{P{5, 1}, P{5, 2}},
    W{P{2, 2}, P{3, 2}},
    W{P{3, 2}, P{4, 2}},
    W{P{4, 2}, P{5, 2}},
    W{P{1, 2}, P{1, 3}},
    W{P{6, 2}, P{6, 3}},
    W{P{0, 3}, P{1, 3}},
    W{P{2, 3}, P{3, 3}},
    W{P{4, 3}, P{5, 3}},
    W{P{6, 3}, P{7, 3}},
    W{P{2, 3}, P{2, 4}},
    W{P{3, 3}, P{3, 4}},
    W{P{0, 4}, P{1, 4}},
    W{P{1, 4}, P{2, 4}},
    W{P{4, 4}, P{5, 4}},
    W{P{5, 4}, P{6, 4}},
    W{P{6, 4}, P{7, 4}},
    W{P{4, 4}, P{4, 5}},
    W{P{2, 5}, P{3, 5}},
    W{P{5, 5}, P{6, 5}},
  }
}

func (m *MatchOptions) buildWallRects() {
  m.WallRects = make([]Rect, 0)
  for _, wall := range m.arenaWallList {
    horizontal := wall.P1.X == wall.P2.X
    var w, h, x, y float64
    if horizontal {
      w = float64(m.ArenaCellSize + 2*m.ArenaBorder)
      h = float64(m.ArenaBorder)
      x = float64(wall.P1.X*(m.ArenaCellSize+m.ArenaBorder) - m.ArenaBorder/2)
      y = float64(MaxInt(wall.P1.Y, wall.P2.Y)*(m.ArenaCellSize+m.ArenaBorder) - m.ArenaBorder/2)
    } else {
      w = float64(m.ArenaBorder)
      h = float64(m.ArenaCellSize + 2*m.ArenaBorder)
      y = float64(wall.P1.Y*(m.ArenaCellSize+m.ArenaBorder) - m.ArenaBorder/2)
      x = float64(MaxInt(wall.P1.X, wall.P2.X)*(m.ArenaCellSize+m.ArenaBorder) - m.ArenaBorder/2)
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
    wall := m.arenaWallList[idx]
    horizontal := wall.P1.X == wall.P2.X
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

func (m *MatchOptions) Collide(r *Rect) bool {
  for _, rect := range m.WallRects {
    if r.X < rect.X+rect.W &&
      r.X+r.W > rect.X &&
      r.Y < rect.Y+rect.H &&
      r.H+r.Y > rect.Y {
      return true
    }
  }
  return false
}

func (m *MatchOptions) RealPosition(p P) RP {
  rp := RP{}
  rp.X = float64((m.ArenaCellSize + m.ArenaBorder)) * (float64(p.X) + 0.5)
  rp.Y = float64((m.ArenaCellSize + m.ArenaBorder)) * (float64(p.Y) + 0.5)
  return rp
}
