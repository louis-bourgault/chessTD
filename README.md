# chess tower defense game

game design

at the start the pieces are reset back to their set placed

Chess based roguelike tower defense game. There is a standard 8x8 board that resets at the start of each round (back to either default chess or chess960)

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