package integration

// func TestCORS(t *testing.T) {
// 	fp, err := filepath.Abs("fixtures" + _const.FILESEPARATOR + "relay.json")
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
// 	b, err := ioutil.ReadFile(fp)
// 	res, err := http.Post("http://"+*dispatchU+"relay", "application/json", bytes.NewBuffer(b))
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
// 	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
// 		t.Fatalf("Access-Control != *")
// 	}
// 	if res.Header.Get("Access-Control-Allow-Headers") != "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization" {
// 		t.Fatalf("Access-Control-Allow-Headers!=Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	}
// }
