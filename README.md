# Chess engine

This is a proof of concept implementation of a chess engine developed while writing my Bachelor's thesis. It implements from scratch all aspects of a basic chess engine: legal moves generation, tree exploration and static evaluation.

The performance of the tree exploration is acceptable and the main focus should now be implementing better heuristics for the static evaluation.

## How to use

Launch the backend server

```
cd server
go run .
```

Launch the frontend UI

```
cd ui
npm run start
```

Open the [localhost:3000](http://localhost:3000)
