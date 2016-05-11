//
//  MatchData.swift
//  admin
//
//  Created by tassar on 5/6/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import ObjectMapper

class PlayerData: Mappable {
	var id: UInt!
	var createdAt: String!
	var name: String?
	var gold: Int!
	var lostGold: Int!
	var energy: Double!
	var combo: Int!
	var grade: String!
	var level: Int!
	var levelData: String!
	var hitCount: Int!
	var cid: String!
	var questionInfo: String!
	var answered: Int!
	var questionCount: Int!
	required init?(_ map: Map) {
	}

	func mapping(map: Map) {
		id <- map["id"]
		createdAt <- map["createdAt"]
		name <- map["name"]
		gold <- map["gold"]
		lostGold <- map["lostGold"]
		energy <- map["energy"]
		combo <- map["combo"]
		grade <- map["grade"]
		level <- map["level"]
		levelData <- map["levelData"]
		hitCount <- map["hitCount"]
		cid <- map["cid"]
		questionInfo <- map["questionInfo"]
		answered <- map["answered"]
		questionCount <- map["questionCount"]
	}

	func getName() -> String {
		if name != nil && name!.characters.count > 0 {
			return name!
		} else {
			return cid.componentsSeparatedByString(":")[1]
		}
	}
}

enum MatchAnswerType: Int {
	case NotAnswer = 0, Answering, Answered
}

class MatchData: Mappable {
	var id: UInt!
	var createdAt: String!
	var mode: String!
	var elasped: Double!
	var gold: Int!
	var member: [PlayerData]!
	var rampageCount: Int!
	var answerType: MatchAnswerType!
	var teamID: String!

	required init?(_ map: Map) {
	}

	func mapping(map: Map) {
		id <- map["id"]
		createdAt <- map["createdAt"]
		mode <- map["mode"]
		elasped <- map["elasped"]
		gold <- map["gold"]
		member <- map["member"]
		rampageCount <- map["rampageCount"]
		answerType <- map["answerType"]
		teamID <- map["teamID"]
	}
}

class MatchResult: Mappable {
	var matchID: Int!
	var teamID: String!
	var matchData: MatchData!
	required init?(_ map: Map) {
	}

	func mapping(map: Map) {
		matchID <- map["matchID"]
		teamID <- map["teamID"]
		matchData <- map["matchData"]
	}
}