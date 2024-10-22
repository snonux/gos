# Gos (Go Social Media)

Gos is a Go-based replacement for Buffer.com, providing the ability to schedule and manage social media posts from the command line.

## Features

* Mastodon and LinkedIn support
* OAuth2 authentication for LinkedIn
* Dry run mode for testing posts without actually publishing.
* Configurable via flags and environment variables.
* Easy to integrate into automated workflows.

## Installation

### Prerequisites

* Go (version 1.23 or later)
* Supported browsers like Firefox, Chrome, etc for oauth2.

### Build and Instal

Clone the repository:

```bash
git clone https://codeberg.org/snonux/gos.git
cd gos
```

Build the binary:

```bash
go build -o gos ./cmd/gos
sudo mv gos ~/go/bin
```

Or if you want to use the `Taskfile`:

```bash
go-task install
```

## Configuration

Gos requires a configuration file to store API secrets and OAuth2 credentials for each supported social media platform. The configuration is managed using a Secrets structure, which is stored as a JSON file in `~/.config/gos/gosec.json`.

Example Configuration File (`secrets.json`):

Below is an example of how your secrets.json file should be structured

```json
{
  "MastodonURL": "https://mastodon.example.com",
  "MastodonAccessToken": "your-mastodon-access-token",
  "LinkedInClientID": "your-linkedin-client-id",
  "LinkedInSecret": "your-linkedin-client-secret",
  "LinkedInRedirectURL": "http://localhost:8080/callback",
}
```

### Configuration Fields

* `MastodonURL`: The base URL of the Mastodon instance you are using (e.g., https://mastodon.social).
* `MastodonAccessToken`: Your access token for the Mastodon API, which is used to authenticate your posts.
* `LinkedInClientID`: The client ID for your LinkedIn app, which is needed for OAuth2 authentication.
* `LinkedInSecret`: The client secret for your LinkedIn app.
* `LinkedInRedirectURL`: The redirect URL configured for handling OAuth2 responses.
* `LinkedInAccessToken`: This will be automatically updated by Gos after successful OAuth2 authentication with LinkedIn.
* `LinkedInPersonID`: This will be automatically updated by Gos after successful OAuth2 authentication with LinkedIn.

### Automatically Managed Fields

Once you finish the OAuth2 setup, some fields—like `LinkedInAccessToken` and `LinkedInPersonID`—will get filled in automatically. To check if everything's working without actually posting anything, you can run the app in dry run mode with the `--dry` option. After OAuth2 is successful, Gos will update the file with `LinkedInClientID` and `LinkedInAccessToken`. And if the access token expires, it’ll just go through the OAuth2 process again.

## Invoking Gos

Gos is a command-line tool that lets you post updates to multiple social media platforms. You can run it with various flags to customize its behavior, such as posting in dry run mode, limiting posts by size, or targeting specific platforms.

Flags are used to control the tool's behavior. Below are several common ways to invoke Gos, along with descriptions of the available flags.

### Common Flags

* `--dry` Enables dry run mode, which simulates the posting process without actually sending posts to the platforms.
* `--gosDir <directory>`: Specifies the directory for the Gos queue/database. Default is `~/.gosdir`
* `--browser <browser>`: Specifies the OAuth2 browser to use for authentication (e.g., firefox, chrome).
* `--platform <platform>`: Specifies a comma separated list which social media platform to post to plus max message size (mastodon:500, linkedin:1000, default is all). 

### Examples

*Dry run mode*

Dry run mode lets you simulate the entire posting process without actually sending the posts. This is useful for testing configurations or seeing what would happen before making real posts.

```bash
./gos --dry
```

*Normal run*

Sharing to all platforms is as simple as the following (assuming it is configured correctly):

```bash
./gos 
```

:-)

However, you will notice that there are no messages queued to be posted yet. Read the next section of this README...

## Composing Messages to Be Posted

To post messages using Gos, you need to create text files that contain the content of the posts. These files are placed inside the directory specified by the --gosDir flag (the default directory is `~/.gosdir`). Each text file represents a single post and must have the .txt extension.

### Basic Structure of a Message File

Each text file should contain the message you want to post on the specified platforms. That's it. Example of a Basic Post File `~/.gosdir/samplepost.txt`:

```bash
This is a sample message to be posted on social media platforms.

Maybe add a link here: https://foo.zone

#foo #cool #gos #golang
```

The message is just arbitrary text, and Gos does not parse any of the content other than ensuring the overall allowed size for the social media platform isn't exceeded. If it exceeds the limit, Gos will prompt you to edit the post using your standard text editor (as specified by the `EDITOR` environment variable). All the hyperlinks, hashtags, etc., are interpreted by the social platforms themselves (e.g., Mastodon, LinkedIn) when posting.

### Adding share tags in the Filename

You can control which platforms the post is shared to and manage other behaviors using tags embedded in the filename.

To target specific platforms, you can add tags in the format share:platform1.-platform2 within the filename. This tells Gos to share the message only to platform1 (e.g., Mastodon) and explicitly exclude platform2 (e.g., LinkedIn). 

You can include multiple platforms by listing them after share:, separated by a .. Use the - symbol to exclude a platform.

*Examples:*

* To share only on Mastodon: `~/.gosdir/foopost.share:mastodon.txt`
* To not share on LinkedIn: `~/.gosdir/foopost.share:-linkedin.txt`
* To explicitly share on both: `~/.gosdir/foopost.share:linkedin:mastodon.txt`
* To explicitly share on only linkedin: `~/.gosdir/foopost.share:linkedin:-mastodon.txt`

### Using the `prio` Tag

Normally, Gos randomly picks any queued message without any specific order or priority. However, you can assign a higher priority to a message. The priority determines the order in which posts are processed, with messages without a priority tag being processed last and those with priority tags being processed first. If there are multiple messages with the priority tag, then a random message will be selected from them.

*Examples using the Priority tag:* 

* To share only on Mastodon: `~/.gosdir/foopost.prio.share:mastodon.txt`
* To not share on LinkedIn: `~/.gosdir/foopost.prio.share:-linkedin.txt`
* To explicitly share on both: `~/.gosdir/foopost.prio.share:linkedin:mastodon.txt`
* To explicitly share on only linkedin: `~/.gosdir/foopost.prio.share:linkedin:-mastodon.txt`

### Summary of Filename Structure

* The text file must be placed in the gosDir.
* Use the `.txt` extension (`.md` works as well as a hack, to make it compatible with Obsidian).
* The optional `share` tag controls which platforms the post goes to.
* The optional `prio` tag controls the priority of the post.

## How Queueing Works in Gos

When you place a message file in the gosDir, Gos processes it by moving the message through a queueing system before posting it to the target social media platforms. The lifecycle of a message includes several key stages, from creation to posting, all managed through the `./db/platforms/PLATFORM` directories.

### Step-by-Step Queueing Process

1. Inserting a Message into gosDir: You start by creating a text file that represents your post (e.g., `foo.txt`) and place it in the gosDir. This file is then processed by Gos when it runs.

2. Moving to the Queue: Upon running Gos, the tool identifies the message in the gosDir and places it into the queue for the specified platform. The message is moved into the appropriate directory for each platform in `./db/platforms/PLATFORM`. During this stage, the message file is renamed to include a timestamp indicating when it was queued and given a `.queued` extension.

*Example: If a message is queued for LinkedIn, the filename might look like this:*

```
./db/platforms/linkedin/foo.share:-mastodon.txt.20241022-102343.queued
```

3. Posting the Message: Once a message is placed in the queue, Gos takes care of posting it to the specified social media platforms. 

4. Renaming to `.posted`: After a message is successfully posted to a platform, the corresponding .queued file is renamed to have a .posted extension and the timestamp in the filename is updated as well. This signals that the post has been processed and published.

*Example: After a successful post to LinkedIn, the message file might look like this:*

```
./db/platforms/linkedin/foo.share:-mastodon.txt.20241112-121323.posted
```

##  How Message Selection Works in Gos

Gos uses a combination of priority, platform-specific tags, and timing rules to decide which messages to post. The message selection process ensures that messages are posted according to your configured cadence and targets, while respecting pauses between posts and previously met goals.

The key factors in message selection are:

* Message Priority: Messages with no priority value are processed after those with priority. If two messages have the same priority, one is selected randomly.
* Pause Between Posts: The `-pauseDays` flag allows you to specify a minimum number of days to wait between posts for the same platform. This prevents oversaturation of content and ensures that posts are spread out over time.
* Target Number of Posts Per Week: The `-target` flag defines how many posts per week should be made to a specific platform. This target helps Gos manage the rate of posting, ensuring that the right number of posts are made without exceeding the desired frequency. 
* Post History Lookback: The `-lookback` flag tells Gos how many days back to look in the post history to calculate whether the weekly post target has already been met. It ensures that previously posted content is considered before deciding to queue up another message.
