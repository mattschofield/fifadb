package main

import (
	"fmt"
)

type Player struct {
	ID                      int    `csv:"id"`
	FirstName               string `csv:"first_name"`
	LastName                string `csv:"last_name"`
	Position                string `csv:"position"`
	Club                    string `csv:"club"`
	League                  string `csv:"league"`
	Nationality             string `csv:"nationality"`
	Rating                  int    `csv:"rating"`
	LowestPrice             int    `csv:"lowest_price"`
	ImageURL                string `csv:"image_url"`
	FUTHeadTotalStats       int    `csv:"futhead_total_stats"`
	FUTHeadTotalInGameStats int    `csv:"futhead_total_in_game_stats"`
	HeightInCM              int    `csv:"height_in_cm"`
	StrongFoot              string `csv:"strong_foot"`
	WorkRateOffensive       string `csv:"workrate_offensive"`
	WorkRateDefensive       string `csv:"workrate_defensive"`
	SkillMoves              int    `csv:"skill_moves"`
	WeakFoot                int    `csv:"weak_foot"`
	Pace                    int    `csv:"pace"`
	PaceAttributes
	Shooting int `csv:"shooting"`
	ShootingAttributes
	Passing int `csv:"passing"`
	PassingAttributes
	Dribbling int `csv:"dribbling"`
	DribblingAttributes
	Defending int `csv:"defending"`
	DefendingAttributes
	Physical int `csv:"physical"`
	PhysicalAttributes
	// TODO(mattschofield): find a way to pull these stats
	// FUTHeadAttackerRating   int
	// FUTHeadCreatorRating    int
	// FUTHeadDefenderRating   int
	// FUTHeadBeastRating      int
	// FUTHeadHeadingRating    int
}

type PaceAttributes struct {
	Acceleration int `csv:"pac_acceleration"`
	SprintSpeed  int `csv:"pac_sprint_speed"`
}

type ShootingAttributes struct {
	Positioning int `csv:"sho_positioning"`
	Finishing   int `csv:"sho_finishing"`
	ShotPower   int `csv:"sho_shot_power"`
	LongShots   int `csv:"sho_long_shots"`
	Volleys     int `csv:"sho_volleys"`
	Penalties   int `csv:"sho_penalties"`
}

type PassingAttributes struct {
	Vision       int `csv:"pas_vision"`
	Crossing     int `csv:"pas_crossing"`
	FreeKick     int `csv:"pas_free_kick"`
	ShortPassing int `csv:"pas_short_passing"`
	LongPassing  int `csv:"pas_long_passing"`
	Curve        int `csv:"pas_curve"`
}

type DribblingAttributes struct {
	Agility     int `csv:"dri_agility"`
	Balance     int `csv:"dri_balance"`
	Reactions   int `csv:"dri_reactions"`
	BallControl int `csv:"dri_ball_control"`
	Dribbling   int `csv:"dri_dribbling"`
	Composure   int `csv:"dri_composure"`
}

type DefendingAttributes struct {
	Interceptions  int `csv:"def_interceptions"`
	Heading        int `csv:"def_heading"`
	Marking        int `csv:"def_marking"`
	StandingTackle int `csv:"def_standing_tackle"`
	SlidingTackle  int `csv:"def_sliding_tackle"`
}

type PhysicalAttributes struct {
	Jumping    int `csv:"phy_jumping"`
	Stamina    int `csv:"phy_stamina"`
	Strength   int `csv:"phy_strength"`
	Aggression int `csv:"phy_aggression"`
}

func (p *Player) Name() string {
	n := p.LastName

	if p.FirstName != "" {
		n = fmt.Sprintf("%s %s", p.FirstName, n)
	}

	return n
}

func (p *Player) String() string {
	return fmt.Sprintf("id:%d - %d %s (%s) - %s, %s", p.ID, p.Rating, p.Name(), p.Nationality, p.League, p.Club)
}
