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

### /ping
in order to get the definition of your character, the engine will ping you. this is also run at the start of each match up to 10 times in order to wake up your server if you're hosting it in serverless land where everything is slightly less reliable.

#### Request:
GET request - no payload
#### Response:
```json
{
	"name": "solvo",
	"description": "the magnificent"
}
```
there will be more things here in the future. stay posted.

### /results
once all players are finished, the engine will send you the results of everyone in the match so you can brag. no response is required.
#### Request:
```json
{
	"players": [{
		"id": "5dd5148f-19fe-4b00-b248-774cece1f196",
		"definition": {
			"name": "schwordler",
			"description": "the brave"
		},
		"state": [{
			"game_id": "896a2327-f772-40cd-b128-11aafb72c93e",
			"guesses": ["cigar", "humph", "stool", "belly", "dwelt", "knelt"],
			"results": [
				[0, 0, 0, 0, 0],
				[0, 0, 0, 0, 0],
				[0, 1, 0, 0, 1],
				[0, 1, 0, 2, 0],
				[0, 0, 2, 2, 2],
				[2, 2, 2, 2, 2]
			],
			"correct": true,
			"times": [726200, 504000, 1278800, 746900, 739300, 739300],
			"total_time": 4734500
		}, {
			"game_id": "60a58c96-8676-4d6b-a4fa-78d46f87f8d4",
			"guesses": ["cigar", "humph", "awake", "abyss", "aloft", "asset"],
			"results": [
				[0, 0, 0, 1, 0],
				[0, 0, 0, 0, 0],
				[2, 0, 0, 0, 1],
				[2, 0, 0, 1, 1],
				[2, 0, 0, 0, 2],
				[2, 2, 2, 2, 2]
			],
			"correct": true,
			"times": [763300, 929300, 1009600, 747500, 839200, 871000],
			"total_time": 5159900
		}, {
			"game_id": "920a3e78-b501-425f-8824-1a238c08cbbf",
			"guesses": ["cigar", "humph", "hasty", "hazel"],
			"results": [
				[0, 0, 0, 1, 0],
				[2, 0, 0, 0, 0],
				[2, 2, 0, 0, 0],
				[2, 2, 2, 2, 2]
			],
			"correct": true,
			"times": [472000, 679100, 919100, 1068700],
			"total_time": 3138900
		}],
		"player_summary": {
			"total_time": 13033300,
			"total_guesses": 16,
			"average_guesses": 5.333333333333333,
			"games_won": 3,
			"disqualified": false
		}
	}, {
		"id": "c0024473-f690-48f2-b7a5-8940b95bc5f7",
		"definition": {
			"name": "solvo",
			"description": "the magnificent"
		},
		"state": [{
			"game_id": "60a58c96-8676-4d6b-a4fa-78d46f87f8d4",
			"guesses": ["bleed", "roger", "murky", "serif", "spill", "sinew"],
			"results": [
				[0, 0, 0, 2, 0],
				[0, 0, 0, 2, 0],
				[0, 0, 0, 0, 0],
				[1, 1, 0, 0, 0],
				[1, 0, 0, 0, 0],
				[1, 0, 0, 2, 0]
			],
			"correct": false,
			"times": [868300, 562800, 313200, 404300, 358700, 286300],
			"total_time": 2793600
		}, {
			"game_id": "896a2327-f772-40cd-b128-11aafb72c93e",
			"guesses": ["murky", "trial", "preen", "lasso", "piper", "blend"],
			"results": [
				[0, 0, 0, 1, 0],
				[1, 0, 0, 0, 1],
				[0, 0, 2, 0, 1],
				[1, 0, 0, 0, 0],
				[0, 0, 0, 1, 0],
				[0, 1, 2, 1, 0]
			],
			"correct": false,
			"times": [814800, 299200, 456700, 514200, 242400, 268100],
			"total_time": 2595400
		}, {
			"game_id": "920a3e78-b501-425f-8824-1a238c08cbbf",
			"guesses": ["known", "drown", "ficus", "flash", "built", "beret"],
			"results": [
				[0, 0, 0, 0, 0],
				[0, 0, 0, 0, 0],
				[0, 0, 0, 0, 0],
				[0, 1, 1, 0, 1],
				[0, 0, 0, 1, 0],
				[0, 0, 0, 2, 0]
			],
			"correct": false,
			"times": [836300, 807800, 291900, 210200, 380300, 172800],
			"total_time": 2699300
		}],
		"player_summary": {
			"total_time": 8088300,
			"total_guesses": 21,
			"average_guesses": 7,
			"games_won": 0,
			"disqualified": false
		}
	}],
	"games": [{
		"id": "896a2327-f772-40cd-b128-11aafb72c93e",
		"answer": "knelt",
		"summary": {
			"start": "2022-02-18T01:56:05.2105036-08:00",
			"end": "2022-02-18T01:56:05.2152943-08:00",
			"fastest": {
				"player_id": "c0024473-f690-48f2-b7a5-8940b95bc5f7",
				"time": 2595400
			},
			"most_accurate": {
				"player_id": "5dd5148f-19fe-4b00-b248-774cece1f196",
				"average_guess_length": 6
			},
			"loudest": {}
		}
	}, {
		"id": "920a3e78-b501-425f-8824-1a238c08cbbf",
		"answer": "hazel",
		"summary": {
			"start": "2022-02-18T01:56:05.2105308-08:00",
			"end": "2022-02-18T01:56:05.2137058-08:00",
			"fastest": {
				"player_id": "c0024473-f690-48f2-b7a5-8940b95bc5f7",
				"time": 2699300
			},
			"most_accurate": {
				"player_id": "5dd5148f-19fe-4b00-b248-774cece1f196",
				"average_guess_length": 4
			},
			"loudest": {}
		}
	}, {
		"id": "60a58c96-8676-4d6b-a4fa-78d46f87f8d4",
		"answer": "asset",
		"summary": {
			"start": "2022-02-18T01:56:05.2104935-08:00",
			"end": "2022-02-18T01:56:05.2157021-08:00",
			"fastest": {
				"player_id": "c0024473-f690-48f2-b7a5-8940b95bc5f7",
				"time": 2793600
			},
			"most_accurate": {
				"player_id": "5dd5148f-19fe-4b00-b248-774cece1f196",
				"average_guess_length": 6
			},
			"loudest": {}
		}
	}],
	"summary": {
		"fastest": {
			"player_id": "c0024473-f690-48f2-b7a5-8940b95bc5f7",
			"time": 8088300
		},
		"most_accurate": {
			"player_id": "5dd5148f-19fe-4b00-b248-774cece1f196",
			"average_guess_length": 5.333333333333333
		},
		"loudest": {},
		"most_correct": {
			"player_id": "5dd5148f-19fe-4b00-b248-774cece1f196",
			"correct_games": 3
		},
		"games_fastest": {
			"player_id": "c0024473-f690-48f2-b7a5-8940b95bc5f7",
			"count": 3
		},
		"games_loudest": {},
		"games_most_accurate": {
			"player_id": "5dd5148f-19fe-4b00-b248-774cece1f196",
			"count": 3
		}
	}
}
```

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