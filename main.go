//
// Copyright (c) 2016 Dennis Chen
//
// This file is part of Clashr.
//
// Clashr is free software: you can redistribute it and/or modify it under the
// terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// Clashr is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License for
// more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Clashr.  If not, see <http://www.gnu.org/licenses/>.
//

package main

import (
	"log"
	"net"
	"net/rpc"

	"github.com/clashr/judgrpcd/api"
)

func init() {
}

func main() {
	judge := new(api.Judge)

	server := rpc.NewServer()
	if err := server.Register(judge); err != nil {
		log.Fatalf("Format of service builder isn't correct. %s", err)
	}

	//start listening for messages on port 1234
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Couldn't start listening on port 1234. Error %s", err)
	}

	log.Println("Serving RPC handler")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error serving: %s", err)
		}

		go server.ServeConn(conn)
	}
}
