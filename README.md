# Checkmatey ♟️
Checkmatey is a web chess application that allows users to play against an AI in the browser. The backend is implemented in Go and the frontend in React JS.

![Screenshot](/assets/chess_board.png)


## Installation
### Prerequisites
- Golang: [Install Golang](https://golang.org/dl/)
- Node: [Install Node](https://nodejs.org/en/download)

### Steps
Clone the repository:
```
git clone git@github.com:mark-b96/checkmatey.git
cd checkmatey
```
Install backend dependencies:
```
cd backend
go get ./...
```
Install dependencies for the frontend:
```
cd frontend
npm install
```

## Usage

```
# In the root directory, give your user permission to run the script
chmod +x run.sh
./run.sh
```

## Tests
```
cd backend
go test ./...
```


## To do
- Extend test coverage
- Implement AI algorithm
- Dockerise
