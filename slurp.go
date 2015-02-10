// +build slurp

package main

import (
	"github.com/omeid/slurp"
	"github.com/omeid/slurp/stages/archive"
	"github.com/omeid/slurp/stages/fs"
	"github.com/omeid/slurp/stages/resources"
	"github.com/omeid/slurp/stages/web"

	"github.com/slurp-contrib/ace"
	"github.com/slurp-contrib/gcss"
	"github.com/slurp-contrib/gin"
	"github.com/slurp-contrib/jsmin"
	"github.com/slurp-contrib/livereload"
	"github.com/slurp-contrib/watch"
)

func init() {
	config.Livereload = ":35729"
}

func Slurp(b *slurp.Build) {

	// Download deps.
	b.Task("libs", nil, func(c *slurp.C) error {
		return web.Get(c,
			//TODO: this can use some lovin'. Perhaps a go-bower?
			"https://github.com/angular/bower-angular/archive/v1.3.12.zip",
			"https://github.com/angular/bower-angular-route/archive/v1.3.12.zip",
			"https://github.com/angular/bower-angular-resource/archive/v1.3.12.zip",
		).Then(
			archive.Unzip(c),
			fs.Dest(c, "libs/"),
		)

	})

	b.Task("libs.js", nil, func(c *slurp.C) error {
		return fs.Src(c,
			"libs/*angular*/angular.min.js",
			"libs/bower-angular-route-*/angular-route.min.js",
			"libs/bower-angular-resource-*/angular-resource.min.js",
		).Then(
			slurp.Concat(c, "libs.js"),
			fs.Dest(c, "./public/assets/"),
		)
	})

	b.Task("gcss", nil, func(c *slurp.C) error {
		return fs.Src(c, "frontend/*.gcss").Then(
			gcss.Compile(c),
			slurp.Concat(c, "style.css"),
			fs.Dest(c, "./public/assets/"),
		)
	})

	b.Task("js", nil, func(c *slurp.C) error {
		return fs.Src(c,
			"frontend/*.js",
		).Then(
			slurp.Concat(c, "app.js"),
			jsmin.JSMin(c),
			fs.Dest(c, "./public/assets/"),
		)
	})

	b.Task("ace", nil, func(c *slurp.C) error {
		return fs.Src(c,
			"frontend/*.ace",
		).Then(
			ace.Compile(c, ace.Options{
				//Because we use {{ and }} for angular.js
				DelimLeft:  "<<",
				DelimRight: ">>",
			}, config),
			fs.Dest(c, "./public"),
		)
	})

	b.Task("gin", nil, func(c *slurp.C) error {
		gin := watch.Watch(c, gin.Gin(c, &gin.Config{}, "-tags=slurp"), "*.go", "*/*.go", "*/*/*.go")

		b.Defer(func() { gin.Close() })

		b.Wait() //Wait forever.
		return nil
	})

	b.Task("frontend", []string{"libs.js", "js", "ace", "gcss"}, func(c *slurp.C) error {
		return nil
	})

	b.Task("watch", []string{"frontend"}, func(c *slurp.C) error {

		g := watch.Watch(c, func(string) { b.Run(c, "gcss") }, "frontend/*.gcss")
		a := watch.Watch(c, func(string) { b.Run(c, "ace") }, "frontend/*.ace")
		j := watch.Watch(c, func(string) { b.Run(c, "js") }, "frontend/*.js")

		//Close all the watchers on exit.
		b.Defer(func() {
			g.Close()
			a.Close()
			j.Close()
		})
		b.Wait() //Wait forever.
		return nil
	})

	b.Task("embed", nil, func(c *slurp.C) error {
		return fs.Src(c,
			"public/*",
			"public/*/*",
		).Then(
			resources.Stage(c, resources.Config{
				Pkg:     "main",
				Var:     "Public",
				Declare: false,
			}),
			fs.Dest(c, "."),
		  )
	})

	b.Task("livereload", nil, func(c *slurp.C) error {

		l := watch.Watch(c, livereload.Start(c, config.Livereload, "public"),
			"public/*",
			"public/assets/*",
		)

		b.Defer(func() {
			l.Close()
		})

		b.Wait() //Wait forever.
		return nil
	})

	// # Special tasks
	// when running this task with "slurp" it will run `go get`
	// for build dependenceis.
	b.Task("init", []string{"libs"}, func(c *slurp.C) error {
		//ideal for checking deps.
		return nil
	})

	//When running slurp with no args, well, the "default" task is run.
	b.Task("default", []string{"livereload", "watch", "gin"}, func(c *slurp.C) error {
		//ideal for clean up.
		return nil
	})
}
