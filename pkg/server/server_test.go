/*
Copyright 2018 The OpenEBS Author

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server_test

/*
func TestStartHTTPServer(t *testing.T) {

	ErrorMessages := make(chan error)
	go func() {
		//Block port 9191 and attempt to start http server at 9191.
		p1, err := net.Listen("tcp", "localhost:9191")
		defer p1.Close()
		if err != nil {
			t.Log(err)
		}
		server.ListenPort = ":9191"
		ErrorMessages <- server.StartHTTPServer()
	}()
	msg := <-ErrorMessages
	if msg != nil {
		t.Log("Try to start http server in a port which is busy.")
		t.Log(msg)
	}
}
*/
