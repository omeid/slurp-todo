# Golang Angular Todo
Simple obligatory Todo app showcasing [Slurp](https://github.com/omeid/slurp).

> Heads up! You need a working GOPATH for this.


- Install Slurp.
  ```bash
  $ go get -v   github.com/omeid/slurp           # Get slurp.
  $ go install  github.com/omeid/slurp/cmd/slurp # Install runner.
  ```

- Clone this repo and then get it, or just clone this one.
  ```bash
  $ go get -v  github.com/omeid/slurp-todo
  ```

- Run 
  ```bash
  $ cd $GOPATH/github.com/YOURUSERNAME/slurp-todo
  $ slurp init # Download deps, may take a little time.
  $ slurp 
  ```

Now point your browser at `http://localhost:8081` and voila!

Now if you edit any of the frontend files, they will be livereloaded\*
and when you change .go files, the server will be restarted.


\*Works best in Firefox.
