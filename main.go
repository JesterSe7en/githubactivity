package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	// using github.com/mrz1836/go-sanitize to sanitize the username
	"github.com/mrz1836/go-sanitize"
)

type GithubEvent struct {
	Type       string
	Actor      Actor // who triggered the event
	Repo       Repo  // repository that was affected
	Created_At string
	Payload    json.RawMessage // to be decoded later
}

type Actor struct {
	ID            int
	Login         string
	Display_Login string
	Url           string
}

type Repo struct {
	ID   int
	Name string
	Url  string
}

type CommitCommentEvent struct {
	Action  string
	Comment json.RawMessage
}

type CreateEvent struct {
	Ref           string
	Ref_Type      string
	Master_Branch string
	Description   string
	Pusher_Type   string
}

type DeleteEvent struct {
	Ref      string
	Ref_Type string
}

type ForkEvent struct {
	Forkee struct {
		Created_At string
		HTML_URL   string
		Name       string
	}
}

type Page struct {
	Page_Name string
	Action    string
	Sha       string
	HTML_URL  string
}
type GollumEvent struct {
	Pages []Page
}

type IssueCommentEvent struct {
	Action  string
	Changes json.RawMessage
	Issue   json.RawMessage
	Comment struct {
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads#issue_comment
		Body       string
		Created_At string
		HTML_URL   string
	}
}

type IssuesEvent struct {
	Action string
	Issue  struct {
		Updated_At string
		Title      string
		HTML_URL   string
	}
	Assignee struct {
		Login string
	}
}

type MemberEvent struct {
	Member struct {
		Login string
	}
	Changes struct {
		Role_Name struct {
			To string
		}
	}
}

type PullRequestEvent struct {
	Action       string
	Pull_Request struct {
		Action   string
		Assignee struct {
			Login string
		}
		HTML_URL string
	}
}

type PullRequestReviewEvent struct {
	Action       string
	Pull_Request struct {
		Action   string
		Assignee struct {
			Login string
		}
	}
	Review struct {
		HTML_URL string
	}
}

type PullRequestReviewCommentEvent struct {
	Action       string
	Pull_Request struct {
		Title string
	}
	Comment struct {
		Body       string
		Created_At string
		HTML_URL   string
	}
}

type PullRequestReviewThreadEvent struct {
	Action       string
	Pull_Request struct {
		Title    string
		HTML_URL string
	}
}

type PushEvent struct {
	Size int
	Ref  string
}

type ReleaseEvent struct {
	Action string
	Relase struct {
		Created_At string
		HTML_URL   string
	}
	Assets []struct {
		ID string
	}
}

type SponsorshipEvent struct {
	Action         string
	Effective_Date string
}

func main() {
	args := os.Args[1:]
	prettyDateFormat := "Mon, Jan 2, 2006 at 3:04 PM"

	if len(args) != 1 {
		fmt.Println("Usage:  githubactivity <username>")
		os.Exit(1)
	}
	// using github.com/mrz1836/go-sanitize to sanitize the username
	username := sanitize.AlphaNumeric(args[0], true)
	url := sanitize.URL(fmt.Sprintf("https://api.github.com/users/%s/events", username))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Error: %s\n", resp.Status)
		os.Exit(1)
	}

	// Grab the body of the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// Grab the events type and parse it out
	var events []GithubEvent
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d events\n", len(events))

	// TODO: Refactor this to match with myabe an adapter pattern
	for i := 0; i < len(events); i++ {
		switch events[i].Type {
		case "CommitCommentEvent":
			var commentEvent CommitCommentEvent
			err = json.Unmarshal(events[i].Payload, &commentEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Comment created: %s\n\n", commentEvent.Comment)
			}

		case "CreateEvent":
			var createEvent CreateEvent
			err = json.Unmarshal(events[i].Payload, &createEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("A %s was created: %s\n\n", createEvent.Ref_Type, createEvent.Ref)
			}

		case "DeleteEvent":
			var deleteEvent DeleteEvent
			err = json.Unmarshal(events[i].Payload, &deleteEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Repository deleted: %s\n\n", deleteEvent.Ref)
			}

		case "ForkEvent":
			// A Fork is created
			var forkEvent ForkEvent
			err = json.Unmarshal(events[i].Payload, &forkEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				parsedTime, err := time.Parse(time.RFC3339, forkEvent.Forkee.Created_At)
				if err != nil {
					fmt.Printf("Repository forked at %s: %s\nURL: %s\n\n", forkEvent.Forkee.Created_At, forkEvent.Forkee.Name, forkEvent.Forkee.HTML_URL)

				} else {
					prettyDate := parsedTime.Local().Format(prettyDateFormat)
					fmt.Printf("Repository forked at %s: %s\nURL: %s\n\n", prettyDate, forkEvent.Forkee.Name, forkEvent.Forkee.HTML_URL)
				}
			}
		case "GollumEvent":
			// A wiki page is created or upadated
			var gollumEvent GollumEvent
			err = json.Unmarshal(events[i].Payload, &gollumEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				for j := 0; j < len(gollumEvent.Pages); j++ {
					fmt.Printf("Page %s: %s\nURL: %s\n\n", gollumEvent.Pages[j].Action, gollumEvent.Pages[j].Page_Name, gollumEvent.Pages[j].HTML_URL)
				}
			}

		case "IssueCommentEvent":
			var commentEvent IssueCommentEvent
			err = json.Unmarshal(events[i].Payload, &commentEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				parsedTime, err := time.Parse(time.RFC3339, commentEvent.Comment.Created_At)
				if err != nil {
					fmt.Printf("Comment created on %s: \n%s\nURL: %s\n\n", commentEvent.Comment.Created_At, commentEvent.Comment.Body, commentEvent.Comment.HTML_URL)
				} else {
					prettyDate := parsedTime.Local().Format(prettyDateFormat)
					fmt.Printf("Comment created on %s: \n%s\nURL: %s\n\n", prettyDate, commentEvent.Comment.Body, commentEvent.Comment.HTML_URL)
				}
			}

		case "IssuesEvent":
			var issueEvent IssuesEvent
			err = json.Unmarshal(events[i].Payload, &issueEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Issue %s: %s\nURL: %s\n\n", issueEvent.Action, issueEvent.Issue.Title, issueEvent.Issue.HTML_URL)
			}
		case "MemberEvent":
			// Activity related to repository collaborators
			var memberEvent MemberEvent
			err = json.Unmarshal(events[i].Payload, &memberEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Member %s changed role to %s\n\n", memberEvent.Member.Login, memberEvent.Changes.Role_Name.To)
			}

		case "PublicEvent":
			// When a private repository is made public
			fmt.Printf("Repository made public: %s\n\n", events[i].Repo.Name)
		case "PullRequestEvent":
			var pullRequestEvent PullRequestEvent
			err = json.Unmarshal(events[i].Payload, &pullRequestEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Pull request %s: \n Assignee: %s\nURL: %s\n\n", pullRequestEvent.Action, pullRequestEvent.Pull_Request.Assignee.Login, pullRequestEvent.Pull_Request.HTML_URL)
			}
		case "PullRequestReviewEvent":
			var pullRequestReviewEvent PullRequestReviewEvent
			err = json.Unmarshal(events[i].Payload, &pullRequestReviewEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Pull request review %s: \n URL: %s\n\n", pullRequestReviewEvent.Action, pullRequestReviewEvent.Review.HTML_URL)
			}
		case "PullRequestReviewCommentEvent":
			var PullRequestReviewCommentEvent PullRequestReviewCommentEvent
			err = json.Unmarshal(events[i].Payload, &PullRequestReviewCommentEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				parsedTime, err := time.Parse(time.RFC3339, PullRequestReviewCommentEvent.Comment.Created_At)
				if err != nil {
					fmt.Printf("Pull Request Review Commented on %s: \n%s\nURL: %s\n\n", PullRequestReviewCommentEvent.Comment.Created_At, PullRequestReviewCommentEvent.Comment.Body, PullRequestReviewCommentEvent.Comment.HTML_URL)
				} else {
					prettyDate := parsedTime.Local().Format(prettyDateFormat)
					fmt.Printf("Pull Request Review Commented on %s: \n%s\nURL: %s\n\n", prettyDate, PullRequestReviewCommentEvent.Comment.Body, PullRequestReviewCommentEvent.Comment.HTML_URL)
				}
			}
		case "PullRequestReviewThreadEvent":
			// Activity related to a comment thread on a pull request being marked as resolved or unresolved.
			var PullRequestReviewThreadEvent PullRequestReviewThreadEvent
			err = json.Unmarshal(events[i].Payload, &PullRequestReviewThreadEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Pull Request marked as %s: \n URL: %s\n\n", PullRequestReviewThreadEvent.Action, PullRequestReviewThreadEvent.Pull_Request.HTML_URL)
			}
		case "PushEvent":
			var pushEvent PushEvent
			err = json.Unmarshal(events[i].Payload, &pushEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Pushed %d commits to %s to the %s repository \n\n", pushEvent.Size, pushEvent.Ref, events[i].Repo.Name)
			}
		case "ReleaseEvent":
			var releaseEvent ReleaseEvent
			err = json.Unmarshal(events[i].Payload, &releaseEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Published release on %s with %d assets.  URL: %s\n\n", releaseEvent.Relase.Created_At, len(releaseEvent.Assets), releaseEvent.Relase.HTML_URL)
			}
		case "SponsorshipEvent":
			var sponsorshipEvent SponsorshipEvent
			err = json.Unmarshal(events[i].Payload, &sponsorshipEvent)
			if err != nil {
				fmt.Printf("Error interpreting event: %s\n\n", err)
			} else {
				fmt.Printf("Sponsorship created for %s.\n\n", events[i].Repo.Name)
			}
		case "WatchEvent":
			parsedTime, err := time.Parse(time.RFC3339, events[i].Created_At)
			if err != nil {
				fmt.Printf("Stared to watch repository @ %s: %s\n", events[i].Created_At, events[i].Repo.Name)
			} else {
				prettyDate := parsedTime.Local().Format(prettyDateFormat)
				fmt.Printf("Stared to watch repository @ %s: %s\n", prettyDate, events[i].Repo.Name)
			}
		default:
			fmt.Println("Unknown event type")
		}
	}
}
