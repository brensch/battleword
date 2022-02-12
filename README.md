# battleword
wordle is cool right now

## setup
download the engine binary from the releases page, and run it from a commandline somewhere 
examples:
- windows: powershell
- mac: terminal
- linux: you know what to do
```
./battleword --names player1,player2 --apis http://localhost:8080,http://localhost:8081
```

## api
### /guess
the engine will hit your api here with the previous results of a game. you are expected to respond with your best guess.
#### Request:
```json
{
	"previous_words": ["beast", "lapse"],
	"previous_results": [
		[0, 1, 1, 2, 0],
		[0, 2, 1, 2, 2]
	],
	"times": [1012000, 336700]
}
```
#### Response:
```json
{
	"guess": "pause",
	"shout": "why is everybody shouting"
}
```
Shouts server no purpose except to intimidate your opponents.

### /results
Once all players are finished, the engine will send you the results of everyone in the match so you can brag. No response is required.
#### Request:
```json
{
	"players": [{
		"name": "brendo",
		"guesses": {
			"previous_words": ["beast", "lapse", "pause"],
			"previous_results": [
				[0, 1, 1, 2, 0],
				[0, 2, 1, 2, 2],
				[2, 2, 2, 2, 2]
			],
			"times": [1012000, 336700, 473626]
		}
	}, {
		"name": "skye",
		"guesses": {
			"previous_words": ["beast", "beast", "beast", "beast", "beast", "beast"],
			"previous_results": [
				[0, 1, 1, 2, 0],
				[0, 1, 1, 2, 0],
				[0, 1, 1, 2, 0],
				[0, 1, 1, 2, 0],
				[0, 1, 1, 2, 0],
				[0, 1, 1, 2, 0]
			],
			"times": [1012000, 336700, 232300, 160800, 282900, 290000]
		}
	}],
	"answer": "pause"
}
```

### notes
 the logic for the yellows when multiple letters are incorrect is not correct RN. I need to figure out what wordle does then make it do the same.

### releasing
 in order to release, run:
```
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
goreleaser release
```

to do a dry run (no upload): 
```
goreleaser release --snapshot --rm-dist
```

`GITHUB_TOKEN` is a standard PAT from github and needs to be set to upload.