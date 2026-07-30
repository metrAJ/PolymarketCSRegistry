# Polymarket CS2 Active markets registry
Back-end server to scrape and compose active markets related to CS2 team matches. Plus CLI and TUI apps to get the feed. 
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
Forceful update events and markets data. The server automatically updates every 15 sec by default. 
### 1. Get events
* URL: `http://localhost:3000/api/events`
* Method: GET
Get all active events and markets in JSON format. 
## Quck use
The project includes make commands to easily use the app. 
1. Start the server: 
```bash
make run-server
```
2. Start the CLI:
```bash
make run-cli
```
3. Start the TUI: 
```bash
make run-tui
```
## Console client
The project includes basic CLI and TUI apps for navigation through events and markets. 
### CLI 
Shows 1 event at a time. Change events by pressing ENTER.
<img width="742" height="598" alt="image" src="https://github.com/user-attachments/assets/f34b422f-d8fa-485e-98e9-c3193b1bbfe7" />

## TUI 
Shows all events as a navigation list and their markets in the right window. Navigate with arrows and TAB ot by mouse. 
<img width="1379" height="379" alt="image" src="https://github.com/user-attachments/assets/7ecdc174-e47f-452b-b86f-a3541ca16a00" />
