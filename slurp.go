// +build slurp

package main

//This file will be only compiled along the project with slurp. So don't put any projec code here.

import (
	"github.com/omeid/slurp"
	"github.com/omeid/slurp/stages/archive"
	"github.com/omeid/slurp/stages/fs"
	"github.com/omeid/slurp/stages/util"
	"github.com/omeid/slurp/stages/web"

	"github.com/slurp-contrib/ace"
	"github.com/slurp-contrib/gcss"
	"github.com/slurp-contrib/gin"
	"github.com/slurp-contrib/jsmin"
	"github.com/slurp-contrib/livereload"
	"github.com/slurp-contrib/resources"
	"github.com/slurp-contrib/watch"
)

func init() {
	config.Livereload = ":35729"
}

// This function is called to allow registering the tasks when slurp is run.
func Slurp(b *slurp.Build) {

	// Download deps, why use a package manger for simply downloading some links and unziping. :)
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

	//This stage concatenates the javascript libraries files.
	b.Task("libs.js", nil, func(c *slurp.C) error {
		return fs.Src(c,
			"libs/*angular*/angular.min.js",
			"libs/bower-angular-route-*/angular-route.min.js",
			"libs/bower-angular-resource-*/angular-resource.min.js",
		).Then(
			util.Concat(c, "libs.js"),
			fs.Dest(c, "./public/assets/"),
		)
	})

	// Compiles the gcss.
	b.Task("gcss", nil, func(c *slurp.C) error {
		return fs.Src(c, "frontend/*.gcss").Then(
			gcss.Compile(c),
			util.Concat(c, "style.css"),
			fs.Dest(c, "./public/assets/"),
		)
	})

	//Minfiy our javascript.
	b.Task("js", nil, func(c *slurp.C) error {
		return fs.Src(c,
			"frontend/*.js",
		).Then(
			util.Concat(c, "app.js"),
			jsmin.JSMin(c),
			fs.Dest(c, "./public/assets/"),
		)
	})

	// Compile the ace templates.
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

	// This will fire-up a gin build server and proxy, it will rebuld the app everytime a go file changes.
	// Uses the slurp tag to allow for package configuration (see the init() func above).
	b.Task("gin", nil, func(c *slurp.C) error {
		gin := gin.NewGin(c, &gin.Config{}, "-tags=slurp")
		watch := watch.Watch(c, gin.Run, "*.go", "*/*.go", "*/*/*.go")

		<-c.Done()
		watch.Close()
		gin.Close()
		return nil
	})

	//Frontend requires the libs.js, js, ace, and gcss tasks, this is basically "grouping" tasks.
	b.Task("frontend", []string{"libs.js", "js", "ace", "gcss"}, func(c *slurp.C) error {
		return nil
	})

	//The name says a lonet.
	b.Task("watch", []string{"frontend"}, func(c *slurp.C) error {

		g := watch.Watch(c, func(string) { b.Run(c, "gcss") }, "frontend/*.gcss")
		a := watch.Watch(c, func(string) { b.Run(c, "ace") }, "frontend/*.ace")
		j := watch.Watch(c, func(string) { b.Run(c, "js") }, "frontend/*.js")

		<-c.Done()
		g.Close()
		a.Close()
		j.Close()
		return nil
	})

	//This will generate the resource file.
	b.Task("embed", nil, func(c *slurp.C) error {
		return fs.Src(c,
			"public/*",
			"public/*/*",
		).Then(
			resources.Build(c, resources.Config{
				Pkg:     "main",
				Var:     "Public",
				Declare: false,
				Tag:     "embed",
			}),
			fs.Dest(c, "."),
		)
	})

	//Start a livereload server and triggered everytime anything in public folder changes.
	b.Task("livereload", nil, func(c *slurp.C) error {

		l := watch.Watch(c, livereload.Start(c, config.Livereload, "public"),
			"public/*",
			"public/assets/*",
		)

		<-c.Done()
		l.Close()

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
