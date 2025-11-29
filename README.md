# Jeopardy ScoreBoard V2

This is version 2.0 of the first sscoreboard I did started last year. While the first version i made achieved all of its goals (replacing score keeping by paper, saving games, allowing of game lengthing, etc), I am not satisified with the quality of the code I produced. Initally starting with HTMX, I transitioned to Vue for the project, and while it definitely served me well, the end result was a lot more code than I was hoping for, and a tricky architecture to follow. I will be simplyifing the tech stack with the end result being a scoreboard that is much easier to use and modify.

Here are the following changes I plan on making:

- For the front end, I will replace Vue with HTMX, Alpine.js and SCSS. On paper, this means I will sacrifice heavy front end tooling for a massively simplified workflow that aligns well the simple styling of the score keeper. Alpine will replace front end state management and non-server updates.

- For the back end, I will stick with go, but instead of it serving JSON to seperate front end, it will instead serve HTML fragments and HTMX pages instead. The HTMX will make request to end points on the server, receive fragments, and place them onto the HTMX page directly. Game state will be saved in a redis cache rather than in local storage to ensure the names of players aren't exposed, and that point incrementing and decrementing are near instant.