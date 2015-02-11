# Golang Angular Todo
Simple obligatory Todo app showcasing [Slurp](https://github.com/omeid/slurp).


> Heads up! You need a working GOPATH for this.


- Install Slurp.
  ```bash
  $ go get -v   github.com/omeid/slurp/cmd/slurp # Get slurp.
  ```

- Clone this repo and then get it, or just clone this one.
  ```bash
  $ go get -v  github.com/omeid/slurp-todo
  ```

- Try and live code!
  ```bash
  $ cd $GOPATH/github.com/YOURUSERNAME/slurp-todo
  $ slurp init             # Download deps
  $ slurp                  # Start default task: gin, watch, livereload. 
  ```

  Now point your browser at `http://localhost:8081` and voila!

  You're now in development, if you edit any of the frontend files, they will be livereloaded\* and when you change .go files  the server  will be restarted.


- Build.
  Once you're happy with the product, build it.

  ```bash
  $ slurp frontend        # Rebuild the frontend, just to be sure.
  $ go build              # Build it!
  ```

  Now you should have a `slurp-todo` binary that will work, it will server the "assets" from the public folder.

  We can make things even more simpler.

- Embed!
  
  ```bash
  $ slurp embed           # Generate the public folder as public_resource.go
  $ go build -tags=embed  # Build the app with public_resource.go
  ```

  Now all you need to run this app is the `slurp-todo`, you don't need to ship the public folder.
  
  But why stop there?

- Go nuts!
  ```bash
    $ go build -a -tags 'netgo embed' --ldflags '-extldflags "-static"'
  ```
    Now you have a static `slurp-todo` that doesn't require any libraries.\*\*


#### TODO

 - [x] Embed resorces (public directory).
 - [ ] Use a database file instead of memory.
 - [ ] Replace SQLite3 with something that is properly static.


##### Ace? gcss?
The frontend is written using [Ace](https://github.com/yosssi/ace) for html templates and [gcss](https://github.com/yosssi/gcss) as css preprocessor, please refer to their documentation for their syntax documentation and more details

\*  Works best in Firefox.  
\*\* This is not entirely true as the app is using SQLite3 and it depends on `dlopen` calls which requires a `glibc`, this can be easily avoided by replacing the SQLite3 database with a pure Go database or something that can be properly staticly linked. 
