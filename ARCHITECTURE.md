# Pacman Architecture
[Pacman video - youtube]()

Pac-Man game were implemented using multithreaded programming using the GO language. You can run Ghost Enemy and Pac-Man separately, here you have the arquitecture of the game ilustratedin a uml diagram

## Diagram architecture
(diagram.jpeg)

## Pacman:

This component is for the player, with this component we track the place of pacman and it is stacked into the labyrinth simultaneously the labyrinth is stacked; we handle its developments with the bolt keys presses with a strategy in the know of the game.

the following keys are used to move the pacman: up,down,left,right

Pacman can have the amount of lives that the user declare at the beginning of the game.

## Ghosts:

The enemies elements works each one in a thread and we generation the movement of the ghost with a random number from and each time the value of this number is one of the following, it will make the correspondes movement, this is a very basic way to make the ghost moves by itself:
		
In the same way as pacman, the ghost handle the collision with the walls, and it acts as a cancelled movement.

	
```go
func fantasmaC() {
	for _, g := range enemy {
		dir := fantasmaM()
		g.row, g.col = moving(g.row, g.col, dir)
		time.Sleep(800 * time.Millisecond)
	}

}
```

## Maze config

The maze for the game is loaded from the *level1.txt* file in the repository and print it in a infinite loop line by line in the console, this file has the necesary information to build the maze in our program an handle the walls and the posible postions for the enemy and pacman.
The maze is cleaned and printed each loop to handle the correct visualization as a game when a ghost or pacman moved. # is used to represent a wall in the file.


Placing pacman, ghosts and candies in level 1.

    -P --> Placing player in the maze 
    -G --> Placing Ghost in the maze 
    -. --> Placing Candies

```go
var level []string
```

**Maze Creation pacing players and gohst:**
```go
 func levelL(file string) error {
	f, badLevel := os.Open(file)
	if badLevel != nil {
		return badLevel
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		level = append(level, line)
	}

	for row, line := range level {
		for col, char := range line {
			switch char {
			case 'P':
				pacman = sprite{row, col}
			case 'G':
				enemy = append(enemy, &sprite{row, col})
			case '.':
				puntos++
			}
		}
	}

	return nil
}
```


	
	
## Jason Configuration

I use a jason file for setting the principal sprites for the game and they can be changed for any icon that the user wants in the code changing the *config.json* file.
	

Functions
-------------

`dibujar()` this prints a new frame of the game

`limpieza():` limpia la consola 

`moving(row, col int)` Moves the cursor to the given coordinates in order to print a change in the board.

`pacmanM(row, col int)` moves paccman to the required field 

`cargarLevel():` Reads information from Json File the configuration of the sprites.

`controles():` in a goroutine the program Reads from standard input and returns the direction for the movement 

`fantasmaM() string` generates de random movement of the ghost by creating a random number 

`fantasmaC()` in a goRoutinee I move the ghost to a corresponding place 

`levelL(file string)` Reads the designed level from the .txt file. 




## How to Run the program 

```go
{
	//place yourself in the main directory were all files are saved and run de following comand:
    go run pacman.go

    //to avoid errors open a big terminal window
}
```
### or use
```go
{
	make test
}
```