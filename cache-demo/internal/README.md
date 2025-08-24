if you have the go dependencies installed
go run main.go

Self Notes: 

create a client: client := &http.Client{Timeout: ...} 
create request: req, err := http.NewRequestWithContext(ctx, method, url, body(generally nil))
execute request: resp, err := client.Do(req)
defer resp.Body.Close() to ensure body is closed after reading
Use structs with  `json:"" ` tags to parse json 
json.NewDecoder(resp.Body).Decode(&gridData)
decode the resp.body into a instance of the json struct 
caches require a mutex, a RWmutex is more precise 
Lock, unlock when editing values, readlock readUnlock when reading. 
Don't defer unlocks in loops, creates problems. 
Bucket rate limiter algorithm works by filling a channel with empty structs which don't exhuast memory 
goroutine runs in the background on a ticker, filling the bucket with more structs if possible on tick. 
On request, shoot out of the channel, empty channels should signify the rate limit has been exceeded. 