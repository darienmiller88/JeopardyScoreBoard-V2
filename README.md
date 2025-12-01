# Jeopardy ScoreBoard V2

This is version 2.0 of the first sscoreboard I did started last year. While the first version i made achieved all of its goals (replacing score keeping by paper, saving games, allowing of game lengthing, etc), I am not satisified with the quality of the code I produced. Initally starting with HTMX, I transitioned to Vue for the project, and while it definitely served me well, the end result was a lot more code than I was hoping for, and a tricky architecture to follow. I will be simplyifing the tech stack with the end result being a scoreboard that is much easier to use and modify.

Here are the following changes I plan on making:

- For the front end, I will replace Vue with HTMX, Alpine.js and SCSS. On paper, this means I will sacrifice heavy front end tooling for a massively simplified workflow that aligns well the simple styling of the score keeper. Alpine will replace front end state management and non-server updates.

- For the back end, I will stick with go, but instead of it serving JSON to seperate front end, it will instead serve HTML fragments and HTMX pages instead. The HTMX will make request to end points on the server, receive fragments, and place them onto the HTMX page directly. Game state will be saved in a redis cache rather than in local storage to ensure the names of players aren't exposed, and that point incrementing and decrementing are near instant.

### New Tech Stack:
* [HTMX](https://vuejs.org/)
* [Sass](https://sass-lang.com/)
* [Alpine.js](https://alpinejs.dev/)
* [Redis](https://redis.io/)
* [Go](https://github.com/go-chi/chi)
* [MongoDB-Atlas](https://www.mongodb.com/cloud/atlas)
* [Railway](https://railway.com)

## Features
This version will still retain the following features from v1:

* Single player Jeopardy. Client can add and remove Users to a specific ADAPT location, and add and subtract points to them.
* Tournament/Team jeopardy mode where the client can choose the host location for the game, and which teams are playing.
* The ability to save games, both single player and team jeopardy.
* Viewing games, which will include the winner of the game, and the total and average amount of points earned during the game.

### Requirements:

* Clone repo using `git clone https://github.com/darienmiller88/JeopardyScoreBoard-V2`
* Migrate the necessary information to your local `.env` as described in the `.env_sample` file
* Run go build to create a root level `JeopardyScoreBoardV2.exe` file, and then run `.\JeopardyScoreBoard-V2` to run the executable. If an executable is not needed, simply input `go run main.go` instead, or `.\fresh` to enable a server restart on change.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Feel free to leave suggestions as well, I'm always looking for ways to improve!

<p align="right">(<a href="#top">back to top</a>)</p>

## License
[MIT](https://choosealicense.com/licenses/mit/)