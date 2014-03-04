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

import (
        "dlib/dbus"
        "dlib/logger"
        "os"
)

func main() {
        defer func() {
                if err := recover(); err != nil {
                        logger.Println("recover err:", err)
                }
        }()

        trieMD5Map = make(map[string]*Trie)
        strsMD5Map = make(map[string][]*TrieInfo)
        nameMD5Map = make(map[string]string)
        m := &Pinyin{}
        err := dbus.InstallOnSession(m)
        if err != nil {
                logger.Println("Install Pinyin DBus Session Failed:", err)
                panic(err)
        }

        t := &PinyinTrie{}
        err = dbus.InstallOnSession(t)
        if err != nil {
                logger.Println("Install Pinyin Trie DBus Session Failed:", err)
                panic(err)
        }
        dbus.DealWithUnhandledMessage()

        //select {}
        if err = dbus.Wait(); err != nil {
                logger.Println("lost dbus session:", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}
