### todo

* I want that the file operations have fmt.style %s inputs, idea: merge src and dst string together with | separator, inject the values (fail is number of items is wrong), then split and pass for processing
* change taskrunner so that it is more functional, like: tr.Cmd().Envs("PROFILE", TEST", "SOME_ENV", SOME_VALUE").Dir(srcDir).AsDaemon().Execute("asd %s asd", someTermToInject)
* still show output when "hiding cleanup" is enabled
* idea for making outputs clearer: when running a daemon it gets a process name. we can print output of multiple daemon processes at once in console; each process has a different output prefix and a different color; this way we know what happens at a specific time and having nice coloring and prefixes for distinguishing the daemons
