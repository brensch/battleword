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
		"definition": {
			"id": "a03753b1-c30c-4bad-8b58-201539980d44",
			"name": "solvo",
			"description": "cool guy"
		},
		"state": [{
			"game_id": "b52ccd3a-6344-4e14-bd2c-87d5502fd128",
			"guesses": ["mouse", "riser", "brown", "polka", "quell", "panel"],
			"results": [
				[0, 2, 0, 2, 0],
				[1, 0, 1, 0, 0],
				[0, 1, 1, 1, 0],
				[0, 2, 0, 0, 0],
				[0, 0, 0, 0, 0],
				[0, 0, 0, 0, 0]
			],
			"shouts": ["what's the point of anything?", "wordle is fun, but for how long?", "what's the point of anything?", "what's the point of anything?", "wordle is fun, but for how long?", "wordle is fun, but for how long?"],
			"times": [618300, 258000, 302800, 292100, 198500, 206600],
			"total_time": 1876300
		}]
	}, {
		"definition": {
			"id": "e64b7af2-6785-4923-b56e-7ab1b7622ed2",
			"name": "bolvo",
			"description": "cool guy"
		},
		"state": [{
			"game_id": "b52ccd3a-6344-4e14-bd2c-87d5502fd128",
			"guesses": ["crick", "endow", "stomp", "grass", "shove", "flora"],
			"results": [
				[0, 1, 0, 0, 0],
				[0, 0, 0, 1, 1],
				[1, 1, 1, 0, 0],
				[0, 1, 0, 2, 0],
				[1, 0, 1, 0, 0],
				[0, 0, 1, 1, 0]
			],
			"shouts": ["there has to be a better strat than this", "what's the point of anything?", "you will one day be dust, but i will always be solvo", "wordle is fun, but for how long?", "there has to be a better strat than this", "you will one day be dust, but i will always be solvo"],
			"times": [647700, 307700, 277200, 344200, 203100, 177000],
			"total_time": 1956900
		}]
	}],
	"games": [{
		"id": "b52ccd3a-6344-4e14-bd2c-87d5502fd128",
		"answer": "worst",
		"result": {
			"start": "2022-02-15T21:51:07.8650467-08:00",
			"end": "2022-02-15T21:51:07.8670158-08:00",
			"fastest": {
				"player": {
					"id": "a03753b1-c30c-4bad-8b58-201539980d44",
					"name": "solvo",
					"description": "cool guy"
				},
				"time": 1876300
			},
			"most_accurate": {
				"player": {
					"id": "a03753b1-c30c-4bad-8b58-201539980d44",
					"name": "solvo",
					"description": "cool guy"
				},
				"average_guess_length": 6
			},
			"loudest": {
				"player": {}
			}
		}
	}]
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