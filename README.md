# Loot
Command Storage TUI App

![image](https://github.com/user-attachments/assets/99685a2a-3110-447e-a8c5-a14c6a774a95)

## Why
I mainly wanted to try out the Charm TUI toolkits, but I noticed myself running a command, then forgetting it. Yes there are alias apps out there that can manage this, but it is not my own, and I don't really care. 

## How
To use this app you can just clone the repo, and run `go run .` for development or `go install` for usage across your system. Note if you are going to use this, you probably want to update where the log and db files get saved to. 

## Future
I want to obviously make this readme less bad, but also add a few more features to it

- [X] Cleanup (I had no idea what I was doing with Charm initially)
- [ ] Write the DB/Log file to a safe location by default
- [ ] Styling 
- [ ] Add ability to add long form documentation to each command
- [ ] Emulate alias functionality for luls, so `loot <command name>`

## Things Used
BBolt - https://github.com/etcd-io/bbolt
Why use Bolt? Its fast and simple. I didn't need anything more than a key and value. 

Charm - https://github.com/charmbracelet
Why use Charm? Well I like all my apps to look good, building a simple terminal interface would have been fine, but not great. 
