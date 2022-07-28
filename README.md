# Directory-Organizer
Organizes the files in a directory. Sorts files into folders by extension.

## About
Written using go1.18.

## Getting Started
- Place test files in ./data/test/input
- `go get .`
- `go build`
- `./Directory-Organizer --help` - make sure you set the output directory to ./data/test/output

## Example Usage
You could run this everyday at 2am via a cron job to organize and cleanup your downloads folder.
For example, add `0 2 * * * ~/MyPrograms/Directory-Organizer -t ~/Downloads -o ~/CleanDownloads -s` via `crontab -e`.

## TODOs
- Add recursive flag
- Add different strategies, not just by extension
- Add no extension folder name flag
