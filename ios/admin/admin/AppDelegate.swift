//
//  AppDelegate.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import XCGLogger

let log = XCGLogger.defaultInstance()

@UIApplicationMain
class AppDelegate: UIResponder, UIApplicationDelegate {

	var window: UIWindow?

	func application(application: UIApplication, didFinishLaunchingWithOptions launchOptions: [NSObject: AnyObject]?) -> Bool {
		UITabBar.appearance().barTintColor = UIColor.clearColor()
		UITabBar.appearance().backgroundImage = UIImage()
		UITabBar.appearance().shadowImage = UIImage()
		WsClient.singleton.connect(PLConstants.getClientWsAddress())
		return true
	}
}
