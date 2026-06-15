# chess tower defense game
basically a tower defense game but with chess pieces
pieces attack in the way that they would normally attack in chess (bishop diagonally, knight L etc)
their locations are reset at the start of each round to the beginnign positions (unless you have a modified otherwise)
enemy path is randomised

its also kind of a roguelike system where there are modifiers that you can use like permanently deleting a piece (useful when pawns are in the way) or changing the dps or making a piece always start the round at a place that you want it to start it.

so Bloons TD/Chess/Balatro/Slay the Spire

its really cool

idk what else to put in this readme

```go run main.go```



# Notes
```cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .```
```GOOS=js GOARCH=wasm go build -o ./out/wasm/game.wasm```

# Todo
- make path gen more aware of where things are and do obstacle gen
- balancing improvements
- card system
- validate piece moves
- make look good
- if you double click on a piece it deletes - why?