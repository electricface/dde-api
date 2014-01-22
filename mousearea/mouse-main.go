/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

// #cgo amd63 386 CFLAGS: -g -Wall
// #cgo pkg-config: x11 xtst glib-2.0
// #include "mouse-record.h"
import "C"

import (
	"dlib"
	"dlib/dbus"
	"dlib/logger"
	"sync"
)

var (
	opMouse *Manager
	lock    sync.Mutex

	genID = func() func() int32 {
		id := int32(0)
		return func() int32 {
			lock.Lock()
			tmp := id
			id += 1
			lock.Unlock()
			return tmp
		}
	}()
)

func (op *Manager) RegisterArea(area []coordinateRange) int32 {
	cookie := genID()
	idRangeMap[cookie] = area

	return cookie
}

func (op *Manager) UnregisterArea(cookie int32) {
	delete(idRangeMap, cookie)
}

func (op *IdleTick) RigisterIdle(name string, timeout int32) int32 {
	cookie := genID()
	info := newTimerInfo(name, cookie, timeout)
	cookieTimerMap[cookie] = info
	go startTimer(cookie)

	return cookie
}

func (op *IdleTick) UnregisterIdle(cookie int32) {
	endTimer(cookie, true)
}

func NewManager() *Manager {
	return &Manager{}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover error:", err)
		}
	}()

	//cookieTimerMap = make(map[int32]*TimerInfo)
	idRangeMap = make(map[int32][]coordinateRange)
	//tmp := coordinateRange{X1: 1266, X2: 1370, Y1: 600, Y2: 767}
	//idRangeMap[0] = []coordinateRange{tmp}
	opMouse = NewManager()
	err := dbus.InstallOnSession(opMouse)
	if err != nil {
		logger.Println("Install DBus Session Failed:", err)
		panic(err)
	}

	/*
		err = dbus.InstallOnSession(idle)
		if err != nil {
			logger.Println("Install DBus Session Failed:", err)
			panic(err)
		}
	*/
	dbus.DealWithUnhandledMessage()
	C.record_init()
	defer C.record_finalize()

	dlib.StartLoop()
}