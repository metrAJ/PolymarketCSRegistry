# Polymarket CS2 Active markets registry
Back end server to scrap and compose active markets related to CS2 team matches. Plus cli and tui apps to get the feed. 
## Tech Stack 
* **Language:** GO 1.25.5
* **PKG:** GAMMA API 
* **TUI:** tview
### Prerequisites
* Git installed
* GO 1.23.0+ 
## API Endpoints
The API runs on `http://localhost:3000` by default.
### 1. Update events
* URL: `http://localhost:3000/api/scrape`
* Method: POST
Forcefull udate events and markets data. Server automatically udating every 15 sec by default. 
### 1. Get events
* URL: `http://localhost:3000/api/events`
* Method: GET
Get all active events and markets in JSON format. 
## Quck use
Project includes make commands to easely use the app. 
1. Start the server: 
```bash
make run-server
```
2. Start cli:
```bash
make run-cli
```
3. Start tui: 
```bash
make run-tui
```
## Console client
Project includes basic CLI and TUI apps for navigation throgh events and markets. 
### CLI 
Shows 1 event at a time. Change events by pressing ENTER.
![alt text](image.png)
## TUI 
Shows all events as a navigation list and their markets and the right window. Navigate with arrows and TAB ot by mouse. 
![alt text](image-1.png)