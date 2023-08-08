package commandlineprocessor

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/tekkamanendless/lab-youareell-2a/internal/youareellclient"
)

// Process a command.
func Process(client *youareellclient.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	arg := args[0]
	args = args[1:]

	var err error
	switch arg {
	case "help":
		err = processHelp(client, args)
	case "ids":
		err = processUsers(client, args)
	case "messages":
		err = processMessages(client, args)
	case "send":
		err = processSend(client, args)
	case "watch":
		err = processWatch(client, args)
	default:
		err = fmt.Errorf("unknown command: %s", arg)
	}
	return err
}

// processHelp processes the help-related commands.
func processHelp(client *youareellclient.Client, args []string) error {
	fmt.Printf("You may run any of these commands:\n")
	fmt.Printf("   help                          | Show this help.\n")
	fmt.Printf("   ids                           | List the users.\n")
	fmt.Printf("   ids <github-id>               | List the user ID for the given GitHub ID.\n")
	fmt.Printf("   messages                      | List the most recent messages.\n")
	fmt.Printf("   messages <github-id>          | List the messages for the given GitHub ID.\n")
	fmt.Printf("   send <from> <message>         | Send a message from the given GitHub ID.\n")
	fmt.Printf("   send <from> <message> to <to> | Send a message from the given GitHub ID to the other GitHub ID.\n")
	fmt.Printf("   watch                         | Watch for new messages.\n")
	fmt.Printf("   watch <github-id>             | Watch for new messages for the given GitHub ID.\n")
	return nil
}

// processMessages processes the message-related commands.
func processMessages(client *youareellclient.Client, args []string) error {
	var githubID string

	if len(args) == 0 {
		// No GitHub ID.
	} else if len(args) == 1 {
		githubID = args[0]
	} else {
		return fmt.Errorf("expected 0 or 1 argument; got %d", len(args))
	}

	endpoint := "/messages"
	if githubID != "" {
		endpoint = "/ids/" + url.PathEscape(githubID) + "/messages"
	}

	var output []youareellclient.Message
	err := client.Raw(http.MethodGet, endpoint, nil, &output)
	if err != nil {
		return err
	}
	for _, message := range output {
		fmt.Printf("Message: %s %s %s -> %s: %s\n", message.Timestamp, message.Sequence, message.FromID, message.ToID, message.Message)
	}

	return nil
}

// processSend processes the send-related commands.
func processSend(client *youareellclient.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("expected 2 args; got %d", len(args))
	}

	fromID := args[0]
	text := args[1]

	input := youareellclient.Message{
		Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"), // Apparently a timestamp is required.
		FromID:    fromID,
		Message:   text,
	}

	args = args[2:]
	if len(args) == 0 {
		// Great.
	} else if len(args) == 2 {
		if args[0] != "to" {
			return fmt.Errorf("expected 'to'; got %s", args[0])
		}
		input.ToID = args[1]
	} else {
		return fmt.Errorf("expected 5 args; got %d", 2+len(args))
	}

	var output youareellclient.Message
	err := client.Raw(http.MethodPost, "/ids/"+url.PathEscape(fromID)+"/messages", input, &output)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", output)

	return nil
}

// processUsers processes the user-related commands.
func processUsers(client *youareellclient.Client, args []string) error {
	if len(args) == 0 {
		var output []youareellclient.User
		err := client.Raw(http.MethodGet, "/ids", nil, &output)
		if err != nil {
			return err
		}
		for _, user := range output {
			fmt.Printf("User: %s: %s (%s)\n", user.UserID, user.GitHubID, user.Name)
		}
	} else {
		if len(args) == 1 {
			githubID := args[0]

			var output string
			err := client.Raw(http.MethodGet, "/ids/"+url.PathEscape(githubID), nil, &output)
			if err != nil {
				return err
			}
			fmt.Printf("ID: %s\n", output)
		} else if len(args) == 2 {
			name := args[0]
			githubID := args[1]

			input := youareellclient.User{
				UserID:   "-",
				Name:     name,
				GitHubID: githubID,
			}
			err := client.Raw(http.MethodPost, "/ids", input, nil)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("expected 1 or 2 args; got %d", len(args))
		}
	}
	return nil
}

// processWatch processes the watch-related commands.
func processWatch(client *youareellclient.Client, args []string) error {
	var githubID string

	if len(args) == 0 {
		// No GitHub ID.
	} else if len(args) == 1 {
		githubID = args[0]
	} else {
		return fmt.Errorf("expected 0 or 1 argument; got %d", len(args))
	}

	endpoint := "/messages"
	if githubID != "" {
		endpoint = "/ids/" + url.PathEscape(githubID) + "/messages"
	}

	var mostRecentSequence string
	for {
		var output []youareellclient.Message
		err := client.Raw(http.MethodGet, endpoint, nil, &output)
		if err != nil {
			return err
		}
		if len(output) > 0 {
			if mostRecentSequence != "" {
				for _, message := range output {
					if message.Sequence == mostRecentSequence {
						break
					}
					fmt.Printf("Message: %s %s %s -> %s: %s\n", message.Timestamp, message.Sequence, message.FromID, message.ToID, message.Message)
				}
			}
			mostRecentSequence = output[0].Sequence
		}

		time.Sleep(1 * time.Second)
	}

	// No return here; we'll be in an infinite loop above.
}
