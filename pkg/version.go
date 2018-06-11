package lexmachine

// Must be set at build via
// -ldflags "-X lexmachine.Version=`cat VERSION`"
// -ldflags "-X lexmachine.Version=`git describe --tags`"
var Version = "0.13.2-dev"
