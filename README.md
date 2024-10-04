# GitHub Activity CLI

`githubactivity` is a simple command-line interface (CLI) tool that fetches and displays the recent public events for a specified GitHub user. This is a lightweight utility to quickly check a user's GitHub activity, such as push events, pull requests, issue comments, and more.

## Features

- Fetches recent public events from GitHub for a specified user.
- Displays events such as pushes, pull requests, issues, stars, and forks.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/githubactivity-cli.git
   cd githubactivity-cli
   ```
2. **Build the CLI**
   Assuming you have Go installed, build the project with the following command:
   ```bash
   go build -o githubactivity

## Usage
After building the executable, the CLI can be invoked like this:
```bash
githubactivity <username>
```

Example:
```bash
githubactivity JesterSe7en
```
This will print the list of recent events for the user "JesterSe7en" on GitHub.

**Ouput Example:**
```
Found 3 events
Stared to watch repository @ Fri, Jul 26, 2024 at 11:32 AM: shauninman/MinUI
Stared to watch repository @ Fri, Jul 26, 2024 at 11:32 AM: OnionUI/Onion
Stared to watch repository @ Sat, Jul 13, 2024 at 8:43 PM: raizam/gamedev_libraries
```

## Dependencies
This CLI tool uses the following:
- Go (version 1.23.2): developed on that version, untested on others

## License
This project is licensed under the MIT License. See the LICENSE file for details.
