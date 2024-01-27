<!-- <image, logo> -->
<!-- <center title> -->
# Local Diary
<!-- <maybe replace the current description below and make it a small description here> -->
<!-- <badges, version badge> -->

Local Diary is a lightweight and intuitive web application developed in Golang for users who want a local alternative to keep their notes organized. With a user-friendly interface accessible through your browser, Local Diary allows you to effortlessly jot down and save your thoughts, ideas, and important information.


<!-- # Preview -->
<!-- <screenshot > gif of app being used> -->

# Features

* **Local:** Your notes are stored locally, ensuring total privacy and accessibility whenever you need them.
* **Browser Access:** Start the local server, and use your favourite browser to connect to Local Diary. No installation needed!
* **Intuitive Interface:** Simple, clean, fast. Designed to be as straight-forward and minimalistic as possible.

<!-- # Installation
Some prebuilt packages are provided on the [release page of the GitHub project repository] [LINK]. -->

# Development Guide
<!-- <maybe add index>

-Clone Repository:
...
-Makefile Scripts:
--Build and run.......
--run tests...
 -->

## Prerequisites
Ensure you have the following tools and dependencies installed on your system before diving into Local Diary development:
* Golang
* Make GNU
* Air CLI
* GCC or MINGW
* Node and NPM

## Setting Up Local Development

**Clone Repository:**
```bash
git clone github.com/UPSxACE/go-local-diary
cd go-local-diary
```
**Install Dependencies:**
```bash
make dep

# Install playwright browsers (optional) for E2E tests
make dep-browsers
```
**Run parcel watcher for js and css:**
```bash
make parcel-watch
```
**Run watcher for other assets:**
```bash
make assets-watch
```
**Delete the watchers cache:**
```bash
make assets-clear
```
**Run in development mode:**
```bash
# Run with live reload (using air)
air

# Or, as alternative, just run using golang
make dev
```
**Run tests:**
```bash
# Normal test output in console
make test

# Test coverage and output html file
make test-coverage

# Note: Don't forget to initialize the server before trying to run e2e tests
# Test end-to-end with playwright 
make test-e2e
# Show end-to-end report
make test-e2e-report
# Test end-to-end with playwright UI
make test-e2e-ui
```
**Build executables:**
```bash
# Compiles project and outputs 3 executables (windows, linux, mac)
make assets-clear
make parcel-build
make assets-build
make build

# Deletes the executables
make clean
```


<!-- ## Requirements -->

<!-- # Config -->


<!-- Contribute
If you want to contribute check the CONTRIBUTING.md -->

<!-- Donate -->