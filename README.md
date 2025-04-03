***Calorific*** is a single page, local webapp written entirely in Go, using Templ templates for its frontend. 

This project was built mostly because I wanted minimalist calorie tracking app and I could not find a decent free app that I liked.

Since it was something I wanted to build, I figured I might as well use it as an introductory project to help learn Go, which I had been watching lectures about for about a month previous to starting. As a result, *a lot* of the project structure could use some improvement (especially the templ side of things, I used internal CSS components in each templ file and used templ more as a way to embed go into HTML -- hopefully I can learn from this in the future). Since I didn't plan on hosting it anywhere, I used sqlite for my db, specifically using the driver found at github.com/mattn/go-sqlite3 alongside the stdlib database/sql package.


Anyways, in order to run this project for yourself if you wanted to try it out for some reason, after cloning simply:

* Install Go version 1.23.4 
* Ensure dependencies needed from within the go mod file are installed
* *(Recommended)* Delete calorific.db file
* Execute make in the project directory
* Navigate to *localhost:8082*

