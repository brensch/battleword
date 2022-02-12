# battleword
wordle is cool right now

## what is this
battleword is a competition to see who can come up with the fastest/most accurate/shoutiest wordle solver.

players host an api, then the battleword engine will make a `POST` request to their api with the state of a wordle (starting empty)

the player's api should then respond to the state of the game in the body of the post with their best guess. as soon as the battleword engine hears back from them, it will send the results of their guess in a new request. it will do that until the player's api guesses correctly, or they reach the guess limit.

## quickstart

1. download the latest release for your OS and unpack
2. run `solvo` (double click) - this starts solvo the solver. he will listen for game states from engine.
3. run `engine` - this starts sending game states to solvo. with every guess solvo makes, engine will send a new request to solvo with the results of his previous guess. solvo will ignore those results and choose a completely random word to send next. your solver should do better than solvo.

## setup
to test your own guesser against the engine, create an api that implements the schema below. once you've done that, run the engine against the api location of your solver like so:

```
./engine --names muchbettersolver --apis http://localhost:8081
```

you can specify multiple solvers to compete against each other:
```
./engine --names muchbettersolver,solvo --apis http://localhost:8081,http://localhost:8080
```

NB these commands are executed in a command line of your choice. Exact syntax may change based on your OS.


## api
this is what all solvers need to implement.
### /guess
the engine will hit your api here with the previous results of a game. you are expected to respond with your best guess.
#### Request:
```json
{
	"guesses": ["beast", "lapse"],
	"results": [
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
		"state": {
			"guesses": ["beast", "lapse", "pause"],
			"results": [
				[0, 1, 1, 2, 0],
				[0, 2, 1, 2, 2],
				[2, 2, 2, 2, 2]
			],
			"times": [1012000, 336700, 473626]
		}
	}, {
		"name": "skye",
		"state": {
			"guesses": ["beast", "beast", "beast", "beast", "beast", "beast"],
			"results": [
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